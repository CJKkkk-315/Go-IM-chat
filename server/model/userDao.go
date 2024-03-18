package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type UserDao struct {
	DB *redis.Client
}
func (ud *UserDao) AddUser(u User) (err error) {
	_, err = ud.GetUser(u.UserId)
	if err != ErrorUserNotExist {
		return ErrorUserExist
	}

	userBytes, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		return
	}
	userString := string(userBytes)
	ud.DB.HSet(context.Background(), "Users", u.UserId, userString)
	return
}

func (ud *UserDao) GetUser(userId int) (u User, err error) {
	userString, err := ud.DB.HGet(context.Background(), "Users", strconv.Itoa(userId)).Bytes()
	if err != nil {
		return u, ErrorUserNotExist
	}
	err = json.Unmarshal(userString, &u)
	return
}