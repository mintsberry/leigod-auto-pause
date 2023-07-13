package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"leigod-auto-pause/api"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var processName = []string{"OneDrive"}

func main() {
	err := api.Login("", "Lucky626")
	if err != nil {
		fmt.Printf("账户登陆失败:%s，程序自动退出", err)
		return
	}
	//processExists([]string{"hpc_pro"})
	//init()
	go timerHandler()
	signalHandler()
}

func init() {

}

func timerHandler() {
	duration := time.Duration(60 * 5)
	ticker := time.NewTicker(duration * time.Second)
	processCheckTicker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("暂停时长")
			api.Pause()
		case <-processCheckTicker.C:
			fmt.Println("检查进程是否存在")
			exists, _ := processExists(processName)
			if exists {
				fmt.Println("进程存在,重制暂停定时器")
				ticker.Reset(duration)
			}
		}
	}
}

func signalHandler() {
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTTIN)
	for {
		s := <-c
		fmt.Println(fmt.Sprintf("程序收到信号：%s", s.String()))
		api.Pause()
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		// TODO app reload
		default:
			return
		}
	}
}

func processExists(processNames []string) (bool, error) {
	processList, err := process.Processes()
	if err != nil {
		return false, err
	}
	for _, p := range processList {
		name, err := p.Name()
		if err == nil {
			for _, targetName := range processNames {
				if name == targetName {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
