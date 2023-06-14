package utils

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

//去掉字符串中的换行符 \n 和 /
func Trim(str string) string {
	return strings.Replace(strings.Replace(str, "\n", "", -1), "/", " ", -1)
}

func StringToObject(str string, obj interface{}) error {
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		return err
	}
	return nil
}

func ObjectToString(obj interface{}) (string, error) {
	res, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func GetRoleNameByKey(objectKey string) string {
	arr := strings.Split(objectKey, "/")
	var roleName string
	for _, v := range arr {
		roleName = strings.TrimSpace(v)
		if roleName == "" {
			continue
		}
		break
	}
	return roleName
}

//转义字符串
//转义成安全的、可用于文件名或路径的字符串
// \n 替换成 ^#^
// \ 替换成 ^##^
// / 替换成 ^###^
// : 替换成 ^####^
// * 替换成 ^#####^
// ? 替换成 ^######^
// " 替换成 ^#######^
// < 替换成 ^########^
// > 替换成 ^#########^
// | 替换成 ^##########^
func PathEscape(str string) string {
	str = strings.ReplaceAll(str, "\n", "^#^")
	str = strings.ReplaceAll(str, `\`, "^##^")
	str = strings.ReplaceAll(str, `/`, "^###^")
	str = strings.ReplaceAll(str, ":", "^####^")
	str = strings.ReplaceAll(str, "*", "^#####^")
	str = strings.ReplaceAll(str, `?`, "^######^")
	str = strings.ReplaceAll(str, `"`, "^#######^")
	str = strings.ReplaceAll(str, "<", "^########^")
	str = strings.ReplaceAll(str, ">", "^#########^")
	str = strings.ReplaceAll(str, "|", "^##########^")
	return str
}

//转义字符串
//转义成安全的、可用于文件名或路径的字符串
// \n 替换成 ^#^
// \ 替换成 ^##^
// / 替换成 ^###^
// : 替换成 ^####^
// * 替换成 ^#####^
// ? 替换成 ^######^
// " 替换成 ^#######^
// < 替换成 ^########^
// > 替换成 ^#########^
// | 替换成 ^##########^
func PathUnEscape(str string) string {
	str = strings.ReplaceAll(str, "^#^", "\n")
	str = strings.ReplaceAll(str, "^##^", `\`)
	str = strings.ReplaceAll(str, "^###^", `/`)
	str = strings.ReplaceAll(str, "^####^", ":")
	str = strings.ReplaceAll(str, "^#####^", "*")
	str = strings.ReplaceAll(str, "^######^", `?`)
	str = strings.ReplaceAll(str, "^#######^", `"`)
	str = strings.ReplaceAll(str, "^########^", "<")
	str = strings.ReplaceAll(str, "^#########^", ">")
	str = strings.ReplaceAll(str, "^##########^", "|")
	return str
}

func IsInt8Contain(items []uint8, item uint8) bool {
	for _, eachItem := range items {
		if item == eachItem {
			return true
		}
	}
	return false
}

// RandStringRunes 返回随机字符串
func RandStringRunes(n int, isAllLower bool) string {
	// 去掉0 I l o 的字符串
	var letterRunes = []rune("123456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ")
	//避免windows时间每次拿到一致的盐值，休息10ms
	time.Sleep(10 * time.Millisecond)
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	randStr := string(b)
	if isAllLower {
		randStr = strings.ToLower(randStr)
	}
	return randStr
}
