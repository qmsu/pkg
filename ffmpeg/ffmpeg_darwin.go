//go:build darwin
// +build darwin

package ffmpeg

import (
	"bufio"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

func GetVideoDetail(videoPath string) (videoDetail VideoDetail, err error) {
	videoPath = strings.ReplaceAll(videoPath, "(", "\\(")
	videoPath = strings.ReplaceAll(videoPath, ")", "\\)")
	videoPath = strings.ReplaceAll(videoPath, "（", "\\（")
	videoPath = strings.ReplaceAll(videoPath, "）", "\\）")
	cmdStr := fmt.Sprintf("ffprobe -i %s", videoPath)
	output, err := exec.Command("sh", "-c", cmdStr).CombinedOutput()
	if err != nil {
		return videoDetail, err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(text, "Duration:") {
			arr := strings.Split(text, ",")
			if len(arr) > 0 {
				durationArr := strings.Split(arr[0], ":")
				if len(durationArr) > 0 {
					// durationArr[0] Duration
					// durationArr[1] 时
					// durationArr[2] 分
					// durationArr[3] 秒.毫秒
					hour, _ := strconv.Atoi(strings.TrimSpace(durationArr[1]))
					minute, _ := strconv.Atoi(strings.TrimSpace(durationArr[2]))
					tmp, err := strconv.ParseFloat(durationArr[3], 64)
					if err != nil {
						return videoDetail, fmt.Errorf("获取视频时长失败")
					}
					seconds := math.Ceil(tmp)
					videoDetail.VideoDuration = int64(hour*60*60 + minute*60 + int(seconds))
					break
				}
			}
		}
	}
	return videoDetail, nil
}
