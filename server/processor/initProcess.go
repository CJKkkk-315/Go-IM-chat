package processor

import (
	"GOchat/server/common"
	"fmt"
	"net"
)
type linkProcessor struct {
	Conn net.Conn
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	lp := &linkProcessor{Conn:conn}
	for {
		tf := &common.Transfer{
			Conn: conn,
		}

		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("Read Pkg Error", err.Error())
			return
		}
		lp.serverProcessMessage(&mes)
	}
}
func (lp *linkProcessor) serverProcessMessage(message *common.Message) {
	//fmt.Println(message.Type)
	switch message.Type {
	case common.LoginMessageType :
		up := userProcessor{Conn:lp.Conn}
		up.ServerProcessLogin(message)
	case common.RegisterMessageType:
		up := userProcessor{Conn:lp.Conn}
		up.ServerProcessRegister(message)
	case common.LogOutType:
		up := userProcessor{Conn:lp.Conn}
		up.ServerProcessLogOut(message)
	case common.OnlineStatusType:
		sp := smsProcessor{Conn:lp.Conn}
		sp.ShowOnlineList(message)
	case common.HeartCheckType:
		up := userProcessor{Conn:lp.Conn}
		up.ServerProcessHeart(message)
	case common.GroupMessageType:
		sp := smsProcessor{Conn:lp.Conn}
		sp.GroupMessage(message)
	case common.OnlineNoticeType:
		sp := smsProcessor{Conn:lp.Conn}
		sp.NoticeOnline(message)
	default:
		fmt.Println("No Exist MessageType")

	}
}