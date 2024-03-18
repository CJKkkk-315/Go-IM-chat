package processor

import (
	"GOchat/server/common"
	"GOchat/server/model"
	"GOchat/server/utils"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)
type userProcessor struct {
	Conn net.Conn
}

func (up *userProcessor) ServerProcessLogin(message *common.Message) {
	loginMes := common.LoginMessage{}
	err := json.Unmarshal([]byte(message.Data), &loginMes)
	if err != nil {
		fmt.Println("Decode Login Message Fail.", err)
	}
	var loginResMes common.LoginResMessage
	u, err := model.Login(loginMes.UserId, loginMes.UserPwd)
	if err == nil {
		fmt.Println(u.UserName, " Login Success!")
		loginResMes = common.LoginResMessage{
			Code:  200,
			User: u,
			Error: "",
		}
	} else {
		loginResMes = common.LoginResMessage{
			Code:  500,
			User: u,
			Error: err.Error(),
		}
	}
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("Encode Fail.", err)
	}
	resMessage := common.Message{
		Type: common.LoginResMessageType,
		Data: string(data),
	}
	tf := common.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(resMessage)
	if err != nil {
		fmt.Println("Write LoginResMessage Fail.", err)
	}
	utils.OnlineMap[loginMes.UserId] = up.Conn
	utils.HeartMap[loginMes.UserId] = time.Now()
}

func (up *userProcessor) ServerProcessLogOut(message *common.Message) {
	userId, _ := strconv.Atoi(message.Data)
	utils.HeartLock.Lock()
	delete(utils.OnlineMap, userId)
	delete(utils.HeartMap, userId)
	utils.HeartLock.Unlock()
}

func (up *userProcessor) ServerProcessHeart(message *common.Message) {
	userId, _ := strconv.Atoi(message.Data)
	utils.HeartLock.Lock()
	utils.HeartMap[userId] = time.Now()
	utils.HeartLock.Unlock()
}

func (up *userProcessor) ServerProcessRegister(message *common.Message) {
	registerMes := common.RegisterMessage{}
	err := json.Unmarshal([]byte(message.Data), &registerMes)
	if err != nil {
		fmt.Println("Decode Register Message Fail.", err)
	}
	var registerResMes common.RegisterResMessage
	err = model.Register(registerMes.UserId, registerMes.UserPwd, registerMes.UserName)
	if err == nil {
		registerResMes = common.RegisterResMessage{
			Code:  200,
			Error: "",
		}
	} else {
		registerResMes = common.RegisterResMessage{
			Code:  500,
			Error: err.Error(),
		}
	}
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("Encode Fail.", err)
	}
	resMessage := common.Message{
		Type: common.RegisterResMessageType,
		Data: string(data),
	}
	tf := common.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(resMessage)
	if err != nil {
		fmt.Println("Write RegisterResMessage Fail.", err)
	}
}