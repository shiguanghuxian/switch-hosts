package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/shiguanghuxian/switch-hosts/program"
	"github.com/skratchdot/open-golang/open"
	"github.com/zserge/webview"
)

// InjectJs 此文件为注入js的对象，可以互操作
type InjectJs struct {
	program      *program.Program
	w            webview.WebView
	Body         string // 右侧hosts内容
	HostList     string // hosts列表
	ClipboardTxt string // 粘贴板内容
}

// NewInjectJs 创建注入对象
func NewInjectJs(p *program.Program) *InjectJs {
	return &InjectJs{
		program: p,
	}
}

// GetHostsByKey 获取指定key的hosts内容
func (js *InjectJs) GetHostsByKey(key string) {
	log.Println("获取hosts内容 " + key)
	basePath := js.GetRootDir()
	hostsPath := fmt.Sprintf("%shosts/%s.hosts", basePath, key)
	hostsBody, err := ioutil.ReadFile(hostsPath)
	if err != nil {
		js.w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "警告", err.Error())
		return
	}
	js.Body = string(hostsBody)
}

// GetHostsList 获取hosts列表
func (js *InjectJs) GetHostsList() {
	js.program.GetDB().ReLoad()
	body, err := json.Marshal(js.program.GetDB().Hosts)
	if err != nil {
		js.w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "警告", err.Error())
		js.HostList = "[]"
		return
	}
	js.HostList = string(body)
}

// ChangeHosts 切换hosts
func (js *InjectJs) ChangeHosts(key string) {
	var checkHosts *program.HostConfig
	var err error
	defer func() {
		if err != nil {
			js.w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "警告", err.Error())
			return
		}
	}()
	// 更新选中key
	for _, v := range js.program.GetDB().Hosts {
		if v.Key == key {
			js.program.GetDB().CheckHosts = v
			checkHosts = v
		}
	}
	if checkHosts == nil {
		err = errors.New("不存在此hosts配置")
		return
	}
	err = js.program.GetDB().SaveDb()
	if err != nil {
		return
	}

	// mac 重启状态栏程序
	js.RestartStateSwitchHosts()

	// 更新系统hosts文件
	err = js.program.UpHosts(checkHosts)
	if err != nil {
		return
	}
}

// AddHosts 添加hosts
func (js *InjectJs) AddHosts(data string) {
	val := new(program.HostConfig)
	var err error
	defer func() {
		if err != nil {
			js.w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "警告", err.Error())
			js.HostList = "[]"
			return
		}
	}()
	err = json.Unmarshal([]byte(data), val)
	if err != nil {
		return
	}
	// 写文件
	basePath := js.GetRootDir()
	hostsPath := fmt.Sprintf("%shosts/%s.hosts", basePath, val.Key)
	isExt, err := program.PathExists(hostsPath)
	if err != nil {
		return
	}
	if isExt == true {
		err = errors.New("已存在此key对应hosts")
		return
	}
	err = ioutil.WriteFile(hostsPath, []byte(""), 0664)
	if err != nil {
		return
	}
	// 保存db
	err = js.program.GetDB().AddHosts(val)
	if err != nil {
		return
	}

	// mac 重启状态栏程序
	js.RestartStateSwitchHosts()
}

// SaveHostsData 保存hosts参数
type SaveHostsData struct {
	Key  string `json:"key,omitempty"`
	Body string `json:"body,omitempty"`
}

// SaveHosts 保存hosts信息
func (js *InjectJs) SaveHosts(val string) {
	data := new(SaveHostsData)
	var err error
	defer func() {
		if err != nil {
			js.w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "警告", err.Error())
			js.HostList = "[]"
			return
		}
	}()
	err = json.Unmarshal([]byte(val), data)
	if err != nil {
		return
	}
	if data.Key == "" {
		err = errors.New("参数错误")
		return
	}
	// 保存数据
	basePath := js.GetRootDir()
	hostsPath := fmt.Sprintf("%shosts/%s.hosts", basePath, data.Key)
	err = ioutil.WriteFile(hostsPath, []byte(data.Body), 0644)
	if err != nil {
		return
	}
	// 双写
	if runtime.GOOS == "darwin" {
		// 如果包含，证明是 switch-hosts.app
		if strings.Index(basePath, "switch-hosts.app") > 0 {
			err = ioutil.WriteFile(fmt.Sprintf("/Applications/SwitchHosts.app/Contents/MacOS/hosts/%s.hosts", data.Key), []byte(data.Body), 0644)
		} else {
			err = ioutil.WriteFile(fmt.Sprintf(basePath+"switch-hosts.app/Contents/MacOS/hosts/%s.hosts", data.Key), []byte(data.Body), 0644)
		}
		if err != nil {
			log.Println(err)
			err = nil
		}
	}
	// 判断修改的是否是当前启用的key则更新hosts
	for _, v := range js.program.GetDB().Hosts {
		if v.Key == data.Key {
			if v.Check == true {
				err = js.program.UpHosts(v)
			}
			break
		}
	}

	js.w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "成功", "成功保存hosts配置，可以点击左侧开关应用此hosts")
}

// DelHosts 删除一个hosts
func (js *InjectJs) DelHosts(key string) {
	var err error
	defer func() {
		if err != nil {
			js.w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "警告", err.Error())
			return
		}
	}()
	if key == "" {
		err = errors.New("删除key不能为空")
		return
	}
	err = js.program.GetDB().DelHosts(key)
	if err != nil {
		return
	}

	// mac 重启状态栏程序
	js.RestartStateSwitchHosts()

	js.w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "成功", "成功删除hosts")
}

// GetRootDir 获取db和hosts配置文件所在目录
func (js *InjectJs) GetRootDir() string {
	basePath := program.GetRootDir()
	return basePath
}

// RestartStateSwitchHosts mac 重启状态栏程序
func (js *InjectJs) RestartStateSwitchHosts() {
	if runtime.GOOS == "darwin" {
		go func() {
			time.Sleep(2 * time.Second)
			err := open.Start(program.GetRootDir() + "switch-hosts.app")
			if err != nil {
				log.Println("打开switch-hosts.app错误")
			}
		}()
	}
}

// CtrlC 写入粘贴板
func (js *InjectJs) CtrlC(val string) {
	if val == "" {
		return
	}
	clipboard.WriteAll(val)
}

// CtrlV 读取粘贴板
func (js *InjectJs) CtrlV() {
	val, err := clipboard.ReadAll()
	if err != nil {
		log.Println(err)
		js.ClipboardTxt = ""
		return
	}
	js.ClipboardTxt = val
}
