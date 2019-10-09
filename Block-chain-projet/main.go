package main

import (
	"crypto/sha256";
	"fmt"
       )


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
		Hash:     []byte{},     //TODO
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
//4.引入区块链
type BlockChain struct {
	//定义一个区块链数组
	blocks []*Block
}
//5.定义一个区块链
func NewBlockChain() *BlockChain  {
	//创建一个传世块，并作为第一个区块添加进区块链中
	genesisblock :=GenesisBlock()
	return &BlockChain{
		blocks:[]*Block{genesisblock},
	}
}

//往区块链塞 传世块
func GenesisBlock() *Block  {
	return  NewBlock("传世块",[]byte{})
}
//6.添加区块
func (bc *BlockChain)AddBlock(data string)  {
	//如何获取前区块哈希hash?
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevHash := lastBlock.Hash
	//创建新的区块
	block := NewBlock(data,prevHash)
	//添加到区块链中
	bc.blocks = append(bc.blocks,block)
}
func main()  {
	bc := NewBlockChain()
	bc.AddBlock("向小U转账了100万")
	bc.AddBlock("我又向小U转账了100万")
	for i, block := range bc.blocks{
		fmt.Printf("当前区块高度：===== %d==== \n",i)
		fmt.Printf("当前区块HASH： %x\n",block.Hash)
		fmt.Printf("前区块HASH： %x\n",block.PrevHash)
		fmt.Printf("区块HASH数据： %s\n",block.Data)

	}


}

