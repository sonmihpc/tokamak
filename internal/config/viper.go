// Package config @Author Zhan 2023/11/27 16:42:00
package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

func Viper(path ...string) Config {
	var configPath string
	var config Config
	if len(path) == 0 {
		flag.StringVar(&configPath, "c", "", "Read configuration from the specified file.")
		flag.Parse()
		if configPath == "" {
			panic("Please set the configuration file by -c [conf_path]")
		} else {
			fmt.Printf("Use the configuration file from command flag.\n")
		}
	} else {
		configPath = path[0]
		fmt.Printf("Use the configuration specified by argument.\n")
	}

	v := viper.New()
	fmt.Printf("config path: %s\n", configPath)
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		panic("Fail to read configuration file.")
	}

	if err := v.Unmarshal(&config); err != nil {
		panic("Fail to unmarshal configuration.")
	}
	return config
}
