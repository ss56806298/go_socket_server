package main

import (
	"./lib"
	"github.com/bitly/go-simplejson"
	"net"
	"strconv"
	"time"
)

func main() {
	startServer("./conf/config.yaml")
}

func startServer(configpath string) {
	//读取本地配置文件
	configmap := lib.GetYamlConfig(configpath)
	host := lib.GetElement("host", configmap)

	//本地创建套接字并监听端口
	//心跳间隔
	timeinterval, err := strconv.Atoi(lib.GetElement("beatinginterval", configmap))
	lib.CheckError(err)
	netListen, err := net.Listen("tcp", host)
	lib.CheckError(err)
	defer netListen.Close()
	lib.Log("Server Starting, wait for clients")

	//for循环处理客户端请求
	for {
		conn, err := netListen.Accept()
		//错误记录并继续运行
		if err != nil {
			lib.CheckError(err)
			continue
		}

		lib.Log(conn.RemoteAddr().String(), " tcp connect success")

		//处理链接请求
		go handleConnection(conn, timeinterval)
	}
}

//处理链接请求
func handleConnection(conn net.Conn, timeout int) {
	//缓冲区数据
	tmpBuffer := make([]byte, 0)
	//客户端提交的数据限制在32字节以内
	buffer := make([]byte, 512)

	// messager := make(chan byte)

	//设置超时时间
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	//创建
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			lib.Log(conn.RemoteAddr().String(), " connection error: ", err)
			//客户端链接出了问题,断开
			conn.Close()
			break
		}
		lib.Log("receive data:", string(buffer))
		tmpBuffer = lib.Depack(append(tmpBuffer, buffer[:n]...))
		lib.Log("receive data string:", string(tmpBuffer))

		//解析JSON
		js, js_err := simplejson.NewJson(tmpBuffer)

		if js_err != nil {
			lib.Log(conn.RemoteAddr().String(), " json analyze error: ", err)
			conn.Write([]byte("fail"))
			continue
		}

		//server or client？
		code, code_err := js.Get("code").Int()
		_, num_err := js.Get("num").Int()
		if code_err == nil {
			user_id, _ := js.Get("user_id").Int()

			if code == 1 {
				//注册成功
				conn.Write(lib.Enpack([]byte("success")))
				go receiveMessage(conn, user_id, timeout)
				break
			} else {
				conn.Write(lib.Enpack([]byte("register first")))
				continue
			}
		} else if num_err == nil {
			lib.ServerTaskDeliver(tmpBuffer, conn)
			conn.Close()
		}

		// //分发任务
		// lib.TaskDeliver(tmpBuffer, conn)

		// //开始心跳
		// go lib.HeartBeating(conn, messager, timeout)
		// //是否收到客户端的消息
		// go lib.GravelChannel(tmpBuffer, messager)
	}
}

//加入到收取消息的队列中
func receiveMessage(conn net.Conn, user_id int, timeout int) {
	//缓冲区数据
	tmpBuffer := make([]byte, 0)
	//客户端提交的数据限制在32字节以内
	buffer := make([]byte, 300)

	messager := make(chan byte)

	//设置超时时间
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	//加入连接池
	lib.AddToPool(user_id, conn)

	//创建
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			lib.Log(conn.RemoteAddr().String(), " connection error: ", err)
			//客户端链接出了问题,断开并从连接池从剔除
			conn.Close()
			lib.DeleteFromPool(user_id, conn)
			break
		}

		lib.Log(string(user_id), " send data:", string(buffer))
		tmpBuffer = lib.Depack(append(tmpBuffer, buffer[:n]...))
		lib.Log(string(user_id), " send data:", string(tmpBuffer))

		//分发任务
		lib.TaskDeliver(tmpBuffer, conn)

		//开始心跳
		go lib.HeartBeating(conn, messager, timeout)
		//是否收到客户端的消息
		go lib.GravelChannel(tmpBuffer, messager)
	}
}
