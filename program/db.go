package program

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

// DB 负责操作db.json
type DB struct {
	filename   string
	CheckHosts *HostConfig   // 选中的hosts
	Hosts      []*HostConfig // 全部hosts
}

// HostConfig 一个自定义hosts
type HostConfig struct {
	Name  string `json:"name,omitempty"`
	Key   string `json:"key,omitempty"`
	Check bool   `json:"check,omitempty"`
}

// NewDB 创建db对象
func NewDB(filename string) (db *DB, err error) {
	if filename == "" {
		filename = GetRootDir() + "db.json"
	}
	db = &DB{
		filename: filename,
	}
	err = db.ReLoad()
	return
}

// ReLoad 重新加载db文件
func (db *DB) ReLoad() (err error) {
	body, err := ioutil.ReadFile(db.filename)
	if err != nil {
		return
	}
	hosts := make([]*HostConfig, 0)
	err = json.Unmarshal(body, &hosts)
	if err != nil {
		return
	}
	db.Hosts = hosts
	// 找到选中的hosts
	for _, v := range hosts {
		if v.Check == true {
			db.CheckHosts = v
		}
	}
	if len(hosts) > 0 {
		if db.CheckHosts == nil {
			db.CheckHosts = hosts[0]
		}
	}

	return
}

// SaveDb 写db文件
func (db *DB) SaveDb() (err error) {
	// 更新选中的hosts
	for _, v := range db.Hosts {
		if v.Key == db.CheckHosts.Key {
			v.Check = true
		} else {
			v.Check = false
		}
	}
	body, err := json.Marshal(db.Hosts)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(db.filename, body, 0644)
	// 同步写另一个 db.json
	if runtime.GOOS == "darwin" {
		var _err error
		// 如果包含，证明是 switch-hosts.app
		if strings.Index(db.filename, "switch-hosts.app") > 0 {
			_err = ioutil.WriteFile("/Applications/SwitchHosts.app/Contents/MacOS/db.json", body, 0644)
		} else {
			_err = ioutil.WriteFile(GetRootDir()+"switch-hosts.app/Contents/MacOS/db.json", body, 0644)
		}
		if _err != nil {
			log.Println(_err)
		}
	}
	return
}

// AddHosts 添加一个hosts配置
func (db *DB) AddHosts(val *HostConfig) (err error) {
	if val == nil {
		err = errors.New("添加hosts信息不能为空")
		return
	}
	db.Hosts = append(db.Hosts, val)
	err = db.SaveDb()
	return
}

// DelHosts 删除一个hosts配置
func (db *DB) DelHosts(key string) (err error) {
	if key == "" {
		err = errors.New("key不能为空")
		return
	}
	hosts := make([]*HostConfig, 0)
	for _, v := range db.Hosts {
		if v.Key != key {
			hosts = append(hosts, v)
		}
	}
	db.Hosts = hosts
	err = db.SaveDb()
	if err != nil {
		return
	}
	// 同步删文件
	err = os.Remove(fmt.Sprintf(GetRootDir()+"hosts/%s.hosts", key))
	if err != nil {
		return
	}
	if runtime.GOOS == "darwin" {
		var _err error
		// 如果包含，证明是 switch-hosts.app
		if strings.Index(db.filename, "switch-hosts.app") > 0 {
			_err = os.Remove(fmt.Sprintf("/Applications/SwitchHosts.app/Contents/MacOS/hosts/%s.hosts", key))
		} else {
			_err = os.Remove(fmt.Sprintf(GetRootDir()+"switch-hosts.app/Contents/MacOS/hosts/%s.hosts", key))
		}
		if _err != nil {
			log.Println(_err)
		}
	}
	return
}
