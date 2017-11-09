package lib

import (
	"encoding/json"
	"net"
)

/*
	连接池
*/

var connpool = make(map[int]net.Conn)

//加入连接池
func AddToPool(uuid int, conn net.Conn) {
	connpool[uuid] = conn
	LogNum("online players:", len(connpool))
}

//剔除连接池
func DeleteFromPool(uuid int, conn net.Conn) {
	delete(connpool, uuid)
	LogNum("online players:", len(connpool))
}

//检查某个ID的玩家是否在线
func checkOnline(uuid int) bool {
	if _, ok := connpool[uuid]; ok {
		return true
	}
	return false
}

//获取某个ID玩家的连接
func getOnlineUserConn(uuid int) net.Conn {
	return connpool[uuid]
}

//看下链接池的连接数
func CountPool() int {
	Log(len(connpool))
	return len(connpool)
}

//全部发送消息
func SendMessageToAll(message string) {
	var data Announcement

	data.Code = 1
	data.Message = message

	body, err := json.Marshal(data)
	Log("send message to All:", string(body))
	if err != nil {
		LogErr(err)
	}

	for _, conn := range connpool {
		conn.Write(Enpack(body))
	}
}

func SendToAll(message []byte) {
	for _, conn := range connpool {
		conn.Write(Enpack(message))
	}
}

//跑马灯的结构体
type Announcement struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
