package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"

	"log"
)

type Wallet struct {

	//私钥
	Private *ecdsa.PrivateKey
	PubKey []byte
}

func NewWallet()*Wallet  {
	curve := elliptic.P256()
	privateKey , err :=ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil{
		log.Panic("1111")
	}
	pubKeyOrig := privateKey.PublicKey
	//公钥
	pubKey := append(pubKeyOrig.X.Bytes(),pubKeyOrig.Y.Bytes()...)
	return &Wallet{Private:privateKey,PubKey:pubKey}
}


//
func (w *Wallet)NewAddress()string  {
	pubKey := w.PubKey
	hash:=sha256.Sum256(pubKey)


	//编码
	rip160hasher := ripemd160.New()
	//rip160hasher := crypto.RIPEMD160.New()
	_, err :=rip160hasher.Write(hash[:])
	if err!=nil{
		log.Panic()
	}
	//返回rip160的哈希结果
	rip160HashValue :=rip160hasher.Sum(nil)
	version := byte(00)
	payload := append([]byte{version},rip160HashValue...)
	//checksum
	//两次sha
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	//前四字节byte
	checkCode := hash2[:4]
	//25字节数据
	payload = append(payload,checkCode...)
	//比特币全节点源码
	address := base58.Encode(payload)
	return address
}

func HashPubKey(data []byte) []byte {
	hash := sha256.Sum256(data)

	//理解为编码器
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])

	if err != nil {
		log.Panic(err)
	}

	//返回rip160的哈希结果
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
}
func CheckSum(data []byte) []byte {
	//两次sha256
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])

	//前4字节校验码
	checkCode := hash2[:4]
	return checkCode
}
func IsValidAddress(address string)bool  {
	addressByet := base58.Decode(address)
	if len(addressByet)<4{
		return false
	}
	//解码
	payload := addressByet[:len(addressByet)-4]
	checksum1 := addressByet[len(addressByet)-4:]
	checksum2 := CheckSum(payload)
	fmt.Printf("checksum1 :%x\n",checksum1)
	fmt.Printf("checksum2 :%x\n",checksum2)
	return bytes.Equal(checksum1,checksum2)
}