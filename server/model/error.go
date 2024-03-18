package model

import "errors"

var (
	ErrorUserNotExist = errors.New("user not exist")
	ErrorUserExist = errors.New("user exist")
	ErrorPassword = errors.New("password error")
)
