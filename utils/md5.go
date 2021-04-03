package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

func GenMd5(srcStr string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(srcStr)))
}

func GenMd5WithSpecificLen(srcStr string, length int) string {
	hash := GenMd5(srcStr)
	var result string
	if len(hash) > length {
		result = hash[:length]
	} else {
		result = hash + strings.Repeat("0", length-len(hash))
	}
	return result
}
