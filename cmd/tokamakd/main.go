// @Author Zhan 2023/11/26 18:55:00
package main

import (
	"fmt"
	"github.com/coreos/go-systemd/daemon"
	"github.com/sonmihpc/tokamak/internal/config"
	"github.com/sonmihpc/tokamak/internal/inspector"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// load the config from file
	cfg := config.Viper()
	// create inspector instance from config
	spector := inspector.NewInspector(cfg)
	if spector == nil {
		fmt.Println("fail to create inspector, exit.")
		return
	}
	go spector.CheckInBackground()
	// daemon
	if _, err := daemon.SdNotify(false, "READY=1"); err != nil {
		log.Printf("notification supported, but failed.")
	}
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range ch {
			switch s {
			case syscall.SIGHUP, syscall.SIGILL, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("tokamak exit.")
				ch <- s
			default:
				log.Printf("received other signal: %s, tokamakd exit.", s)
			}
		}
	}()
	<-ch
}
