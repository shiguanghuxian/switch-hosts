package program

import (
	"fmt"
	"io/ioutil"
	"os"
)

// 写mac hosts
func (p *Program) writeHosts(body []byte) (err error) {
	// 写临时文件
	err = ioutil.WriteFile(GetRootDir()+"hosts.cache", body, 0644)
	if err != nil {
		return
	}
	// 执行 applacript - 可以解决无权限写hosts文件问题
	script := fmt.Sprintf(`do shell script "cp %shosts.cache /private/etc/hosts" with administrator privileges`,
		GetRootDir(),
	)
	_, err = RunApplacript(script)
	if err != nil {
		return
	}
	// 删除临时文件
	err = os.Remove(GetRootDir() + "hosts.cache")
	return
}
