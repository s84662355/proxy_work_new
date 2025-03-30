package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"go.uber.org/zap"
	"mproxy/config"
	"mproxy/log"
	"mproxy/server"
)

func init() {
	var (
		configPath string
		mark       string
	)
	flag.StringVar(&configPath, "config_path", "", "配置文件路径")
	flag.StringVar(&mark, "mark", "", "标记")
	flag.Parse()

	fmt.Println(configPath)

	config.Load(configPath)
	fmt.Println(filepath.Join(config.GetConf().LogDir, config.ServiceName, mark))
	log.Init(filepath.Join(config.GetConf().LogDir, config.ServiceName, mark))
	log.Info(fmt.Sprintf("%s服务启动", config.ServiceName))
	go gohttp()
}

func main() {
	if err := server.Start(); err != nil {
		log.Fatal(fmt.Sprintf("启动服务%s失败", config.ServiceName), zap.Any("error", err))
		return
	}

	defer server.Stop()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-signalChan
}
