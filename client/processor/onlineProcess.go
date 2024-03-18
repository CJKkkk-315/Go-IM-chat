package processor

import (
	"GOchat/client/common"
	"GOchat/client/model"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

type OnlineProcessor struct {
	Conn net.Conn
	User model.User
	Flag bool
}

var Op *OnlineProcessor
func (op *OnlineProcessor) HeartCheck() {
	for op.Flag {
		time.Sleep(time.Second)
		tf := common.Transfer{
			Conn: op.Conn,
			Buf:  [8092]byte{},
		}

		heartMes := common.Message{
			Type: common.HeartCheckType,
			Data: strconv.Itoa(op.User.UserId),
		}
		_ = tf.WritePkg(heartMes)
	}

}

func Online(up *UserProcessor) {
	Op = &OnlineProcessor{
		Conn:   up.Conn,
		User:   up.User,
		Flag:   true,
	}

}
func (op *OnlineProcessor) Keep() {
	tf := common.Transfer{
		Conn: op.Conn,
		Buf:  [8092]byte{},
	}
	for op.Flag {
		message, err := tf.ReadPkg()
		if err != nil {
			return
		}
		switch message.Type {
		case common.OnlineStatusResType :
			op.OnlineUserPrint(message)
		case common.GroupMessageType:
			op.ShowGroupMessage(message)

		default:
			fmt.Println("No Exist MessageType", message)

		}
	}
}

func (op *OnlineProcessor) ShowGroupMessage(message common.Message) {
	sms := common.ShortMessage{}
	_ = json.Unmarshal([]byte(message.Data), &sms)
	fmt.Println(sms.SendUser.UserName, "对大家说：", sms.Content)
}

func (op *OnlineProcessor) SendGroupMessage(content string) {
	sms := common.ShortMessage{
		SendUser: op.User,
		Content:  content,
	}
	tf := common.Transfer{
		Conn: op.Conn,
		Buf:  [8092]byte{},
	}
	smsStr, _ := json.Marshal(sms)
	mes := common.Message{
		Type: common.GroupMessageType,
		Data: string(smsStr),
	}
	_ = tf.WritePkg(mes)
}

func (op *OnlineProcessor) OnlineUserPrint(message common.Message) {
	//onlineList := make([]model.User, 0)
	onlineStatusRes := common.OnlineStatusRes{}

	_ = json.Unmarshal([]byte(message.Data), &onlineStatusRes)
	fmt.Println("当前在线列表:")
	for _,v := range onlineStatusRes.UsersList {
		fmt.Println("用户id:", v.UserId, "用户名:", v.UserName)
	}
}

func (op *OnlineProcessor) GetOnlineUser() {
	tf := common.Transfer{
		Conn: op.Conn,
	}
	osMes := common.Message{
		Type: common.OnlineStatusType,
		Data: "",
	}
	_ = tf.WritePkg(osMes)
}

func (op *OnlineProcessor) NoticeOnline() {
	tf := common.Transfer{
		Conn: op.Conn,
	}
	noticeMes := common.OnlineNotice{User:op.User}
	data, _ := json.Marshal(noticeMes)
	osMes := common.Message{
		Type: common.OnlineNoticeType,
		Data: string(data),
	}
	_ = tf.WritePkg(osMes)
}