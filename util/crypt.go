package util

import (
	"crypto/md5"
	"encoding/hex"
)

const salt = "w$1&kU2"

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text + salt))
	return hex.EncodeToString(hasher.Sum(nil))
}
