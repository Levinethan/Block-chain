package main

import "crypto/sha256"

//1.定义区块结构
type Block struct {
	//1.前区块hash
	PrevHash []byte
	//2.当前区块HASH
	Hash []byte
	//3.区块数据
	Data []byte
}

//2.创建区块
func NewBlock(data string,prevBlockHash []byte) *Block  {
	block :=Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	block.SetHash()
	return &block

}

//3.生成HASH
func (block *Block) SetHash()  {
	//拼装数据
	blockInfo := append(block.PrevHash,block.Data...)
	//sha256加密
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
