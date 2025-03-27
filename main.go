package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"mproxy/config"
	"mproxy/log"
	"mproxy/service"
)

func init() {
	var configPath string
	flag.StringVar(&configPath, "config_path", "", "配置文件路径")
	flag.Parse()

	config.Load(configPath)
	log.Init(filepath.Join(config.ConfData.LogDir, config.ServiceName))
	log.Infof("%s服务启动", config.ServiceName)
	go gohttp()
}

func main() {
	if err := service.Start(); err != nil {
		log.Fatalf("启动服务%s失败 err=%+v", config.ServiceName, err)
		return
	}

	defer service.Stop()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-signalChan
}
