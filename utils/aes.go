package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	tsgutils "github.com/typa01/go-utils"
	"time"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func Encode(str string, key string) string {
	pass := []byte(str)
	xpass, err := AesEncrypt(pass, []byte(key))
	if err != nil {
		return str
	}
	pass64 := base64.StdEncoding.EncodeToString(xpass)
	return pass64

}

func Decode(str string, key string) string {
	bytesPass, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return str
	}
	tpass, err := AesDecrypt(bytesPass, []byte(key))
	if err != nil {
		return str
	}

	return string(tpass)
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
