package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/shiguanghuxian/switch-hosts/program"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	// 系统日志显示文件和行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// 打开状态栏程序
	go func() {
		err := open.Start(program.GetRootDir() + "switch-hosts.app")
		if err != nil {
			log.Println("打开switch-hosts.app错误")
		}
	}()
	// 状态栏程序事例
	p, err := program.New("")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// 编辑界面程序
	err = NewEdit(p).Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("启动程序成功")

	// 监听退出信号
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	p.Stop()
	log.Println("Exit")
}
