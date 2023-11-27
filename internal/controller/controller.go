// Package controller @Author Zhan 2023/11/26 20:02:00
package controller

import (
	"fmt"
	"github.com/containerd/cgroups/v3/cgroup1"
	"github.com/opencontainers/runtime-spec/specs-go"
	"log"
)

const prefix = "tokamak"

type UserController struct {
	Uid      int32
	Path     string
	CGroup   cgroup1.Cgroup
	Resource *specs.LinuxResources
}

func NewUserController(uid int32, resource *specs.LinuxResources) (*UserController, error) {
	path := fmt.Sprintf("/%s-%v", prefix, uid)
	// if cgroup has existed
	var cgroup cgroup1.Cgroup
	cgroup, err := cgroup1.Load(cgroup1.StaticPath(path))
	if err == nil {
		log.Printf("cgroup %s has existed, add to user controller\n", path)
		if err2 := cgroup.Update(resource); err2 != nil {
			log.Println(err2)
		}
	} else {
		log.Printf("cgroup %s not exist, create new one\n", path)
		cgroup, err = cgroup1.New(cgroup1.StaticPath(path), resource)
		if err != nil {
			return nil, err
		}
	}
	return &UserController{
		Uid:      uid,
		Path:     path,
		CGroup:   cgroup,
		Resource: resource,
	}, nil
}

func (u *UserController) AddProcess(pid int) error {
	return u.CGroup.Add(cgroup1.Process{Pid: pid})
}

func (u *UserController) DeleteCGroup() error {
	return u.CGroup.Delete()
}
