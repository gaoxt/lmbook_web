package config

import (
	"errors"
)

const (
	Succ              = 0
	ErrServiceBusy    = 1001
	ErrBookId         = 2001
	ErrBookIdOverFlow = 2002
)

var ErrMsg = map[int]string{
	Succ:              "",
	ErrServiceBusy:    "服务器错误",
	ErrBookId:         "书籍ID错误",
	ErrBookIdOverFlow: "书籍ID溢出",
}

func GetErrMsg(code int) error {
	return errors.New(ErrMsg[code])
}
