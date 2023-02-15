package utils

import (
	"crypto/md5"
	"github.com/anaskhan96/go-password-encoder"
)

var PassOP password.Options

func init() {
	PassOP = password.Options{SaltLen: 10, Iterations: 10000, KeyLen: 50, HashFunction: md5.New}
}

func EncodeMd5Salt(pwd string) (string, string) {
	return password.Encode(pwd, &PassOP)
}

func VerifyMd5Salt(pwd, slat, enPwd string) bool {
	return password.Verify(pwd, slat, enPwd, &PassOP)
}
