package lib

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net"
	"strconv"
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

//服务器任务分发
func ServerTaskDeliver(postdata []byte, conn net.Conn) {
	//解析Json
	js, js_err := simplejson.NewJson(postdata)
	if js_err != nil {
		Log(js_err)
		return
	}
	num, _ := js.Get("num").Int()
	switch num {
	//跑马灯
	case 1001:
		message, _ := js.Get("message").String()
		SendMessageToAll(message)
		break

	//获取公会在线玩家ID
	case 1002:
		var girudo_member_ids []interface{}
		girudo_member_ids, _ = js.Get("girudo_member_ids").Array()

		online_member_ids := make([]int, 0)
		for _, member_id := range girudo_member_ids {
			x, _ := member_id.(json.Number)
			y, _ := strconv.ParseFloat(string(x), 64)

			member_id_int := int(y)

			if checkOnline(member_id_int) {
				online_member_ids = append(online_member_ids, member_id_int)
			}
		}

		js_s := simplejson.New()
		js_s.Set("online_member_ids", online_member_ids)

		byte_message, _ := js_s.Encode()
		conn.Write(byte_message)
		break

	//公会相关的信息
	case 1003:
		var girudo_member_ids []interface{}
		girudo_member_ids, _ = js.Get("girudo_member_ids").Array()
		id, _ := js.Get("id").Int()
		nickname, _ := js.Get("nickname").String()
		type_g, _ := js.Get("type").Int()
		create_time, _ := js.Get("create_time").Int()

		//json打包推送信息
		js_s_client := simplejson.New()

		js_s_client.Set("code", 2)
		js_s_client.Set("id", id)
		js_s_client.Set("nickname", nickname)
		js_s_client.Set("type", type_g)
		js_s_client.Set("create_time", create_time)

		switch type_g {
		case 1:
			message, _ := js.Get("message").String()
			girudo_class, _ := js.Get("girudo_class").Int()
			user_id, _ := js.Get("user_id").Int()

			js_s_client.Set("message", message)
			js_s_client.Set("girudo_class", girudo_class)
			js_s_client.Set("user_id", user_id)
			break
		case 2:
		case 3:
		case 4:
		case 5:
			break
		case 6:
			dungeon_id, _ := js.Get("dungeon_id").String()

			js_s_client.Set("dungeon_id", dungeon_id)
			break
		}
		byte_message, _ := js_s_client.Encode()

		for _, member_id := range girudo_member_ids {
			x, _ := member_id.(json.Number)
			y, _ := strconv.ParseFloat(string(x), 64)

			member_id_int := int(y)

			if checkOnline(member_id_int) {
				member_conn := getOnlineUserConn(member_id_int)
				member_conn.Write(Enpack(byte_message))
			}
		}
		break

	//全服聊天
	case 1004:
		message, _ := js.Get("message").String()
		nickname, _ := js.Get("nickname").String()
		user_id, _ := js.Get("user_id").Int()

		//json打包推送信息
		js_s_client := simplejson.New()

		js_s_client.Set("code", 3)
		js_s_client.Set("message", message)
		js_s_client.Set("nickname", nickname)
		js_s_client.Set("user_id", user_id)

		byte_message, _ := js_s_client.Encode()

		SendToAll(byte_message)
		break

	default:
		break
	}
}

//任务分发
func TaskDeliver(postdata []byte, conn net.Conn) {
	// for _, v := range routers {
	// pred := v[0]
	// act := v[1]
	// var entermsg Msg
	// err := json.Unmarshal(postdata, &entermsg)
	//解析Json
	js, js_err := simplejson.NewJson(postdata)
	if js_err != nil {
		Log(js_err)
		return
	}
	// if pred.(func(entermsg Msg) bool)(entermsg) {
	// 	act.(Controller).Excute(entermsg, conn)
	// 	return
	// }
	code, code_js_err := js.Get("code").Int()

	if code_js_err != nil {
		Log(code_js_err)
		return
	}

	switch code {
	//公会聊天的信息
	case 2:
		girudo_member_ids, _ := js.Get("girudo_member_ids").Array()
		message, _ := js.Get("message").String()

		//json打包推送信息
		js_s_client := simplejson.New()

		js_s_client.Set("code", 2)
		js_s_client.Set("message", message)

		byte_message, _ := js_s_client.Encode()

		for _, member_id := range girudo_member_ids {
			memer_id_int := member_id.(int)

			if checkOnline(memer_id_int) {
				member_conn := getOnlineUserConn(memer_id_int)

				member_conn.Write(byte_message)
			}
		}
		break
	default:
		conn.Write([]byte("error"))
		break
	}
	// }
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
