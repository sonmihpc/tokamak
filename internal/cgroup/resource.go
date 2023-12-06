// Package cgroup @Author Zhan 2023/12/6 14:11:00
package cgroup

type Resource struct {
	CPU    *CPU
	Memory *Memory
}

type CPU struct {
	Quota  int64
	Period uint64
}

type Memory struct {
	Max              int64
	SwapMax          int64
	DisableOOMKiller bool
}
