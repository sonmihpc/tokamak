// Package config @Author Zhan 2023/11/26 20:01:00
package config

type CGroupCfg struct {
	CheckIntervalMs  int     `mapstructure:"check-interval-ms" yaml:"check-interval-ms"`
	UserCpuPercent   float64 `mapstructure:"user-cpu-percent" yaml:"user-cpu-percent"`
	UserMemPercent   float64 `mapstructure:"user-mem-percent" yaml:"user-mem-percent"`
	UserSwapPercent  float64 `mapstructure:"user-swap-percent" yaml:"user-swap-percent"`
	DisableOOMKiller bool    `mapstructure:"disable-oom-killer" yaml:"disable-oom-killer"`
	ExcludeUids      []int32 `mapstructure:"exclude-uids" yaml:"exclude-uids"`
}

type Config struct {
	CGroup CGroupCfg `mapstructure:"cgroup" yaml:"cgroup"`
}
