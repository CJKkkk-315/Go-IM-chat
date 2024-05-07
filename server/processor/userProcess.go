package processor

import (
	"GOchat/server/common"
	"GOchat/server/model"
	"GOchat/server/utils"
	"fmt"
	"net"
	"strconv"
	"time"
)

type UserProcessor struct {
	Conn net.Conn
}

func (up UserProcessor) ServerProcessLogin(loginMes *common.LoginMessage, loginResMes *common.LoginResMessage) (err error) {

	u, err := model.Login(loginMes.UserId, loginMes.UserPwd)
	if err == nil {
		fmt.Println(u.UserName, " Login Success!")
		loginResMes.Code = 200
		loginResMes.User = u
		loginResMes.Error = ""
	} else {
		loginResMes.Code = 500
		loginResMes.User = u
		loginResMes.Error = err.Error()
	}

	utils.OnlineMap[loginMes.UserId] = up.Conn
	utils.HeartMap[loginMes.UserId] = time.Now()
	return
}

func (up UserProcessor) ServerProcessLogOut(message *common.Message, reply *common.Message) (err error) {
	userId, _ := strconv.Atoi(message.Data)
	utils.HeartLock.Lock()
	delete(utils.OnlineMap, userId)
	delete(utils.HeartMap, userId)
	utils.HeartLock.Unlock()
	return
}

func (up UserProcessor) ServerProcessHeart(message *common.Message, reply *common.Message) (err error) {
	userId, _ := strconv.Atoi(message.Data)
	utils.HeartLock.Lock()
	utils.HeartMap[userId] = time.Now()
	utils.HeartLock.Unlock()
	return
}

func (up UserProcessor) ServerProcessRegister(registerMes *common.RegisterMessage, registerResMes *common.RegisterResMessage) (err error) {
	err = model.Register(registerMes.UserId, registerMes.UserPwd, registerMes.UserName)
	if err == nil {
		registerResMes.Code = 200
		registerResMes.Error = ""
	} else {
		registerResMes.Code = 500
		registerResMes.Error = err.Error()
	}
	return
}
