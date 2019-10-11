package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//1.定义一个工作量证明结构

type ProofOFWork struct {
	//block
	block *Block
	//目标值
	//可以存储一个很大的数 可以有很多方法 赋值方法 比较方法
	target *big.Int
	//
}
//创建POW的函数
func NewProofOFWork(block *Block)*ProofOFWork  {
	pow :=ProofOFWork{
		block:  block,
		//target: nil,
	}
	//需要一个指定的难度值，现在是一个string类型需要转换
	targetStr :="0000f00000000000000000000000000000000000000000000000000000000000"
	//引入辅助变量，目的是把上面的难度转换成一个big。int
	tmpInt := big.Int{}
	//将难度赋值给big。int 指定16进制格式
	tmpInt.SetString(targetStr,16)
	pow.target =&tmpInt
	return &pow
}
// 3.提供不断计算hash的函数
// -run（）
func (pow *ProofOFWork) Run() ([]byte,uint64)  {
	//1.拼接数据 join 区块的数据和不断变化的随机数Nonce
	//2.进行hash运算
	//3.在target big int进行比较
		//a。找到了 退出返回
		//b 没找到，继续找，改变nonce的值
		//c。找到了就退出
	var nonce uint64
	var hash [32]byte
	block := pow.block
	for {
		//1.拼接数据
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
			block.Data,

		}
		blockInfo :=bytes.Join(tmp,[]byte{})
		//2.hash运算
		hash = sha256.Sum256(blockInfo)
		//3.找到与POW中的target进行比较
		tmpInt := big.Int{}
		//将得到的hash数组转换成一个big。int
		tmpInt.SetBytes(hash[:])
		//比较当前hash与目标hash 如果小于阈值则说明找到，否则继续
		if tmpInt.Cmp(pow.target) == -1{
			fmt.Printf("find it hash : %x ,nonce :%d \n",hash,nonce)

			//return  hash[:],nonce
			break

		} else {
			 nonce++
		}


	}
	return  hash[:],nonce

}
