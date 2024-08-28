package main

import (
	"dbcheck/config"
	"dbcheck/service"
	"flag"
	"log"
)

func main() {
	// 加载配置文件
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.Parse()


	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf( "%v", err)
	}

	// 执行巡检函数
	err = service.CheckMysql(cfg)

	if err != nil {
		log.Printf( "%v", err)
	}
}
