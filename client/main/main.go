package main

import (
	"GOchat/client/processor"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("欢迎来到Go聊天系统")
		fmt.Println("请选择操作：")
		fmt.Println("1: 登录")
		fmt.Println("2: 注册")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		conn, err := net.Dial("tcp", "127.0.0.1:8889")
		if err != nil {
			fmt.Println(err)
		}
		up := processor.UserProcessor{Conn: conn}
		switch choice {
		case "1":
			fmt.Println("您选择了登录")
			fmt.Print("请输入用户ID: ")
			userName, _ := reader.ReadString('\n')
			userName = strings.TrimSpace(userName)
			userId, _ := strconv.Atoi(userName)
			fmt.Print("请输入密码: ")
			userPwd, _ := reader.ReadString('\n')
			userPwd = strings.TrimSpace(userPwd)
			err := up.Login(userId, userPwd)
			if err != nil {
				fmt.Println("登录失败！", err)
				break
			}
			fmt.Println("登录成功！")



			up.MainInterface()
		case "2":
			fmt.Println("您选择了注册")
			fmt.Print("请设置用户名: ")
			userName, _ := reader.ReadString('\n')
			userName = strings.TrimSpace(userName)
			fmt.Print("请设置用户ID: ")
			userIdStr, _ := reader.ReadString('\n')
			userIdStr = strings.TrimSpace(userIdStr)
			userId, _ := strconv.Atoi(userIdStr)

			fmt.Print("请设置密码: ")
			userPwd, _ := reader.ReadString('\n')
			userPwd = strings.TrimSpace(userPwd)

			err := up.Register(userId, userPwd, userName)
			if err != nil {
				fmt.Println("注册失败！", err)
				break
			}
			fmt.Println("注册成功！")

		default:
			fmt.Println("未知选择，请重新运行程序")
			return
		}
	}
}




