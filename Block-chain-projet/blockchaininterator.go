package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainInterator struct {
	db *bolt.DB
	//数据库
	currentHashPointer []byte
	//当前指针 用来寻找区块 代替for循环
}

func (bc *BlockChain) NewIterator()*BlockChainInterator {
	return &BlockChainInterator{
		db:                 bc.db,
		currentHashPointer: bc.tail, // 最初指向区块链的最后一个区块，随着Next的调用 不断变化
	}
}
func (it *BlockChainInterator) Next()  *Block {
	var block  Block
	//返回当前区块，然后指针前移
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil{
			log.Panic("迭代器遍历时,bucket不应该为空,请检查")
		}
		blockTem := bucket.Get(it.currentHashPointer)
		block =Deserialize(blockTem)
		it.currentHashPointer = block.PrevHash

		return nil
	})
	return &block
}