package program

import "io/ioutil"

// 写非windows hosts
func (p *Program) writeHosts(body []byte) (err error) {
	err = ioutil.WriteFile("C:\\windows\\system32\\drivers\\etc\\hosts", body, 0644)
	return
}
