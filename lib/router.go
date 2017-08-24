package lib

import (
	"encoding/json"
	"fmt"
	"net"
)

//消息接口
type Msg struct {
	Code    int                    `json:"code"`
	Content map[string]interface{} `json:"content"`
}

//控制器接口
type Controller interface {
	Excute(message Msg, conn net.Conn)
}

//路线接口
var routers [][2]interface{}

//按路线执行程序
func Route(pred interface{}, controller Controller) {
	switch pred.(type) {
	case func(entry Msg) bool:
		{
			var arr [2]interface{}
			arr[0] = pred
			arr[1] = controller
			routers = append(routers, arr)
		}
	// case map[string]interface{}:
	// 	{
	// 		defaultPred := func(entry Msg) bool {
	// 			for keyPred, valPred := range pred.(map[string]interface{}) {
	// 				val, ok := entry.Code
	// 				if !ok {
	// 					return false
	// 				}
	// 				if val != valPred {
	// 					return false
	// 				}
	// 			}
	// 			return true
	// 		}
	// 		var arr [2]interface{}
	// 		arr[0] = defaultPred
	// 		arr[1] = controller
	// 		routers = append(routers, arr)
	// 		fmt.Println(routers)
	// 	}
	default:
		fmt.Println("requested controller not found")
	}
}

//任务分发
func TaskDeliver(postdata []byte, conn net.Conn) {
	for _, v := range routers {
		pred := v[0]
		act := v[1]
		var entermsg Msg
		err := json.Unmarshal(postdata, &entermsg)
		if err != nil {
			Log(err)
		}
		if pred.(func(entermsg Msg) bool)(entermsg) {
			act.(Controller).Excute(entermsg, conn)
			return
		}
	}
}

type EchoController struct {
}

type RegisterController struct {
}

type PushController struct {
}

type PnumController struct {
}

/*

	任务分发到控制器上进行分别处理
	1.Register:注册任务
	2.Push:全服推送任务

*/
func (this *EchoController) Excute(message Msg, conn net.Conn) {
	mirrormsg, err := json.Marshal(message)
	Log("echo the message:", string(mirrormsg))
	CheckError(err)
}

func (this *RegisterController) Excute(message Msg, conn net.Conn) {
	//注册的玩家ID
	uuid := message.Content["user_id"]
	uid := int(uuid.(float64))
	AddToPool(uid, conn)
}

func (this *PushController) Excute(message Msg, conn net.Conn) {
	imessage := message.Content["message"]
	smessage := imessage.(string)
	SendMessageToAll(smessage)
	LogMsg(smessage)
}

func (this *PnumController) Excute(message Msg, conn net.Conn) {
	num := CountPool()
	conn.Write(IntToBytes(num))
}

func init() {
	//最大20个控制器
	routers = make([][2]interface{}, 0, 20)

	//输出控制 code:100
	var echo EchoController
	Route(func(entry Msg) bool {
		if entry.Code == 100 {
			return true
		}
		return false
	}, &echo)

	//code:1 注册控制,加入链接池
	var register RegisterController
	Route(func(entry Msg) bool {
		if entry.Code == 1 {
			return true
		}
		return false
	}, &register)

	/*
		PHP服务器传过来的推送信息
	*/
	//code:1001
	var push PushController
	Route(func(entry Msg) bool {
		if entry.Code == 1001 {
			return true
		}
		return false
	}, &push)

	/*
		查看在线人数
		code:1002
	*/
	var pnum PnumController
	Route(func(entry Msg) bool {
		if entry.Code == 1002 {
			return true
		}
		return false
	}, &pnum)
}
