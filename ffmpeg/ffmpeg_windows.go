//go:build windows
// +build windows

package ffmpeg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func GetVideoDetail(videoPath string) (videoDetail VideoDetail, err error) {
	tmpFileDir := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().Unix()))
	defer os.RemoveAll(tmpFileDir)
	dir := filepath.Dir(videoPath)
	cmd := exec.Command("cmd", "/c", "cd", "/d", dir, "&&", "ffprobe", "-i", videoPath, "-print_format", "json", "-show_format", ">", tmpFileDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err = cmd.CombinedOutput()
	if err != nil {
		return videoDetail, err
	}
	b, err := ioutil.ReadFile(tmpFileDir)
	if err != nil {
		return videoDetail, err
	}
	var ffprobeData Ffprobe
	err = json.Unmarshal(b, &ffprobeData)
	if err != nil {
		return videoDetail, err
	}
	f, err := strconv.ParseFloat(ffprobeData.Format.Duration, 64)
	if err != nil {
		return videoDetail, err
	}
	videoDetail.VideoDuration = int64(math.Ceil(f))
	videoDetail.FileSize, _ = strconv.ParseInt(ffprobeData.Format.Size, 10, 64)
	return videoDetail, nil
}
