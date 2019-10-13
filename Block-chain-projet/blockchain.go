package main

import (
	"github.com/boltdb/bolt"
	"log"
)
//4.引入区块链   blockchain 结构重写
type BlockChain struct {
	//定义一个区块链数组
	//blocks []*Block
	db *bolt.DB
	tail []byte
}
const  blockChainDb ="blockChain.db"
const  blockBucket ="blockBucket"
//5.定义一个区块链
func NewBlockChain() *BlockChain  {
	//return &BlockChain{
	//	blocks:[]*Block{genesisBlock},
	//}
	//1.open bolt database
	var lastHash   []byte //读取区块最后一个hash
	db , err := bolt.Open(blockChainDb,0600,nil)
	//defer db.Close()
	if err !=nil{
		log.Panic("open database failed")
	}
	//2.fing bucket if it does not exist then create a database
	_ = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//bucket does not exist create
			bucket, err = tx.CreateBucket([]byte(blockBucket)) //变量b1 通常用一个文件管理
			if err != nil {
				log.Panic("创建bucket失败了")
			}
			//创建一个传世块，并作为第一个区块添加进区块链中
			genesisBlock := GenesisBlock()
			//hash 作为key  block的字节流作为value  这一步是写数据
			_ = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			_ = bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
			//测试一下解码的功能   先。。。。创建一个
			//blockBytes := bucket.Get(genesisBlock.Hash)  //
			//block := Deserialize(blockBytes)             //
			//fmt.Printf("block info : %v\n",block)//
		} else {
			//如果有bucket就返回  db 和 tail
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain {db,lastHash}
}
//往区块链塞 传世块
func GenesisBlock() *Block  {
	return  NewBlock("传世块",[]byte{})
}
//6.添加区块
func (bc *BlockChain)AddBlock(data string)  {
	//如何获取前区块哈希hash?
	db := bc.db
	lastHash := bc.tail
	_ = db.Update(func(tx *bolt.Tx) error {
		//完成数据的添加
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket 不应该为空")
		}
		//这里是创建区块
		block := NewBlock(data, lastHash) //区块创建完成  然后是gob里的写入
		//hash 作为key  block的字节流作为value  这一步是写数据
		_ = bucket.Put(block.Hash, block.Serialize())
		//上下这两部是添加到区块链db数据库中
		_ = bucket.Put([]byte("LastHashKey"), block.Hash)

		//数据库更新完成   然后是更新内存中的区块链 tail
		bc.tail = block.Hash
		return nil
	})
}
