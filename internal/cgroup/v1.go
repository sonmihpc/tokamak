// Package cgroup @Author Zhan 2023/12/6 14:09:00
package cgroup

import (
	"fmt"
	"github.com/containerd/cgroups/v3/cgroup1"
	"github.com/opencontainers/runtime-spec/specs-go"
	"log"
)

var V1PREFIX = "sonmiv1"

type CGroupV1 struct {
	Uid      int32
	Path     string
	CGroup   cgroup1.Cgroup
	Resource *specs.LinuxResources
}

func (C *CGroupV1) AddProcess(pid uint64) error {
	return C.CGroup.Add(cgroup1.Process{Pid: int(pid)})
}

func NewCGroupV1(uid int32, resource *Resource) (CGroup, error) {
	path := fmt.Sprintf("/%s-%v", V1PREFIX, uid)
	log.Printf("Quota: %v Period: %v", resource.CPU.Quota, resource.CPU.Period)
	log.Printf("MemoryMax: %v SwapMax: %v", resource.Memory.Max, resource.Memory.SwapMax)
	res := &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Quota:  &resource.CPU.Quota,
			Period: &resource.CPU.Period,
		},
		Memory: &specs.LinuxMemory{
			Limit:            &resource.Memory.Max,
			Swap:             &resource.Memory.SwapMax,
			DisableOOMKiller: &resource.Memory.DisableOOMKiller,
		},
	}
	var cgroup cgroup1.Cgroup
	cgroup, err := cgroup1.Load(cgroup1.StaticPath(path))
	if err == nil {
		log.Printf("cgroup %s has existed, add to user controller\n", path)
		if err2 := cgroup.Update(res); err2 != nil {
			log.Println(err2)
		}
	} else {
		log.Printf("cgroup %s not exist, create new one\n", path)
		cgroup, err = cgroup1.New(cgroup1.StaticPath(path), res)
		if err != nil {
			return nil, err
		}
	}
	return &CGroupV1{
		Uid:      uid,
		Path:     path,
		CGroup:   cgroup,
		Resource: res,
	}, nil
}
