package program

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 公共函数库

// GetRootDir 获取当前可执行程序路径
func GetRootDir() string {
	// 文件不存在获取执行路径
	file, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		file = fmt.Sprintf(".%s", string(os.PathSeparator))
	} else {
		file = fmt.Sprintf("%s%s", file, string(os.PathSeparator))
	}
	return file
}

// PathExists 判断文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// RunApplacript 运行 applacript
func RunApplacript(script string) (string, error) {
	return command("osascript")(script)
}

// 执行指令
func command(name string, args ...string) func(...string) (string, error) {
	return func(input ...string) (string, error) {
		cmd := exec.Command(name, args...)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", err
		}

		io.WriteString(stdin, strings.Join(input, " "))
		stdin.Close()

		b, err := cmd.Output()

		return strings.TrimSpace(string(b)), err
	}
}
