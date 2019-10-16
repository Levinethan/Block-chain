package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
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
	//Data []byte
	Transactions []*Transaction   //真实交易数组
}

//2.创建区块
func NewBlock(txs []*Transaction ,prevBlockHash []byte) *Block  {
	block :=Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{},
		Transactions :txs,
		Version:00,
		MerkelRoot:[]byte{},
		TimeStamp:uint64(time.Now().Unix()),
		Nonce:0,               //随机值  无法知道区块阈值
		Difficulty:0,        //无效值   无法知道困难梯度
	}
	//block.SetHash()
	block.MerkelRoot = block.MakeMerkelRoot()
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
func (block *Block)Serialize()[]byte  {  //编码工作 序列化  转成字节流
	var buffer  bytes.Buffer
	//使用一个辅助变量buffer  编码好的数据放进  buffer
	//分别定义一个  编码器  和解码器
	//1.编码器
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&block)

	if err != nil {
		log.Panic("编码出错，小明失踪")
	}

	return buffer.Bytes()
}

func Deserialize(data []byte) Block  {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	//那么数据在解码器中   接下来解码
	var block Block

	err := decoder.Decode(&block)
	if err != nil{
		log.Panic("deconder failed")
	}
	return block
}






//4.实现一个辅助函数 用来转换数据类型
func Uint64ToByte(num uint64 ) []byte {      //func XXX()括号里面 传入参数 []byte 返回一个byte类型数据
	var buffer  bytes.Buffer
	err := binary.Write(&buffer,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}
	return buffer.Bytes()

}
func (block *Block)MakeMerkelRoot() []byte  {
	//TODO
	var info  []byte

	for _,tx := range block.Transactions{
		info = append(info,tx.TXID...)

	}
	hash := sha256.Sum256(info)

	return hash[:]
} //对交易数据进行简单拼接 不做二叉树处理