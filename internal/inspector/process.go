// Package inspector @Author Zhan 2023/12/6 14:12:00
package inspector

import (
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sonmihpc/tokamak/internal/cgroup"
	"log"
	"sync"
)

type Process struct {
	Updated bool
	Dirty   bool
	Process *process.Process
}

type UserProcessGroup struct {
	Uid       int32
	enabled   bool
	Processes map[int32]*Process
	CGroup    cgroup.CGroup
	mu        sync.RWMutex
}

func (u *UserProcessGroup) SetUpdated(pid int32, updated bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Processes[pid].Updated = updated
}

func (u *UserProcessGroup) SetDirty(pid int32, dirty bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Processes[pid].Dirty = dirty
}

func (u *UserProcessGroup) beforeUpdate() {
	for pid, _ := range u.Processes {
		u.SetUpdated(pid, false)
	}
}

func (u *UserProcessGroup) afterUpdate() {
	for _, p := range u.Processes {
		if p.Updated == false {
			log.Printf("process had been removed: uid=%v pid=%v\n", u.Uid, p.Process.Pid)
			u.Delete(p.Process.Pid)
		}
	}
}

func (u *UserProcessGroup) Push(pid int32, p *Process) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Processes[pid] = p
}

func (u *UserProcessGroup) Delete(pid int32) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.Processes, pid)
}

func (u *UserProcessGroup) Existed(pid int32) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	_, exist := u.Processes[pid]
	return exist
}
