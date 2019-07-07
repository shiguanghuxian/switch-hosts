package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/shiguanghuxian/switch-hosts/program"
)

func main() {
	// 系统日志显示文件和行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// 程序实例
	p, err := program.New("")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// 运行
	log.Println("启动程序成功")
	p.Run()

	// 监听退出信号
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	p.Stop()
	log.Println("Exit")
}
