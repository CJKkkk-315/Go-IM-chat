package processor

import (
	"GOchat/client/common"
	"GOchat/client/model"
	KgoRpc "GOchat/krpc"
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type UserProcessor struct {
	RpcClient *KgoRpc.Client
	Conn      net.Conn
	User      model.User
}

func (up *UserProcessor) Login(userId int, userPwd string) error {

	loginMes := common.LoginMessage{
		UserId:  userId,
		UserPwd: userPwd,
	}
	mesRes := common.LoginResMessage{}

	err := up.RpcClient.Call(context.Background(), "UserProcessor.ServerProcessLogin", &loginMes, &mesRes)
	if err != nil {
		fmt.Println(err)
	}
	if mesRes.Code == 200 {
		up.User = mesRes.User
		Online(up)
		//Op.NoticeOnline()
		go Op.HeartCheck()
		//go Op.Keep()
		return nil
	} else {
		return errors.New(mesRes.Error)
	}
}

func (up *UserProcessor) Register(userId int, userPwd, userName string) error {
	RegisterMes := common.RegisterMessage{
		UserId:   userId,
		UserPwd:  userPwd,
		UserName: userName,
	}
	mesRes := common.LoginResMessage{}
	err := up.RpcClient.Call(context.Background(), "UserProcessor.ServerProcessRegister", &RegisterMes, &mesRes)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if mesRes.Code == 200 {
		return nil
	} else {
		return errors.New(mesRes.Error)
	}
}
func (up *UserProcessor) LogOut() {
	logOutMes := common.Message{
		Type: common.LogOutType,
		Data: strconv.Itoa(up.User.UserId),
	}
	err := up.RpcClient.Call(context.Background(), "UserProcessor.ServerProcessLogOut", &logOutMes, nil)

	if err != nil {
		fmt.Println(err)
	}
	_ = up.Conn.Close()
	Op.Flag = false
	Op = nil
}
func showMenu() {
	fmt.Println("请选择操作：")
	fmt.Println("1: 获取在线人数")
	fmt.Println("2: 群发消息")
	fmt.Println("3: 私发消息")
	fmt.Println("4: 退出登录")
}
func (up *UserProcessor) MainInterface() {
	reader := bufio.NewReader(os.Stdin)
	showMenu()
	for {
		operation, _ := reader.ReadString('\n')
		operation = strings.TrimSpace(operation)

		fmt.Println(operation)
		switch operation {
		case "1":
			fmt.Println("您选择了获取在线人数")
			Op.GetOnlineUser()
		case "2":
			fmt.Println("您选择了群发消息")
			fmt.Print("请输入群发消息内容: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message) // 去除输入字符串的前后空格
			Op.SendGroupMessage(message)

		case "3":
			fmt.Println("您选择了私发消息")
			fmt.Print("请输入收件人用户名: ")
			recipient, _ := reader.ReadString('\n')
			recipient = strings.TrimSpace(recipient) // 去除输入字符串的前后空格

			fmt.Print("请输入消息内容: ")
			privateMessage, _ := reader.ReadString('\n')
			privateMessage = strings.TrimSpace(privateMessage) // 去除输入字符串的前后空格
			// 实际应用中，在这里添加私发消息的代码
			fmt.Printf("收件人：%s，消息内容：%s\n", recipient, privateMessage)
		case "4":
			up.LogOut()
			return
		default:
			fmt.Println("输入有误！")
			showMenu()
		}
	}

}
