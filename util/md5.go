package util

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func MakeMd5(srcStr string) string {
	salt := "dongne"
	srcStr += salt
	return fmt.Sprintf("%x", md5.Sum([]byte(srcStr)))
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
