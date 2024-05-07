package main

import (
	KgoRpc "GOchat/krpc"
	"GOchat/server/model"
	"GOchat/server/processor"
	"GOchat/server/utils"
	"fmt"
	"net"
	"time"
)

func main() {
	utils.InitAll()
	model.InitSuperUser()
	fmt.Println("Server Listen in Port:8889")
	ln, err := net.Listen("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("Listen Error:", err)
	}
	var upIns processor.UserProcessor
	var spIns processor.SmsProcessor

	err = KgoRpc.Register(&upIns)
	if err != nil {
		fmt.Println(err)
	}

	err = KgoRpc.Register(&spIns)
	if err != nil {
		fmt.Println(err)
	}

	go KgoRpc.Accept(ln)

	for {
		time.Sleep(100 * time.Second)
	}
	//for {
	//	if err != nil {
	//		fmt.Println("One Conn Error", err)
	//		continue
	//	}
	//	fmt.Println("One Conn Link. Addr:", conn.RemoteAddr().String())
	//	go processor.HandleConnection(conn)
	//}
}
