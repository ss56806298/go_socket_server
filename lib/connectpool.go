package lib

import (
	"encoding/json"
	"net"
)

/*
	连接池
*/

var connpool = make(map[int]net.Conn)

//断开连接
func Disconnect(conn net.Conn) {
	conn.Close()
	DeleteFromPool(conn)
}

//加入连接池
func AddToPool(uuid int, conn net.Conn) {
	connpool[uuid] = conn
	LogNum("online players:", len(connpool))
}

//剔除连接池
func DeleteFromPool(conn net.Conn) {
	for key, value := range connpool {
		if value == conn {
			delete(connpool, key)
			LogNum("online players:", len(connpool))
			break
		}
	}
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

//跑马灯的结构体
type Announcement struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
