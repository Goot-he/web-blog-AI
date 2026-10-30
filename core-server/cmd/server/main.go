package main

import (
	"core-server/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//正确的运行方式
//目录										命令						说明
//D:\web-blog-AI\core-server				go run ./cmd/server		推荐，按包编译
//D:\web-blog-AI\core-server\cmd\server		go run .				等价，. 表示当前包
//任意										go run main.go			❌ 单文件，必挂

func main() {
	// 初始化配置文件
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("load config error: ", err)
	}

	// 初始化server端
	srv, err := InitServer(cfg)
	if err != nil {
		log.Fatal("init server error: ", err)
	}

	// 启动server端
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down core-server...")
	srv.Stop()
}
