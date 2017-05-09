package main

import (
	"./lib"
	"net"
	"strconv"
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
	buffer := make([]byte, 300)

	messager := make(chan byte)

	//创建
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			lib.Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		lib.Log("receive data:", string(buffer))
		tmpBuffer = lib.Depack(append(tmpBuffer, buffer[:n]...))
		lib.Log("receive data string:", string(tmpBuffer))
		//分发任务
		lib.TaskDeliver(tmpBuffer, conn)

		//开始心跳
		go lib.HeartBeating(conn, messager, timeout)
		//是否收到客户端的消息
		go lib.GravelChannel(tmpBuffer, messager)
	}
	//断开连接
	defer lib.Disconnect(conn)
}
