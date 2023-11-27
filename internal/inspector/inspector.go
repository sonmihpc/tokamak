// Package inspector @Author Zhan 2023/11/26 19:25:00
package inspector

import (
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sonmihpc/tokamak/internal/config"
	"github.com/sonmihpc/tokamak/internal/controller"
	"log"
	"time"
)

var PERIOD = uint64(100000)

type Process struct {
	Updated bool
	Dirty   bool
	Process *process.Process
}

type SysUserInfo struct {
	Uid            int32
	enable         bool
	Processes      map[int32]*Process // index by pid
	UserController *controller.UserController
}

func (s *SysUserInfo) beforeUpdate() {
	for _, p := range s.Processes {
		p.Updated = false
	}
}

func (s *SysUserInfo) afterUpdate() {
	for _, p := range s.Processes {
		if p.Updated == false {
			log.Printf("process had been removed: uid=%v pid=%v", s.Uid, p.Process.Pid)
			delete(s.Processes, p.Process.Pid)
		}
	}
}

type Inspector struct {
	Interval int // update interval
	Users    map[int32]*SysUserInfo
	Resource *specs.LinuxResources
}

func (s *Inspector) beforeUpdate() {
	for _, u := range s.Users {
		u.beforeUpdate()
	}
}

func (s *Inspector) afterUpdate() {
	for _, u := range s.Users {
		u.afterUpdate()
	}
}

func (s *Inspector) update() {
	s.beforeUpdate()
	defer s.afterUpdate()
	processes, err := process.Processes()
	if err != nil {
		log.Println(err)
		return
	}
	for _, p := range processes {
		uids, err := p.Uids()
		if err != nil {
			log.Println(err)
			continue
		}
		uid := uids[0]
		_, ok := s.Users[uid]
		if !ok {
			enabled := false
			if uid >= 1000 {
				enabled = true
			}
			var usrCtrl *controller.UserController
			if enabled {
				usrCtrl, err = controller.NewUserController(uid, s.Resource)
				if err != nil {
					log.Printf("fail to create new usercontroller for uid=%v\n", uid)
					fmt.Println(err)
					continue
				}
			}
			s.Users[uid] = &SysUserInfo{
				Uid:            uid,
				enable:         enabled,
				Processes:      make(map[int32]*Process),
				UserController: usrCtrl,
			}
			log.Printf("create SysUserInfo: uid=%v\n", uid)
		}
		_, ok = s.Users[uid].Processes[p.Pid]
		if !ok { // new process need update cgroup
			s.Users[uid].Processes[p.Pid] = &Process{
				Updated: true,
				Dirty:   true,
				Process: p,
			}
			if s.Users[uid].enable {
				if err := s.Users[uid].UserController.AddProcess(int(p.Pid)); err != nil {
					log.Printf("fail to add process %v to CGroup\n", p.Pid)
				}
				log.Printf("add process %v to CGroup\n", p.Pid)
			}
			log.Printf("create Process: uid=%v pid=%v\n", uid, p.Pid)
			continue
		}
		// if process still running
		s.Users[uid].Processes[p.Pid].Updated = true
		s.Users[uid].Processes[p.Pid].Dirty = false
		//log.Printf("update process: pid=%v\n", p.Pid)
	}
}

func (s *Inspector) CheckInBackground() {
	ticker := time.NewTicker(time.Duration(s.Interval) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			s.update()
		default:
		}
	}
}

func NewInspector(config config.Config) *Inspector {
	// compute the cgroup cpu & memory setting params
	logicCnt, err := cpu.Counts(true)
	if err != nil {
		log.Printf("fail to get the number of logic cores.\n")
		return nil
	}
	// cpu_quota = logicCnt * user-cpu-percent / 100 * PERIOD
	cpuQuota := int64((float64(logicCnt) * config.CGroup.UserCpuPercent / 100) * float64(PERIOD))
	log.Printf("CPU Quota: %v\n", cpuQuota)
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("fail to get the virtual memory of the system.\n")
		return nil
	}
	memLimit := int64(float64(vm.Total) * config.CGroup.UserMemPercent / 100)
	log.Printf("Memory Limit: %v Total: %v\n", memLimit, vm.Total)
	disableOOMKiller := config.CGroup.DisableOOMKiller
	log.Printf("DisableOOMKiller: %v\n", disableOOMKiller)
	return &Inspector{
		Interval: config.CGroup.CheckIntervalMs,
		Users:    make(map[int32]*SysUserInfo),
		Resource: &specs.LinuxResources{
			CPU: &specs.LinuxCPU{
				Quota:  &cpuQuota,
				Period: &PERIOD,
			},
			Memory: &specs.LinuxMemory{
				Limit:            &memLimit,
				DisableOOMKiller: &disableOOMKiller,
			},
		},
	}
}
