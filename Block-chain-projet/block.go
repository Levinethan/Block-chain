package main

import (
	"bytes"

	"encoding/binary"
	"log"
	"time"
)

//1.定义区块结构
type Block struct {
	//1.版本号
	Version uint64
	//2.Merkel根
	MerkelRoot []byte
	//3.时间戳
	TimeStamp uint64
	//4.难度梯度
	Difficulty uint64
	//5.随机数 噪音
	Nonce uint64
	//6.前区块hash
	PrevHash []byte
	//当前区块HASH
	Hash []byte
	//区块数据
	Data []byte
}

//2.创建区块
func NewBlock(data string,prevBlockHash []byte) *Block  {
	block :=Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{},
		Data:     []byte(data),
		Version:00,
		MerkelRoot:[]byte{},
		TimeStamp:uint64(time.Now().Unix()),
		Nonce:0,               //随机值  无法知道区块阈值
		Difficulty:0,        //无效值   无法知道困难梯度
	}
	//block.SetHash()
	pow := NewProofOFWork(&block)

	//创建一个pow对象
	hash,nonce := pow.Run()
	//查找随机数 不停的进行hash运算
	block.Hash = hash
	//根据挖矿结果对区块进行更行或者补充。
	block.Nonce = nonce
	return &block
	//最后返回给block


}

//3.生成HASH
//func (block *Block) SetHash()  {
	//var blockInfo  []byte
	//拼装数据
	/*
	blockInfo =append(blockInfo,Uint64ToByte(block.Version)...)
	blockInfo =append(blockInfo,block.PrevHash...)
	blockInfo =append(blockInfo,block.MerkelRoot...)
	blockInfo =append(blockInfo,Uint64ToByte(block.TimeStamp)...)
	blockInfo =append(blockInfo,Uint64ToByte(block.Difficulty)...)
	blockInfo =append(blockInfo,Uint64ToByte(block.Nonce)...)
	blockInfo =append(blockInfo,block.Data...)

	 */ //优化 用join函数
	/* tmp := [][]byte{
		 Uint64ToByte(block.Version),
		 block.PrevHash,
		 block.MerkelRoot,
		 Uint64ToByte(block.TimeStamp),
		 Uint64ToByte(block.Difficulty),
		 Uint64ToByte(block.Nonce),
		 block.Data,

	 }
	 blockInfo :=bytes.Join(tmp,[]byte{})
	 //将二维的数组链接起来，返回一个一维的切片


	//sha256加密
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

	 */
//4.实现一个辅助函数 用来转换数据类型
func Uint64ToByte(num uint64 ) []byte {      //func XXX()括号里面 传入参数 []byte 返回一个byte类型数据
	var buffer  bytes.Buffer
	err := binary.Write(&buffer,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}
	return buffer.Bytes()

}
