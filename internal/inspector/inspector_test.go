// @Author Zhan 2023/11/26 19:55:00
package inspector

import (
	"fmt"
	"github.com/sonmihpc/tokamak/internal/cgroup"
	"testing"
	"time"
)

func TestNewInspector(t *testing.T) {
	res := &cgroup.Resource{
		CPU: &cgroup.CPU{
			Quota:  400000,
			Period: PERIOD,
		},
		Memory: &cgroup.Memory{
			Max:              5000000,
			SwapMax:          5000000,
			DisableOOMKiller: true,
		},
	}
	inspector := NewInspector(1000, 1, res, []int32{})
	fmt.Println("first running...")
	inspector.Run()
	time.Sleep(time.Second * 1000)
	fmt.Println("stopping...")
	inspector.Stop()
	time.Sleep(time.Second * 10)
	fmt.Println("running again...")
	inspector.Run()
	time.Sleep(time.Second * 10)
}
