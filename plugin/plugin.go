package plugin

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Plugins struct {
	Name       string                `json:"name"`       //插件名称，和插件目录同名
	Alias      string                `json:"alias"`      //别名，用于在插件页面上展示
	Desc       string                `json:"desc"`       //插件介绍
	Exe        string                `json:"exe"`        //运行文件名
	HideOnList bool                  `json:"hideOnList"` //是否在列表隐藏
	Platform   map[string]PluginsRun `json:"platform"`   //key yun
}

type PluginsRun struct {
	Exe      string `json:"exe"`      //二进制文件或APP
	FileType int    `json:"fileType"` //文件类型：0APP，1命令行运行文件
}

type PluginsData struct {
	Dir       string
	Token     string
	CollegeId string
	ServerId  string
	Plugin    Plugins
	Data      string
}

type PluginsArgs struct {
	AppVersion string `json:"appVersion"`
	ServerUrl  string `json:"serverUrl"`
}

type Plugin struct {
	Data PluginsData
}

func Unmarshal(cmdData *string) (data PluginsData, err error) {
	if cmdData == nil {
		return data, errors.New("参数错误")
	}
	b, err := base64.StdEncoding.DecodeString(*cmdData)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
