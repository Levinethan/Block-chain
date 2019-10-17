package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strings"
)
const reward  =  50
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
	//Sig string
	Signature []byte
	PubKey []byte
	//注意这里是公钥 不是hash 也不是地址
}
type TXOutput struct {
	//转账金额
	Value float64
	//锁定脚本  用地址模拟
	PukKeyHash []byte
	//这个是收款方公钥的hash  注意是公钥的hasht 不是地址

}
//由于现在存储字段是地址的公钥hash 无法直接创建Txoutput
//为了能得到公钥hash 根据思维导图 逆向  写一个Lock函数
func (output *TXOutput)Lock(address string)  {
	//1解码
	//2截取公钥hash 去除version 1字节 去除校验码  4字节
	output.PukKeyHash = GetPubKeyFromAddress(address)
	//真正的锁定动作
}
//给TXOutput 提供一个创建方法 否则无法调用Lock
func NewTXOutput(value float64,address string) *TXOutput  {
	output := TXOutput{
		Value:      value,

	}
	output.Lock(address)
	return &output
}
//完成input output 的改写
//我们重新定义了input和output 的结构struck  把原先用来模拟的 sig签名string  把version 和 检验 5byte去掉 剩下的就是公钥hash
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
	//签名先填的空 后面创建完整交易后 最后做一次签名就可以

	input := TXInput{[]byte{},-1,nil,[]byte(data)}
	output := NewTXOutput(reward,address)
	tx :=Transaction{[]byte{},[]TXInput{input},[]TXOutput{*output}}
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
	//创建交易之后  要进行数字签名  需要私钥 - 打开wallet 从内存中加载NewWallets()
	//然后是找到自己的钱包  根据自己的地址 找到自己的钱包
	//返回钱包  公钥私钥都有了
	ws := NewWallets()
	wallet := ws.WalletsMap[from]
	if wallet ==nil{
		fmt.Printf("无法找到该钱包的私钥，交易失败!\n")
		return nil
	}
	pubKey := wallet.PubKey

	privateKey := wallet.Private
	pubKeyHash := HashPubKey(pubKey)
	//===================================================================
	//1. 找到最合理UTXO集合 map[string][]uint64
	//TODO
	utxos ,resValue:= bc.FindNeedUTXOs(pubKeyHash, amount)  //[]byte  字节流
	//utxos ,resValue:= make(map[string][]uint64),0.0   //临时给个0 编译测试一下
	if resValue < amount {
		fmt.Printf("余额不足，交易失败!")
		return nil
	}
	var inputs []TXInput
	var outputs []TXOutput
	//2. 创建交易输入, 将这些UTXO逐一转成inputs
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{[]byte(id), int64(i), nil,pubKey}
			inputs = append(inputs, input)
		}
	}
	//创建交易输出
	//output := TXOutput{amount, to}
	output := NewTXOutput(amount,to)
	outputs = append(outputs, *output)
	//找零
	if resValue > amount {
		output = NewTXOutput(resValue-amount,from)
		outputs = append(outputs,*output)
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()
	//签名，交易创建后进行签名
	bc.SignTransaction(&tx,privateKey)
	return &tx
}
func (tx *Transaction)Sign(privateKey *ecdsa.PrivateKey,prevTXs map[string]Transaction)  {
	//具体签名
	if tx.IsCoinbase(){
		return
	}
	//1创建一个当前交易的copy  TrimmeCopy  TXCopy 把Signature 和 Pubkey
	txCopy:= tx.TrimmedCopy()
	//2循环遍历txCopy的inputs 得到 input 索引 output的公钥hash
	for i,input := range txCopy.TXInputs{
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) ==0{
			log.Panic("引用的交易无效")
		}
		//3生成要签名的数据。签一个hash值
		//对每个input都要签名一次 由当前input引用的output的哈希
		//不要对input进行赋值 这是一个副本 要对txCopy。TXInputs【】进行操作否则无法把pubkey传进来
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PukKeyHash
		//对拼接好的txCopy进行哈希处理，Sethash得到TXID ，这个TXID是我们要签名的id
		//3生成要签名的数据。签一个hash值
		txCopy.SetHash()
		//还原
		txCopy.TXInputs[i].PubKey=nil
		signDataHash := txCopy.TXID
		//4	签名动作得到r，s字节流 byte[]
		r,s,err :=ecdsa.Sign(rand.Reader,privateKey,signDataHash)
		if err != nil{
			log.Panic(err)
		}
		//5放到我们所签名的input的Signature中
		signature := append(r.Bytes(),s.Bytes()...)
		tx.TXInputs[i].Signature = signature


	}






}
func (tx *Transaction)TrimmedCopy()Transaction  {
	var inputs   []TXInput
	var outputs  []TXOutput
	for _,input := range tx.TXInputs{inputs = append(inputs,TXInput{input.TXid,input.Index,nil,nil})

	}
	for _,output :=range tx.TXOutputs{
		outputs = append(outputs,output)
	}
	return Transaction{tx.TXID,inputs,outputs}

}
//校验
//对每个签名过的input进行校验
//把coinbase除外
func (tx *Transaction)Verify(prevTXs map[string]Transaction) bool{
	if tx.IsCoinbase(){
		return true
	}
	//1得到签名数据
	txCopy := tx.TrimmedCopy()

	//2.得到signature 反推回r，s
	for i,input := range tx.TXInputs{
		prevTXs := prevTXs[string(input.TXid)]
		if len(prevTXs.TXID)==0 {
			log.Panic("引用无效")
		}
		txCopy.TXInputs[i].PubKey = prevTXs.TXOutputs[input.Index].PukKeyHash
		txCopy.SetHash()

		dataHash := txCopy.TXID
		//3拆解pubkey
		signature := input.Signature  //拆  r s
		pubKey := input.PubKey       //拆 X Y
		//定义两个 big int
		r := big.Int{}
		s := big.Int{}

		//拆分signature 平均分给r和s
		r.SetBytes(signature[0:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])
		X := big.Int{}
		Y := big.Int{}

		//拆分signature 平均分给r和s
		X.SetBytes(pubKey[0:len(pubKey)/2])
		Y.SetBytes(pubKey[len(pubKey)/2:])
		pubKeyOrign := ecdsa.PublicKey{elliptic.P256(),&X,&Y}
		//verify
		if !ecdsa.Verify(&pubKeyOrign,dataHash,&r,&s){
			return false
		}
	}

	return true

}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.TXID))

	for i, input := range tx.TXInputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TXid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.TXOutputs{
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x",output.PukKeyHash ))
	}

	return strings.Join(lines, "\n")
}
