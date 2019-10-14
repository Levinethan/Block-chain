package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward  =500
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
	Value float64
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
func (tx *Transaction)IsCoinbase() bool  {
	//1.交易的input只有一个
	//2.交易id为空
	//3.交易的index 为-1
	//同时满足这三个条件
	//if len(tx.TXInputs)==1{
	//	input := tx.TXInputs[0]
	//	if !bytes.Equal(input.TXid , []byte{}) || input.Index!=-1{
	//		return false
	//	}
	//}
	//return true
	if len(tx.TXInputs)==1 &&len(tx.TXInputs[0].TXid)==0 &&tx.TXInputs[0].Index==-1{
		return true
	}
	return false
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

//创建普通的转账交易
//找到合理的UTXO集合 map[string][]uint64   这里是最难的
//找到的UTXO 逐一转成input
//input 创建完成  然后创建output
//创建完成普通转账  如果有找零则返回自己账号
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {

	//1. 找到最合理UTXO集合 map[string][]uint64
	utxos ,resValue:= bc.FindNeedUTXOs(from, amount)

	if resValue < amount {
		fmt.Printf("余额不足，交易失败!")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	//2. 创建交易输入, 将这些UTXO逐一转成inputs
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{[]byte(id), int64(i), from}
			inputs = append(inputs, input)
		}
	}

	//创建交易输出
	output := TXOutput{amount, to}
	outputs = append(outputs, output)

	//找零
	if resValue > amount {
		outputs = append(outputs, TXOutput{resValue - amount, from})
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()
	return &tx
}
