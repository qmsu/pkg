//go:build windows
// +build windows

package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func CheckPort(port string) bool {
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr %s", port)
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = &outBytes
	cmd.Run()
	resStr := outBytes.String()
	r := regexp.MustCompile(`\s\d+\s`).FindAllString(resStr, -1)
	if len(r) > 0 {
		_, err := strconv.Atoi(strings.TrimSpace(r[0]))
		if err != nil {
			fmt.Println("err", err.Error())
			return true
		}
		//if pid != -1 {
		//	c := exec.Command("taskkill.exe", "-F", "-PID", fmt.Sprintf("%d", pid))
		//	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		//	//c := exec.Command("taskkill.exe", "/f", "/im", "toolbox.exe")
		//	err = c.Start()
		//	if err != nil {
		//		return true
		//	}
		//}
		return true
	}
	return false
}

func GetTCPCount() (count int, err error) {
	var outBytes bytes.Buffer
	cmd := exec.Command("cmd", "/c", "netstat -ant | findstr /C /I TCP|findstr ESTABLISHED")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = &outBytes
	err = cmd.Run()
	if err != nil {
		return 0, err
	}
	scan := bufio.NewScanner(strings.NewReader(outBytes.String()))
	for scan.Scan() {
		count++
	}
	return count, nil
}

func GetNetReception() (download, upload int64, err error) {
	info, err := net.IOCounters(false)
	if err != nil {
		return download, upload, err
	}
	if len(info) == 0 {
		return download, upload, nil
	}
	time.Sleep(time.Second)
	info2, err := net.IOCounters(false)
	if err != nil {
		return download, upload, err
	}
	if len(info2) == 0 {
		return download, upload, nil
	}
	download = int64(info2[0].BytesRecv - info[0].BytesRecv) //下载
	upload = int64(info2[0].BytesSent - info[0].BytesSent)   //上传
	return download, upload, nil
}

/**
C:\Users\tool>certutil -hashfile C:\Users\tool\Desktop\小艺帮下载器2.12.4.138.exe
SHA1 的 C:\Users\tool\Desktop\小艺帮下载器2.12.4.138.exe 哈希:
1b8934d799661a04e705c46d9d07fe4df5c96706
CertUtil: -hashfile 命令成功完成。
*/
func GetFileMd5(fileDir string) (string, error) {
	dir := filepath.Dir(fileDir)
	cmdStr := fmt.Sprintf("cd /d %s && certutil -hashfile %s", dir, filepath.Base(fileDir))
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	index := 0
	var fileMd5 string
	for scanner.Scan() {
		index++
		if index == 2 {
			fileMd5 = strings.TrimSpace(scanner.Text())
			break
		}
	}
	if fileMd5 == "" {
		return "", fmt.Errorf("生成md5失败")
	}
	return fileMd5, nil
}
