// Package inspector @Author Zhan 2023/11/26 19:25:00
package inspector

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sonmihpc/tokamak/internal/cgroup"
	"github.com/sonmihpc/tokamak/internal/config"
	"log"
	"sync"
	"time"
)

var PERIOD = uint64(100000)

type Inspector struct {
	Interval            int
	Version             int
	status              bool
	closeCh             chan interface{}
	UserProcessGroupMap map[int32]*UserProcessGroup
	mapMu               sync.RWMutex
	Resource            *cgroup.Resource
	closeMu             sync.Mutex
	excludeUids         []int32
}

func (i *Inspector) Push(uid int32, u *UserProcessGroup) {
	i.mapMu.Lock()
	defer i.mapMu.Unlock()
	i.UserProcessGroupMap[uid] = u
}

func (i *Inspector) Delete(uid int32) {
	i.mapMu.Lock()
	defer i.mapMu.Unlock()
	delete(i.UserProcessGroupMap, uid)
}

func (i *Inspector) Existed(uid int32) bool {
	i.mapMu.RLock()
	defer i.mapMu.RUnlock()
	_, existed := i.UserProcessGroupMap[uid]
	return existed
}

func (i *Inspector) GetUserProcessGroup(uid int32) *UserProcessGroup {
	i.mapMu.RLock()
	defer i.mapMu.RUnlock()
	return i.UserProcessGroupMap[uid]
}

func (i *Inspector) beforeUpdate() {
	for _, u := range i.UserProcessGroupMap {
		u.beforeUpdate()
	}
}

func (i *Inspector) afterUpdate() {
	for _, u := range i.UserProcessGroupMap {
		u.afterUpdate()
	}
}

func (i *Inspector) update() {
	i.beforeUpdate()
	defer i.afterUpdate()
	processes, err := process.Processes()
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(processes))
	for _, p := range processes {
		go i.updateProcess(p, &wg)
	}
	wg.Wait()
}

func (i *Inspector) updateProcess(p *process.Process, wg *sync.WaitGroup) {
	defer wg.Done()
	uids, err := p.Uids()
	if err != nil {
		log.Println(err)
		return
	}
	uid := uids[0]
	if !i.Existed(uid) {
		enabled := false
		if uid >= 1000 && i.outOfExcludeUids(uid) {
			enabled = true
		}
		var cg cgroup.CGroup
		if enabled {
			cg, err = cgroup.NewCGroup(uid, i.Resource, i.Version)
			if err != nil {
				log.Printf("fail to create new cg for uid=%v\n", uid)
				log.Println(err)
				return
			}
		}
		i.Push(uid, &UserProcessGroup{
			Uid:       uid,
			enabled:   enabled,
			Processes: make(map[int32]*Process),
			CGroup:    cg,
		})
		//log.Printf("create user process group: uid=%v\n", uid)
	}
	upg := i.GetUserProcessGroup(uid)
	if !upg.Existed(p.Pid) {
		upg.Push(p.Pid, &Process{
			Updated: true,
			Dirty:   true,
			Process: p,
		})
		if upg.enabled {
			if err := upg.CGroup.AddProcess(uint64(p.Pid)); err != nil {
				log.Println(err)
				log.Printf("fail to add process %v to CGroup\n", p.Pid)
			}
			//log.Printf("add process %v to CGroup\n", p.Pid)
		}
		//log.Printf("create process: uid=%v pid=%v\n", uid, p.Pid)
		return
	}
	upg.SetUpdated(p.Pid, true)
	upg.SetDirty(p.Pid, false)
	//log.Printf("update process: pid=%v\n", p.Pid)
}

func (i *Inspector) RunInBackground() {
	ticker := time.NewTicker(time.Duration(i.Interval) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			i.update()
		case <-i.closeCh:
			log.Println("close CGroup update process")
			return
		default:
			time.Sleep(time.Duration(i.Interval) * time.Millisecond)
			// do nothing
		}
	}
}

func (i *Inspector) Run() bool {
	i.closeMu.Lock()
	defer i.closeMu.Unlock()
	i.status = false
	i.closeCh = make(chan interface{})
	i.UserProcessGroupMap = make(map[int32]*UserProcessGroup)
	go i.RunInBackground()
	i.status = true
	return i.status
}

func (i *Inspector) Stop() bool {
	i.closeMu.Lock()
	defer i.closeMu.Unlock()
	i.closeCh <- true
	i.status = false
	return i.status
}

func (i *Inspector) GetStatus() bool {
	return i.status
}

func (i *Inspector) outOfExcludeUids(uid int32) bool {
	for _, u := range i.excludeUids {
		if uid == u {
			return false
		}
	}
	return true
}

func NewInspector(interval int, version int, resource *cgroup.Resource, uids []int32) *Inspector {
	return &Inspector{
		Interval:    interval,
		Version:     version,
		Resource:    resource,
		excludeUids: uids,
	}
}

func NewInspectorFromCfg(conf *config.Config) *Inspector {
	cfg := conf.CGroup
	logicCnt, err := cpu.Counts(true)
	if err != nil {
		return nil
	}
	cpuQuota := (float64(logicCnt) * cfg.UserCpuPercent / 100) * float64(PERIOD)
	log.Printf("CPU Quota: %v\n", cpuQuota)
	vm, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	memLimit := int64(float64(vm.Total) * cfg.UserMemPercent / 100)
	log.Printf("Memory Limit: %v of total %v\n", memLimit, vm.Total)
	swapMax := int64(float64(vm.SwapTotal) * cfg.UserSwapPercent / 100)
	log.Printf("Swap LImit: %v of total %v\n", swapMax, vm.SwapTotal)
	disableOOMKiller := cfg.DisableOOMKiller
	log.Printf("Disable OOM Killer: %v", disableOOMKiller)
	res := &cgroup.Resource{
		CPU: &cgroup.CPU{
			Quota:  int64(cpuQuota),
			Period: PERIOD,
		},
		Memory: &cgroup.Memory{
			Max:              memLimit,
			SwapMax:          swapMax,
			DisableOOMKiller: disableOOMKiller,
		},
	}
	version := cgroup.GetCGroupVersion()
	return NewInspector(cfg.CheckIntervalMs, int(version), res, cfg.ExcludeUids)
}
