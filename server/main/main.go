package main

import (
	"GOchat/server/model"
	"GOchat/server/processor"
	"GOchat/server/utils"
	"fmt"
	"net"
)
func main() {
	utils.InitAll()
	model.InitSuperUser()
	fmt.Println("Server Listen in Port:8889")
	ln, err := net.Listen("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("Listen Error:", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("One Conn Error", err)
			continue
		}
		fmt.Println("One Conn Link. Addr:", conn.RemoteAddr().String())
		go processor.HandleConnection(conn)
	}
}
