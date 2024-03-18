package model

import (
	"GOchat/server/utils"
)

type User struct {
	UserId int `json:"user_id"`
	UserPwd string `json:"user_pwd"`
	UserName string `json:"user_name"`
}

func Login(userId int, userPwd string) (User, error) {
	ud := UserDao{
		DB:	utils.RB,
	}
	u, err := ud.GetUser(userId)
	if err != nil {
		return u, err
	}
	if u.UserPwd != userPwd {
		return u, ErrorPassword
	}
	return u, nil
}

func Register(userId int, userPwd, userName string) error {
	ud := UserDao{
		DB:	utils.RB,
	}
	err := ud.AddUser(User{
		UserId:   userId,
		UserPwd:  userPwd,
		UserName: userName,
	})

	return err
}
