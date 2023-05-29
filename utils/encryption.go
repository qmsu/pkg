package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	uuid "github.com/satori/go.uuid"
	tsgutils "github.com/typa01/go-utils"
	"time"
)

func MD5(rawString string) string {
	h := md5.New()
	_, err := h.Write([]byte(rawString))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(str string) []int8 {
	sum224 := sha256.Sum256([]byte(str))
	s := make([]int8, sha256.Size)
	for i := range sum224 {
		s[i] = int8(sum224[i])
	}
	return s
}

func Bytes2Hex(bts []int8) string {
	var des string
	for i := 0; i < len(bts); i++ {
		//int(bts[i])&0xFF 就是将int8转换为uint8
		tmp := fmt.Sprintf("%x", int(bts[i])&0xFF)
		// fmt.Println("tmp:", int(bts[i])&0xFF)
		if len(tmp) == 1 {
			des = fmt.Sprintf("%s0", des)
		}
		des = fmt.Sprintf("%s%s", des, tmp)
	}
	return des
}

func GetUuid() (tokenString string) {
	return fmt.Sprintf("%s", uuid.Must(uuid.NewV4(), nil))
}

const baseDigits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const baseDigitsLen int64 = int64(len(baseDigits))

//构建36位的唯一主键码
func BuildKeyID() string {
	//62进制时间戳+UUID截取 定长36位，保证有序
	currentMillis := time.Now().UnixNano() / 1000000 //毫秒
	var sb []byte

	for {
		if currentMillis == 0 {
			break
		}
		sb = append(sb, baseDigits[int(currentMillis%baseDigitsLen)])
		currentMillis /= baseDigitsLen
	}

	sbLen := len(sb)
	for i := 0; i < (sbLen / 2); i++ {
		sb[i], sb[sbLen-i-1] = sb[sbLen-i-1], sb[i]
	}

	uuid := tsgutils.GUID()
	timePrefix := string(sb)
	return timePrefix + uuid[0:(36-len(timePrefix))]
}
