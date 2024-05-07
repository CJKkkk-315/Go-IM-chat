package processor

import (
	"GOchat/client/common"
	"GOchat/client/model"
	KgoRpc "GOchat/krpc"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

type OnlineProcessor struct {
	RpcClient *KgoRpc.Client
	Conn      net.Conn
	User      model.User
	Flag      bool
}

var Op *OnlineProcessor

func (op *OnlineProcessor) HeartCheck() {
	for op.Flag {
		time.Sleep(time.Second)

		heartMes := common.Message{
			Type: common.HeartCheckType,
			Data: strconv.Itoa(op.User.UserId),
		}
		err := op.RpcClient.Call(context.Background(), "UserProcessor.ServerProcessHeart", &heartMes, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

	}

}

func Online(up *UserProcessor) {
	Op = &OnlineProcessor{
		RpcClient: up.RpcClient,
		Conn:      up.Conn,
		User:      up.User,
		Flag:      true,
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
		case common.OnlineStatusResType:
			//op.OnlineUserPrint(message)
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
	smsStr, _ := json.Marshal(sms)
	mes := common.Message{
		Type: common.GroupMessageType,
		Data: string(smsStr),
	}
	err := op.RpcClient.Call(context.Background(), "SmsProcessor.GroupMessage", &mes, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (op *OnlineProcessor) OnlineUserPrint(message common.Message) {
	//onlineList := make([]model.User, 0)
	onlineStatusRes := common.OnlineStatusRes{}

	_ = json.Unmarshal([]byte(message.Data), &onlineStatusRes)
	fmt.Println("当前在线列表:")
	for _, v := range onlineStatusRes.UsersList {
		fmt.Println("用户id:", v.UserId, "用户名:", v.UserName)
	}
}

func (op *OnlineProcessor) GetOnlineUser() {
	reply := common.OnlineStatusRes{}
	err := op.RpcClient.Call(context.Background(), "SmsProcessor.ShowOnlineList", &common.Message{}, &reply)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range reply.UsersList {
		fmt.Println("用户id:", v.UserId, "用户名:", v.UserName)
	}

}

func (op *OnlineProcessor) NoticeOnline() {
	noticeMes := common.OnlineNotice{User: op.User}
	err := op.RpcClient.Call(context.Background(), "SmsProcessor.NoticeOnline", &noticeMes, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}
