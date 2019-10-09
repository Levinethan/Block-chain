package main

import (
	"crypto/sha256"
	"fmt"
)
//定义区块结构
type Block struct {
	//1.前区块hash
	PrevHash []byte
	//2.当前区块HASH
	Hash []byte
	//3.区块数据
	Data []byte
}

//创建区块
func NewBlock(data string,prevBlockHash []byte) *Block  {
	block :=Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	block.SetHash()
	return &block

}
//输出区块信息
func main()  {
	block :=NewBlock("转让了一枚比特币",[]byte{})
	fmt.Printf("当前区块HASH： %x\n",block.Hash)
	fmt.Printf("前区块HASH： %x\n",block.PrevHash)
	fmt.Printf("区块HASH数据： %s\n",block.Data)

}

//生成HASH
func (block *Block) SetHash()  {
	blockInfo := append(block.PrevHash,block.Data...)
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
