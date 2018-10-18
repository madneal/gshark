package misc

import (
	"crypto/md5"
	"fmt"
)

func MakeMd5(srcStr string) string {
	salt := "dongne"
	srcStr += salt
	return fmt.Sprintf("%x", md5.Sum([]byte(srcStr)))
}
