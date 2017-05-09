package lib

import (
	"bytes"
	"encoding/binary"
)

/*
   一个简单的通讯协议，由 header + 信息长度 ＋ 信息内容组成
*/

const (
	ConstHeader       = "wallHeader"
	ConstHeaderLength = 10
	ConstMLength      = 4
)

func Enpack(message []byte) []byte {
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}

func Depack(message []byte) []byte {
	length := len(message)

	Log("receive data length:", length)

	var i int
	data := make([]byte, 300)
	for i = 0; i < length; i = i + 1 {
		//超出长度不再循环
		if length < i+ConstHeaderLength+ConstMLength {
			break
		}
		if string(message[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(message[i+ConstHeaderLength : i+ConstHeaderLength+ConstMLength])

			Log("receive data message length:", messageLength)
			if length < i+ConstHeaderLength+ConstMLength+messageLength {
				break
			}
			data = message[i+ConstHeaderLength+ConstMLength : i+ConstHeaderLength+ConstMLength+messageLength]
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return data
}

//字节转为INT
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)

}

//INT转为字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
