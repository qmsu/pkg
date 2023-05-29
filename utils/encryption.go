package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
