package utils

import (
	"github.com/shirou/gopsutil/v3/disk"
)

//判断是否是Windows C盘
func CheckWinC(p string) bool {
	if p == "c:" || p == "C:" {
		return true
	}
	return false
}

type Disk struct {
	Free        uint64  `json:"free"`         //空闲空间
	Used        uint64  `json:"used"`         //已使用空间
	Total       uint64  `json:"total"`        //总空间
	Path        string  `json:"path"`         //盘符路径
	UsedPercent float64 `json:"used_percent"` //已使用百分比
}

//获取系统磁盘使用情况
func Dsik() ([]*Disk, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	var usage []*Disk
	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}
		if u.Total/1024/1024/1024 < 1 {
			continue
		}
		//fmt.Println(part.Mountpoint + "_" + part.Device + "_" + part.String())
		//fmt.Println(u.Path + "\t" + strconv.FormatFloat(u.UsedPercent, 'f', 2, 64) + "% full.") //37.29%
		//fmt.Println("Total: " + strconv.FormatUint(u.Total/1024/1024/1024, 10) + " GiB")
		//fmt.Println("Free:  " + strconv.FormatUint(u.Free/1024/1024/1024, 10) + " GiB")
		//fmt.Println("Used:  " + strconv.FormatUint(u.Used/1024/1024/1024, 10) + " GiB")
		usage = append(usage, &Disk{
			Used:        u.Used,
			Total:       u.Total,
			Free:        u.Free,
			Path:        u.Path,
			UsedPercent: u.UsedPercent,
		})
	}
	return usage, nil
}
