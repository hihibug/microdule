package utils

import (
	"crypto/md5"

	"github.com/anaskhan96/go-password-encoder"
)

var PassOP password.Options

func init() {
	PassOP = password.Options{SaltLen: 10, Iterations: 10000, KeyLen: 50, HashFunction: md5.New}
}

// EncodeMd5Salt 生成md5和salt
func EncodeMd5Salt(pwd string) (string, string) {
	return password.Encode(pwd, &PassOP)
}

// VerifyMd5Salt 验证MD5和salt
func VerifyMd5Salt(pwd, slat, enPwd string) bool {
	return password.Verify(pwd, slat, enPwd, &PassOP)
}
