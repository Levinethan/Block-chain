package gob_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"testing"
)
//以下实现序列化 和反序列化
/*
gob包  是golang自带的一个数据结构序列化的编码和解码工具
Decoder   Encoder





 */
type Person struct {    //定义一个结构 使用gob编码 得到字节流 然后解码的功能
	Name string
	age  uint
}

func Test_funGob(t *testing.T)  {
	var xiaoMing  Person
	xiaoMing.Name = "小明"
	xiaoMing.age =  20
	var buffer  bytes.Buffer
	//使用一个辅助变量buffer  编码好的数据放进  buffer
	//分别定义一个  编码器  和解码器
	//1.编码器
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&xiaoMing)

	if err != nil {
		log.Panic("编码出错，小明失踪")
	}
	fmt.Printf("编码成功 : %v \n",buffer.Bytes())
	//使用解码器解码
	//接受一个Reader（）那就New 一个Reader  然后里面放buffer数据
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	//那么数据在解码器中   接下来解码
	var BigMing  Person
	err = decoder.Decode(&BigMing)
	if err != nil{
		log.Panic("deconder failed")
	}
	fmt.Printf("解码后的小明：%v \n",&BigMing)



}
