package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

func GenMd5(srcStr string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(srcStr)))
}
