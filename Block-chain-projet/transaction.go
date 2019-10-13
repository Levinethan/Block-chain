package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const reward  =12.5
//1定义交易结构
type Transaction struct {
	TXID []byte
	TXInputs []TXInput  //交易输入数组
	TXOutputs []TXOutput  //交易输出数组
}
type TXInput struct {
	//引用的交易ID
	TXid []byte
	//output 的索引值
	Index int64
	//解锁脚本 用地址模拟
	Sig string
}
type TXOutput struct {
	//转账金额
	value float64
	//锁定脚本  用地址模拟
	PukKeyHash string
}
//设置交易ID
func (tx *Transaction) SetHash()  {
	var buffer  bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err :=encoder.Encode(tx)
	if err!=nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}
//2提供创建交易方法 (挖矿交易)(转账交易)
func NewCoinbaseTX(address string,data string) *Transaction  {
	//挖矿交易特点  1只有一个input 2 无需引用id 3 无需引用index 由于挖矿无需指定签名  所以sig字段可以是任何东西
	input := TXInput{[]byte{},-1,data}
	output := TXOutput{reward,address}
	tx :=Transaction{[]byte{},[]TXInput{input},[]TXOutput{output}}
	tx.SetHash()
	return &tx
}
//3创建挖矿交易
//4根据交易调整程序
