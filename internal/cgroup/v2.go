// Package cgroup @Author Zhan 2023/12/6 14:09:00
package cgroup

import (
	"fmt"
	"github.com/containerd/cgroups/v3/cgroup2"
	"log"
)

var V2PREFIX = "sonmiv2"

type CGroupV2 struct {
	Uid      int32
	Path     string
	CGroup   *cgroup2.Manager
	Resource *Resource
}

func (C *CGroupV2) AddProcess(pid uint64) error {
	return C.CGroup.AddProc(pid)
}

func NewCGroupV2(uid int32, resource *Resource) (CGroup, error) {
	path := fmt.Sprintf("%s-%v.slice", V2PREFIX, uid)
	log.Printf("Quota: %v Period: %v", resource.CPU.Quota, resource.CPU.Period)
	log.Printf("MemoryMax: %v SwapMax: %v", resource.Memory.Max, resource.Memory.SwapMax)
	res := cgroup2.Resources{
		CPU: &cgroup2.CPU{
			Max: cgroup2.NewCPUMax(&resource.CPU.Quota, &resource.CPU.Period),
		},
		Memory: &cgroup2.Memory{
			Swap: &resource.Memory.SwapMax,
			Max:  &resource.Memory.Max,
		},
	}
	m, err := cgroup2.NewSystemd("/", path, -1, &res)
	if err != nil {
		return nil, err
	}
	_ = m.Update(&res)
	return &CGroupV2{
		Uid:      uid,
		Path:     path,
		CGroup:   m,
		Resource: resource,
	}, nil
}
