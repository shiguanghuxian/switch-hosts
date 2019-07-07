package main

import (
	"log"
	"net"
	"net/http"

	"github.com/shiguanghuxian/switch-hosts/program"
	"github.com/zserge/webview"
)

const (
	windowWidth  = 800
	windowHeight = 550
)

// Edit 编辑各个环境hosts程序
type Edit struct {
	program  *program.Program
	injectJs *InjectJs // js注入对象
	addr     string    // http服务监听地址
}

// NewEdit 创建修改hosts窗口
func NewEdit(p *program.Program) (et *Edit) {
	return &Edit{
		program:  p,
		injectJs: NewInjectJs(p),
	}
}

// Run 运行程序
func (et *Edit) Run() (err error) {
	err = et.startServer()
	if err != nil {
		return
	}
	log.Println(et.addr + "/index.html")
	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  "编辑HOSTS",
		Resizable:              true,
		URL:                    et.addr + "/index.html",
		ExternalInvokeCallback: et.handleRPC,
	})
	defer w.Exit()
	et.injectJs.w = w // 注入对象负值
	w.Dispatch(func() {
		// Inject controller
		w.Bind("injectJs", et.injectJs)
	})
	w.Run()
	return
}

// http服务
func (et *Edit) startServer() (err error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return
	}
	go func() {
		defer ln.Close()
		http.Handle("/", http.FileServer(http.Dir(program.GetRootDir()+"html/")))
		log.Fatal(http.Serve(ln, nil))
	}()
	et.addr = "http://" + ln.Addr().String()
	return nil
}

// 处理js调用
func (et *Edit) handleRPC(w webview.WebView, data string) {
	switch {
	default:
		log.Println(data)
	}
}
