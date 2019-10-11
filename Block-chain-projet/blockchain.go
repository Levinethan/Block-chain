package main

//4.引入区块链   blockchain 结构重写
type BlockChain struct {
	//定义一个区块链数组
	blocks []*Block
	//db *bolt.DB


	tail []byte
}
//5.定义一个区块链
func NewBlockChain() *BlockChain  {
	//创建一个传世块，并作为第一个区块添加进区块链中
	genesisBlock :=GenesisBlock()
	return &BlockChain{
		blocks:[]*Block{genesisBlock},
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
