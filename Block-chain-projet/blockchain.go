package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
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
	for _,tx:= range txs{
		if !bc.VerifyTransaction(tx){
			//矿工发现无效交易
			fmt.Printf("交易无效")
			return
		}
	}
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
//找到指定地址的所有的utxo
func (bc *BlockChain) FindUTXOs(pubKeyHash []byte) []TXOutput {
	var UTXO []TXOutput

	txs := bc.FindUTXOTransactions(pubKeyHash)

	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if bytes.Equal(pubKeyHash, output.PukKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}

	return UTXO
}
//找到合理的UTXO
func (bc *BlockChain) FindNeedUTXOs(senderPubKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	//找到的合理的utxos集合
	utxos := make(map[string][]uint64)
	var calc float64

	txs := bc.FindUTXOTransactions(senderPubKeyHash)

	for _, tx := range txs {
		for i, output := range tx.TXOutputs {
			//if from == output.PubKeyHash {
			//两个[]byte的比较
			//直接比较是否相同，返回true或false
			if bytes.Equal(senderPubKeyHash, output.PukKeyHash) {
				//fmt.Printf("222222")
				//UTXO = append(UTXO, output)
				//fmt.Printf("333333 : %f\n", UTXO[0].Value)
				//我们要实现的逻辑就在这里，找到自己需要的最少的utxo
				//3. 比较一下是否满足转账需求
				//   a. 满足的话，直接返回 utxos, calc
				//   b. 不满足继续统计

				if calc < amount {
					//1. 把utxo加进来，
					//utxos := make(map[string][]uint64)
					//array := utxos[string(tx.TXID)] //确认一下是否可行！！
					//array = append(array, uint64(i))
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
					//2. 统计一下当前utxo的总额
					//第一次进来: calc =3,  map[3333] = []uint64{0}
					//第二次进来: calc =3 + 2,  map[3333] = []uint64{0, 1}
					//第三次进来：calc = 3 + 2 + 10， map[222] = []uint64{0}
					calc += output.Value

					//加完之后满足条件了，
					if calc >= amount {
						//break
						fmt.Printf("找到了满足的金额：%f\n", calc)
						return utxos, calc
					}
				} else {
					fmt.Printf("不满足转账金额,当前总额：%f， 目标金额: %f\n", calc, amount)
				}
			}
		}
	}

	return utxos, calc
}

func (bc *BlockChain) FindUTXOTransactions(senderPubKeyHash []byte) []*Transaction {
	var txs []*Transaction //存储所有包含utxo交易集合
	//我们定义一个map来保存消费过的output，key是这个output的交易id，value是这个交易中索引的数组
	//map[交易id][]int64
	spentOutputs := make(map[string][]int64)

	//创建迭代器
	it := bc.NewIterator()

	for {
		//1.遍历区块
		block := it.Next()

		//2. 遍历交易
		for _, tx := range block.Transactions {
			//fmt.Printf("current txid : %x\n", tx.TXID)

		OUTPUT:
			//3. 遍历output，找到和自己相关的utxo(在添加output之前检查一下是否已经消耗过)
			//	i : 0, 1, 2, 3
			for i, output := range tx.TXOutputs {
				//fmt.Printf("current index : %d\n", i)
				//在这里做一个过滤，将所有消耗过的outputs和当前的所即将添加output对比一下
				//如果相同，则跳过，否则添加
				//如果当前的交易id存在于我们已经表示的map，那么说明这个交易里面有消耗过的output

				//map[2222] = []int64{0}
				//map[3333] = []int64{0, 1}
				//这个交易里面有我们消耗过得output，我们要定位它，然后过滤掉
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						//[]int64{0, 1} , j : 0, 1
						if int64(i) == j {
							//fmt.Printf("111111")
							//当前准备添加output已经消耗过了，不要再加了
							continue OUTPUT
						}
					}
				}

				//这个output和我们目标的地址相同，满足条件，加到返回UTXO数组中
				//if output.PubKeyHash == address {
				if bytes.Equal(output.PukKeyHash, senderPubKeyHash) {
					//fmt.Printf("222222")
					//UTXO = append(UTXO, output)

					//!!!!!重点
					//返回所有包含我的outx的交易的集合
					txs = append(txs, tx)

					//fmt.Printf("333333 : %f\n", UTXO[0].Value)
				} else {
					//fmt.Printf("333333")
				}
			}

			//如果当前交易是挖矿交易的话，那么不做遍历，直接跳过

			if !tx.IsCoinbase() {
				//4. 遍历input，找到自己花费过的utxo的集合(把自己消耗过的标示出来)
				for _, input := range tx.TXInputs {
					//判断一下当前这个input和目标（李四）是否一致，如果相同，说明这个是李四消耗过的output,就加进来
					//if input.Sig == address {
					//if input.PubKey == senderPubKeyHash  //这是肯定不对的，要做哈希处理
					pubKeyHash := HashPubKey(input.PubKey)
					if bytes.Equal(pubKeyHash, senderPubKeyHash) {
						//spentOutputs := make(map[string][]int64)
						//indexArray := spentOutputs[string(input.TXid)]
						//indexArray = append(indexArray, input.Index)
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
						//map[2222] = []int64{0}
						//map[3333] = []int64{0, 1}
					}
				}
			} else {
				//fmt.Printf("这是coinbase，不做input遍历！")
			}
		}

		if len(block.PrevHash) == 0 {
			break
			fmt.Printf("区块遍历完成退出!")
		}
	}

	return txs
}

func (bc *BlockChain)FindTransactionnByTXid(id []byte)(Transaction,error)  {
	//给 id   返回 Transaction
	//1遍历区块链
	//2遍历交易
	//3比较交易 找到了退出
	//4如果没找到 返回空Transaction  同时返回错误状态
	it := bc.NewIterator()
	for{
		block :=it.Next()
		for _,tx:= range block.Transactions{
			if bytes.Equal( tx.TXID,id){
				return *tx,nil
			}
		}
		if len(block.PrevHash)==0{
			fmt.Printf("区块遍历结束！\n")
			break

		}

	}
	return Transaction{},errors.New("无法找到交易")
}
func (bc *BlockChain)SignTransaction(tx *Transaction,privateKey *ecdsa.PrivateKey)  {

	prevTXs := make(map[string]Transaction)
	//找到所有引用的交易 UTXO
	for _,input := range tx.TXInputs{
		tx ,err:= bc.FindTransactionnByTXid(input.TXid)
		if err!=nil{
			log.Panic(err)
		}
		prevTXs[string(input.TXid)] = tx

	}
	tx.Sign(privateKey,prevTXs)
}

func (bc *BlockChain)VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase(){
		return true
	}
	prevTXs := make(map[string]Transaction)
	//找到所有引用的交易 UTXO
	for _, input := range tx.TXInputs {
		tx, err := bc.FindTransactionnByTXid(input.TXid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[string(input.TXid)] = tx

	}
	return tx.Verify(prevTXs)
}

