package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"virtual-target/internal"
	"virtual-target/log"
	"virtual-target/util"
)

type Input struct {
	Host       string
	ConfigFile string
}

func ParseInput() Input {
	var result Input
	var config string
	var host string
	var help bool
	flag.StringVar(&host, "H", "127.0.0.1", "Host (Default:127.0.0.1)")
	flag.StringVar(&config, "c", "", "Use Config File")
	flag.BoolVar(&help, "h", false, "Help Information")
	flag.Parse()
	if help == true {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if config != "" {
		if !util.FileExist(config) {
			log.Error("config file not found")
			os.Exit(0)
		} else {
			temp, _ := ioutil.ReadFile(config)
			result.ConfigFile = string(temp)
		}
	}
	result.Host = host
	return result
}

func main() {
	in := ParseInput()
	testCases := util.ResolveConfigFile(in.ConfigFile)
	if len(testCases) == 0 {
		internal.StartNoConfig(in.Host)
	} else {
		internal.StartUseConfig(in.Host, testCases)
	}
	wait()
}

func wait() {
	sign := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sign
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done
}
