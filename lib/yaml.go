package lib

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//获取配置文件的map集合
func GetYamlConfig(path string) map[interface{}]interface{} {
	//读取文件的二进制流
	data, err := ioutil.ReadFile(path)
	m := make(map[interface{}]interface{})
	if err != nil {
		LogErr(err)
	}
	//转换为集合
	err = yaml.Unmarshal([]byte(data), &m)
	return m
}

//将集合中的value根据key值提取出来
func GetElement(key string, themap map[interface{}]interface{}) string {
	if value, ok := themap[key]; ok {
		return fmt.Sprint(value)
	}

	Log("can not find the config file")
	return ""
}
