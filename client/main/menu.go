package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func mainInterface() {
	reader := bufio.NewReader(os.Stdin)

	// 二级菜单
	fmt.Println("请选择操作：")
	fmt.Println("1: 获取在线人数")
	fmt.Println("2: 群发消息")
	fmt.Println("3: 私发消息")

	operation, _ := reader.ReadString('\n')
	operation = strings.TrimSpace(operation) // 去除输入字符串的前后空格

	switch operation {
	case "1":
		fmt.Println("您选择了获取在线人数")
		// 实际应用中，在这里添加获取在线人数的代码

	case "2":
		fmt.Println("您选择了群发消息")
		fmt.Print("请输入群发消息内容: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message) // 去除输入字符串的前后空格
		// 实际应用中，在这里添加群发消息的代码
		fmt.Printf("群发消息内容是：%s\n", message)

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

	default:
		fmt.Println("未知操作，请重新运行程序")
		return
	}
}

