package utils

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const passwordLength = 256

var ERROR_ILLEGAL_PASSWORD = errors.New("Illegal password")

type password [passwordLength]byte

func init() {
	// 更新随机种子, 防止生成一样的随机密码
	rand.Seed(time.Now().Unix())
}

// 采用 base64 编码把密码转换为字符串
func (password *password) String() string {
	return base64.StdEncoding.EncodeToString(password[:])
}

// 解析采用 base64 编码的字符串获取密码
func ParsePassword(passwordString string) (*password, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(passwordString))
	if err != nil || len(bs) != passwordLength {
		return nil, ERROR_ILLEGAL_PASSWORD
	}
	psword := password{}
	copy(psword[:], bs)
	bs = nil
	return &psword, nil
}

// 产生 256个byte随机组合的 密码，最后会使用base64编码为字符串存储在配置文件中
// 不能出现任何一个重复的byte位，必须又 0-255 组成，并且都需要包含
func RandPassword() string {
	// 随机生成一个由 0 ~ 255 组成的 byte 数组
	intArr := rand.Perm(passwordLength)
	psword := &password{}
	for i, v := range intArr {
		psword[i] = byte(v)
		if i == v {
			// 确保不好出现任何一个byte位出现重复
			return RandPassword()
		}
	}
	return psword.String()
}
