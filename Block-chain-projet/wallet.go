package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
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