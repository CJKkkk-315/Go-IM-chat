package common

import "GOchat/client/model"

const (
	LoginMessageType = "LoginMessage"
	LoginResMessageType = "LoginResMessage"
	RegisterMessageType = "RegisterMessage"
	RegisterResMessageType = "RegisterResMessage"
	GroupMessageType = "GroupMessage"
	PrivateMessageType = "PrivateMessage"
	OnlineNoticeType = "OnlineNotice"
	OnlineStatusType = "OnlineStatus"
	OnlineStatusResType = "OnlineStatusRes"
	LogOutType = "LogOut"
	HeartCheckType = "HeartCheck"
)
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMessage struct {
	UserId int `json:"user_id"`
	UserPwd string `json:"user_pwd"`
}

type LoginResMessage struct {
	Code int `json:"code"`// 500:未注册 200:成功
	User model.User `json:"user"`
	Error string `json:"error"`
}

type RegisterMessage struct {
	UserId int `json:"user_id"`
	UserPwd string `json:"user_pwd"`
	UserName string `json:"user_name"`
}

type RegisterResMessage struct {
	Code int `json:"code"`// 500:未注册 200:成功
	Error string `json:"error"`
}

type OnlineStatusRes struct {
	UsersList []model.User `json:"users_list"`
}

type ShortMessage struct {
	SendUser model.User `json:"send_user"`
	Content string `json:"content"`
}

type OnlineNotice struct {
	User model.User `json:"user"`
}