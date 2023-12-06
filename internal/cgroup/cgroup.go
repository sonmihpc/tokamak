// Package controller @Author Zhan 2023/11/26 20:02:00
package cgroup

import (
	"errors"
	"github.com/containerd/cgroups/v3"
)

type CGroup interface {
	AddProcess(pid uint64) error
}

func NewCGroup(uid int32, resource *Resource, version int) (CGroup, error) {
	if version == 1 {
		return NewCGroupV1(uid, resource)
	}
	if version == 2 {
		return NewCGroupV2(uid, resource)
	}
	return nil, errors.New("invalid CGroup version")
}

type Version int

const (
	V1 Version = 1
	V2 Version = 2
)

func GetCGroupVersion() Version {
	if cgroups.Mode() == cgroups.Unified {
		return V2
	}
	return V1
}
