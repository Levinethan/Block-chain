package main

import (
	"bytes"
	"fmt"
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
func NewBlockChain(address string) *BlockChain  {
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
			genesisBlock := GenesisBlock(address)
			fmt.Printf("gensisBlock :%s\n",genesisBlock)
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
func GenesisBlock(address string) *Block  {
	coinbase := NewCoinbaseTX(address,"传世块")
	return  NewBlock([]*Transaction{coinbase},[]byte{})
}
//6.添加区块
func (bc *BlockChain)AddBlock(txs []*Transaction)  {
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
		block := NewBlock(txs, lastHash) //区块创建完成  然后是gob里的写入
		//hash 作为key  block的字节流作为value  这一步是写数据
		_ = bucket.Put(block.Hash, block.Serialize())
		//上下这两部是添加到区块链db数据库中
		_ = bucket.Put([]byte("LastHashKey"), block.Hash)
		//数据库更新完成   然后是更新内存中的区块链 tail
		bc.tail = block.Hash
		return nil
	})
}
func (bc *BlockChain)Printchain() {
	blockHeight := 0
	_ = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockBucket"))
		_ = b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastHashKey")) {
				return nil
			}
			block := Deserialize(v)
			fmt.Printf("=============区块高度 :%d =============\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号 %d\n", block.Version)
			fmt.Printf("前区块HASH： %x\n", block.PrevHash)
			fmt.Printf("merkel根 %x\n", block.MerkelRoot)
			fmt.Printf("时间戳 : %d\n", block.TimeStamp)
			fmt.Printf("随机数 :%d\n", block.Nonce)
			fmt.Printf("难度值 :%d\n", block.Difficulty)
			fmt.Printf("当前区块HASH： %x\n", block.Hash)
			//fmt.Printf("区块HASH数据： %s\n", block.Transactions[0].TXInputs[0].Sig)
			return nil
		})
		return nil
	})
}
//找到指定地址所有的UTXO 遍历数组【】TXOutput
func (bc *BlockChain) FindUTXOs (address string)[]TXOutput{
	var UTXO []TXOutput
	spentOutputs := make(map[string][]int64)
	//map[交易id][]int64
	//遍历区块 交易 output input 剩余的utxo 和花了的utxo
	//创建迭代器
	it := bc.NewInterator()
	for   {
		block := it.Next()
		for _, tx := range block.Transactions{
			fmt.Printf("current txid :%x\n",tx.TXID)
	OUTPUT:
			for i,output:= range tx.TXOutputs{    //遍历output 找到和自己相关的utxo 检查
				fmt.Printf("current index :%d\n",i)
				//这里进行一个过滤，过滤消耗过的outputs和当前的添加的output对比一下
				//如果相同 跳过 否则跳过
				//如果当前交易id存在与已经标识的map 说明交易有消耗过的output
				if spentOutputs[string(tx.TXID)]!=nil{
					for _,j :=range spentOutputs[string(tx.TXID)]{
						if int64(i) == j {
							//当前准备添加的output已经消耗过 不用添加
							fmt.Printf("11111")
							continue OUTPUT
						}
					}
				}
				if output.PukKeyHash == address {
					//fmt.Printf("22222")
					UTXO = append(UTXO, output)
					//fmt.Printf("33333 :%f \n",UTXO[0].Value)
				}else {
					//fmt.Printf("33333")
				}
			}
			//如果当前交易是挖矿交易，那么
			if !tx.IsCoinbase() {
			//遍历input 找到自己花费过的utxo集合
				for _,input := range tx.TXInputs{
				//判断一下这个input和目标 是否一致，如果相同说明是目标消耗过的output
					if input.Sig == address {
					//indexArray := spentOutputs[string(input.TXid)]
					//indexArray = append(indexArray,input.Index)
						spentOutputs[string(input.TXid)]= append(spentOutputs[string(input.TXid)],input.Index)
					}
				}
			}else {
				fmt.Printf("这是coinbase，不做遍历！")
			}
			//定义一个map来保存消费过的output key是这个output的交易id value是这个交易的索引数组
			//map[交易id][]int64
		}
		if len(block.PrevHash)==0{
			break
			fmt.Printf("区块遍历完成 ，并退出")
		}
	}
	return UTXO
}
//找到合理的UTXO
func (bc *BlockChain)FindNeedUTXOs(from string,amount float64)( map[string][]uint64,float64)  {
	utxos  := make(map[string][]uint64) //找到合理的utxo合集
	var calc  float64  //找到utxo里的总数
	spentOutputs := make(map[string][]int64)
	//111111111111111111111111111111111111
	it := bc.NewInterator()
	for   {
		block := it.Next()
		for _, tx := range block.Transactions{
			fmt.Printf("current txid :%x\n",tx.TXID)
		OUTPUT:
			for i,output:= range tx.TXOutputs{    //遍历output 找到和自己相关的utxo 检查
				//fmt.Printf("current index :%d\n",i)
				//这里进行一个过滤，过滤消耗过的outputs和当前的添加的output对比一下
				//如果相同 跳过 否则跳过
				//如果当前交易id存在与已经标识的map 说明交易有消耗过的output
				if spentOutputs[string(tx.TXID)]!=nil{
					for _,j :=range spentOutputs[string(tx.TXID)]{
						if int64(i) == j {
							//当前准备添加的output已经消耗过 不用添加
							///fmt.Printf("11111")
							continue OUTPUT
						}
					}
				}
				if output.PukKeyHash == from {
					//fmt.Printf("22222")
					//UTXO = append(UTXO, output)
					//fmt.Printf("33333 :%f \n",UTXO[0].Value)
					//逻辑实现处 找到自己需要的最少UTXO
					//TODO
					//把UTXO加进来
					//统计一下utxo
					//第一次进来.calc=3 map[33333]=[]uint64{0}
					//2次 calc=3+2 map[33333] = []uint64{0,1}
					//3次 calc = 3+2+10  map[22222]=[]uint64{0}
					//比较一下是否满足转账需求
					if calc < amount{
						utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)],uint64(i))
						calc += output.Value
						if calc >= amount{
							fmt.Printf("找到满足的余额 :%f\n",calc)
							return  utxos,calc
						}
						//array :=utxos[string(tx.TXID)]
						//array = append(array, uint64(i))
					}else {
						fmt.Printf("不满足转账金额，当前总额:%f ,目标金额:%f \n",calc,amount)
					}
					//满足的话返回utxos calc
					//否则继续统计
				}else {
					//fmt.Printf("3333333")
				}
			}
			//如果当前交易是挖矿交易，那么
			if !tx.IsCoinbase() {
				//遍历input 找到自己花费过的utxo集合
				for _,input := range tx.TXInputs{
					//判断一下这个input和目标 是否一致，如果相同说明是目标消耗过的output
					if input.Sig == from {
						//indexArray := spentOutputs[string(input.TXid)]
						//indexArray = append(indexArray,input.Index)
						spentOutputs[string(input.TXid)]= append(spentOutputs[string(input.TXid)],input.Index)
					}
				}
			}else {
				fmt.Printf("这是coinbase，不做遍历！")
			}
			//定义一个map来保存消费过的output key是这个output的交易id value是这个交易的索引数组
			//map[交易id][]int64
		}
		if len(block.PrevHash)==0{
			break
			fmt.Printf("区块遍历完成 ，并退出")
		}
	}
	//TODO
	return utxos,calc
}
	//22222222222222222222222222222222222

