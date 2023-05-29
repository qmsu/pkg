//go:build darwin
// +build darwin

package utils

import (
	"bufio"
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func CheckPort(port string) bool {
	checkStatement := fmt.Sprintf(`netstat -anp | grep -q %s ; echo $?`, port)
	output, err := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	if err != nil {
		return false
	}
	result, err := strconv.Atoi(strings.TrimSuffix(string(output), "\n"))
	if err != nil {
		return false
	}
	if result == 0 {
		return true
	}
	return false
}

func GetTCPCount() (count int, err error) {
	output, err := exec.Command("sh", "-c", "netstat -ant|grep tcp|wc -l").CombinedOutput()
	if err != nil {
		return 0, err
	}
	result, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(string(output), "\n")))
	if err != nil {
		return 0, err
	}
	return result, nil
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

//mac/linux 生成md5
//md5sum /Users/zg/Downloads/NFgkwpYXBI5HXd21RbHnh0zVtHbaJZnRYPfy.mp4
//6a983ad48f9e97aee0bfa7d34e64bbed /Users/zg/Downloads/NFgkwpYXBI5HXd21RbHnh0zVtHbaJZnRYPfy.mp4
func GetFileMd5(fileDir string) (string, error) {
	fileDir = strings.ReplaceAll(fileDir, "(", "\\(")
	fileDir = strings.ReplaceAll(fileDir, ")", "\\)")
	fileDir = strings.ReplaceAll(fileDir, "（", "\\（")
	fileDir = strings.ReplaceAll(fileDir, "）", "\\）")
	cmdStr := fmt.Sprintf("md5sum %s", fileDir)
	output, err := exec.Command("sh", "-c", cmdStr).CombinedOutput()
	if err != nil {
		return "", err
	}
	var fileMd5 string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		arr := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(arr) > 0 {
			fileMd5 = strings.TrimSpace(arr[0])
		}
	}
	if fileMd5 == "" {
		return "", fmt.Errorf("生成md5失败")
	}
	return fileMd5, nil
}
