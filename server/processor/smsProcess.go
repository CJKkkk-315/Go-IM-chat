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

type smsProcessor struct {
	Conn net.Conn
}

func (sp *smsProcessor) ShowOnlineList(message *common.Message) {
	tf := common.Transfer{
		Conn: sp.Conn,
		Buf:  [8092]byte{},
	}
	onlineKeys := make([]string, 0, len(utils.OnlineMap))
	for k := range utils.OnlineMap {
		onlineKeys = append(onlineKeys, strconv.Itoa(k))
	}

	userStrs, _ := utils.RB.HMGet(context.Background(), "Users", onlineKeys...).Result()

	usersList := make([]model.User, 0, len(userStrs))
	for _,v := range userStrs {
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
	onlineStatusRes := common.OnlineStatusRes{UsersList:usersList}
	data, _ := json.Marshal(onlineStatusRes)
	mes := common.Message{
		Type: common.OnlineStatusResType,
		Data: string(data),
	}
	_ = tf.WritePkg(mes)
}

func (sp *smsProcessor) GroupMessage(message *common.Message) {
	groupMes := common.ShortMessage{}
	_ = json.Unmarshal([]byte(message.Data), &groupMes)
	sendId := groupMes.SendUser.UserId
	for k,v := range utils.OnlineMap {
		if k == sendId {continue}
		tf := common.Transfer{
			Conn: v,
			Buf:  [8092]byte{},
		}
		_ = tf.WritePkg(*message)
	}
}

func (sp *smsProcessor) NoticeOnline(message *common.Message) {
	noticeMes := common.OnlineNotice{}
	_ = json.Unmarshal([]byte(message.Data), &noticeMes)
	sendUser := noticeMes.User
	for k,v := range utils.OnlineMap {
		if k == sendUser.UserId {continue}
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

}
