package program

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
	notify "github.com/getlantern/notifier"
	"github.com/getlantern/systray"
	"github.com/shiguanghuxian/switch-hosts/program/icon"
	"github.com/skratchdot/open-golang/open"
)

// Program 程序实体
type Program struct {
	db        *DB
	MenuItems []*MenuItem // 菜单
}

// MenuItem 菜单对象
type MenuItem struct {
	menuItem *systray.MenuItem
	cfg      *HostConfig
}

// New 创建程序实例
func New(dbPath string) (*Program, error) {
	db, err := NewDB(dbPath)
	if err != nil {
		return nil, err
	}
	return &Program{
		db: db,
	}, nil
}

// Run 启动程序
func (p *Program) Run() {
	go p.BackupSystemHosts() // 备份系统hosts
	go p.watcher()           // 监听db.json变化
	systray.Run(p.onReady, p.Stop)
}

// Stop 程序结束要做的事
func (p *Program) Stop() {

}

// GetDB 获取db对象
func (p *Program) GetDB() *DB {
	return p.db
}

// AddMenuItems 添加菜单
func (p *Program) AddMenuItems() {
	// 将各个环境的host加入菜单
	for _, v := range p.db.Hosts {
		one := &MenuItem{
			cfg: v,
		}
		one.menuItem = systray.AddMenuItem(v.Name, v.Name)
		// 选中之前配置值
		if v.Key == p.db.CheckHosts.Key {
			one.menuItem.Check()
		}
		p.MenuItems = append(p.MenuItems, one)
	}
}

// 状态菜单和事件处理
func (p *Program) onReady() {
	systray.SetIcon(icon.Data)
	// systray.SetTitle("Switch Hosts")
	systray.SetTooltip("Switch Hosts")

	// 在协程中处理菜单事件
	go func() {
		// 添加菜单
		p.AddMenuItems()

		// 分割线
		systray.AddSeparator()
		mEdit := systray.AddMenuItem("编辑Hosts配置", "编辑各个环境hosts配置")

		// 分割线
		systray.AddSeparator()
		mAbout := systray.AddMenuItem("关于", "关于应用信息")
		mQuit := systray.AddMenuItem("退出程序", "")

		// 绑定事件
		for _, v := range p.MenuItems {
			go func(item MenuItem) {
				for {
					select {
					case <-item.menuItem.ClickedCh:
						p.changeHosts(&item)
					}
				}
			}(*v)
		}

		// 是否退出
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
				return
			case <-mAbout.ClickedCh:
				p.Notify("关于", "此程序可以更新hosts文件，用于开发测试使用。")
			case <-mEdit.ClickedCh:
				// err := open.Start(GetRootDir() + "hosts")
				// 打开编辑hosts界面程序
				// err := open.Start(GetRootDir() + "edit-hosts")
				var err error
				if runtime.GOOS == "darwin" {
					err = open.Start("/Applications/SwitchHosts.app")
				} else {
					err = open.Start("./SwitchHosts.exe")
				}
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

// 切换hosts文件
func (p *Program) changeHosts(item *MenuItem) {
	for _, v := range p.MenuItems {
		if v.cfg.Key != item.cfg.Key {
			v.menuItem.Uncheck()
		} else {
			v.menuItem.Check()
			// 更新db文件
			p.db.CheckHosts = item.cfg
			err := p.db.SaveDb()
			if err != nil {
				log.Println("写db.json错误", err)
			}
			// 更新hosts文件
			err = p.UpHosts(item.cfg)
			if err != nil {
				log.Println("更新hosts文件错误", err)
				p.Notify("失败", err.Error())
			} else {
				p.Notify("成功", fmt.Sprintf("已更新hosts为 %s", v.cfg.Name))
			}
		}
	}
}

// UpHosts 写入hosts
func (p *Program) UpHosts(v *HostConfig) (err error) {
	// 拼接最终需要设置的hosts信息
	basePath := GetRootDir()
	baseHostsPath := fmt.Sprintf("%shosts/base.hosts", basePath)
	baseHostsBody, err := ioutil.ReadFile(baseHostsPath)
	if err != nil {
		return
	}
	// 读取选中的hosts
	hostsPath := fmt.Sprintf("%shosts/%s.hosts", basePath, v.Key)
	hostsBody, err := ioutil.ReadFile(hostsPath)
	if err != nil {
		return
	}
	baseHostsBody = append(baseHostsBody, []byte(fmt.Sprintf("\n\n# %s \n", v.Name))...)
	hostsBody = append(baseHostsBody, hostsBody...)

	// log.Println(string(hostsBody))
	err = p.writeHosts(hostsBody)
	return
}

// Notify 弹框提示
func (p *Program) Notify(title string, message string) {
	n := notify.NewNotifications()
	msg := &notify.Notification{
		Title:   title,
		Message: message,
		Sender:  "cn.zuoxiupeng.switchhosts",
	}
	n.Notify(msg)
}

// Watcher 监听配置文件变化
func (p *Program) watcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("配置文件发生变化")
					if runtime.GOOS == "darwin" {
						// 配置变化直接退出，由界面程序拉起 - 需要排除状态栏程序
						if strings.Index(GetRootDir(), "switch-hosts.app") > 0 {
						} else {
							systray.Quit()
							os.Exit(0)
						}

					} else {
						err = p.GetDB().ReLoad()
						if err != nil {
							log.Println(err)
						} else {
							for _, v := range p.MenuItems {
								isCheck := false
								for _, vv := range p.GetDB().Hosts {
									if vv.Check == true && vv.Key == v.cfg.Key {
										isCheck = true
										break
									}
								}
								if isCheck == false {
									v.menuItem.Uncheck()
								} else {
									v.menuItem.Check()
								}
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(p.GetDB().filename)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// BackupSystemHosts 备份系统hosts
func (p *Program) BackupSystemHosts() {
	systemHostsPath := fmt.Sprintf(GetRootDir() + "hosts/system.hosts")
	isExt, _ := PathExists(systemHostsPath)
	if isExt == true {
		return
	}
	path := "/etc/hosts"
	if runtime.GOOS == "windows" {
		path = "C:\\windows\\system32\\drivers\\etc\\hosts"
	}
	body, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}
	err = ioutil.WriteFile(systemHostsPath, body, 0644)
	if err != nil {
		log.Println(err)
		return
	}
}
