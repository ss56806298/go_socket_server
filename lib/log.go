package lib

import (
	"log"
	"os"
)

const errLogPath = "./log/error.log"
const serverLogPath = "./log/server.log"
const msgLogPath = "./log/msg.log"
const numLogPath = "./log/num.log"

//错误记录
func LogErr(v ...interface{}) {
	logfile := os.Stdout
	log.Println(v...)
	logger := log.New(logfile, "\r\n", log.Llongfile|log.Ldate|log.Ltime)
	logger.SetPrefix("[Error]")
	logger.Println(v...)

	logfile2, err := os.OpenFile(errLogPath, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		os.Exit(-1)
	}
	log.SetOutput(logfile2)
	log.SetPrefix("[Error]")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println(v...)

	defer logfile.Close()
	defer logfile2.Close()
}

//运行记录
func Log(v ...interface{}) {
	//后台输出
	logfile := os.Stdout
	log.Println(v...)
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime)
	logger.SetPrefix("[Info]")
	logger.Println(v...)

	//LOG输出
	logfile2, err := os.OpenFile(serverLogPath, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		os.Exit(-1)
	}
	log.SetOutput(logfile2)
	log.SetPrefix("[Info]")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println(v...)

	defer logfile.Close()
	defer logfile2.Close()
}

//全服广播记录
func LogMsg(v ...interface{}) {
	//LOG输出
	logfile, err := os.OpenFile(msgLogPath, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		os.Exit(-1)
	}
	log.SetOutput(logfile)
	log.SetPrefix("[Msg]")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println(v...)

	defer logfile.Close()
}

//全服人数记录
func LogNum(v ...interface{}) {
	//LOG输出
	logfile, err := os.OpenFile(numLogPath, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		os.Exit(-1)
	}
	log.SetOutput(logfile)
	log.SetPrefix("[Num]")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println(v...)

	defer logfile.Close()
}

//检查错误
func CheckError(err error) {
	if err != nil {
		LogErr(os.Stderr, "Fatal error: %s", err.Error())
	}
}
