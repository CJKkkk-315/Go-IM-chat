package processor

import (
	"GOchat/server/common"
	"GOchat/server/model"
	"GOchat/server/utils"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type SmsProcessor struct {
	Conn net.Conn
}

func (sp *SmsProcessor) ShowOnlineList(message *common.Message, reply *common.OnlineStatusRes) (err error) {
	onlineKeys := make([]string, 0, len(utils.OnlineMap))
	for k := range utils.OnlineMap {
		onlineKeys = append(onlineKeys, strconv.Itoa(k))
	}

	userStrs, _ := utils.RB.HMGet(context.Background(), "Users", onlineKeys...).Result()

	usersList := make([]model.User, 0, len(userStrs))
	for _, v := range userStrs {
		u := model.User{}
		fmt.Println(v.(string))
		err := json.Unmarshal([]byte(v.(string)), &u)
		if err != nil {
			fmt.Println(err)
		}
		u.UserPwd = ""
		usersList = append(usersList, u)
	}
	fmt.Println(usersList)
	reply.UsersList = usersList
	return

}

func (sp *SmsProcessor) GroupMessage(message *common.Message, reply *common.Message) (err error) {
	groupMes := common.ShortMessage{}
	_ = json.Unmarshal([]byte(message.Data), &groupMes)
	sendId := groupMes.SendUser.UserId
	for k, v := range utils.OnlineMap {
		if k == sendId {
			continue
		}
		tf := common.Transfer{
			Conn: v,
			Buf:  [8092]byte{},
		}
		_ = tf.WritePkg(*message)
	}
	return
}

func (sp *SmsProcessor) NoticeOnline(noticeMes *common.OnlineNotice, reply *common.Message) (err error) {
	sendUser := noticeMes.User
	for k, v := range utils.OnlineMap {
		if k == sendUser.UserId {
			continue
		}
		tf := common.Transfer{
			Conn: v,
			Buf:  [8092]byte{},
		}
		content := "用户id:" + strconv.Itoa(sendUser.UserId) + " " + sendUser.UserName + "上线了！"
		sms := common.ShortMessage{
			SendUser: model.SuperUser,
			Content:  content,
		}
		smsStr, _ := json.Marshal(sms)
		mes := common.Message{
			Type: common.GroupMessageType,
			Data: string(smsStr),
		}
		_ = tf.WritePkg(mes)
	}
	return
}
