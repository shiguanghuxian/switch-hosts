package program

import "io/ioutil"

// 写linux hosts
func (p *Program) writeHosts(body []byte) (err error) {
	err = ioutil.WriteFile("/etc/hosts", body, 0644)
	return
}
