package ecdsa_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"testing"
)

func TestEcdsa(t *testing.T)  {
	curve := elliptic.P256()
	privateKey , err :=ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil{
		log.Panic("1111")
	}
	pubKey := privateKey.PublicKey
	//公钥
	data := "hello world"
	hash := sha256.Sum256([]byte(data))
	r,s,err:=ecdsa.Sign(rand.Reader,privateKey,hash[:])
	if err != nil{
		log.Panic("22222")
	}
	//把 r，s 进行序列化传输

	//定义两个辅助的big。int
	r1 := big.Int{}
	s1 := big.Int{}
	signature := append(r.Bytes(),s.Bytes()...)
	//拆分signature 平均分给r和s
	r1.SetBytes(signature[0:len(signature)/2])
	s1.SetBytes(signature[len(signature)/2:])

	//signature := append(r.Bytes(),s.Bytes()...)
	fmt.Printf("pubKey :%v\n ",pubKey)
	fmt.Printf("r :%v,len :%d\n ",r.Bytes(),len(r.Bytes()))
	fmt.Printf("s :%v,len :%d\n ",s.Bytes(),len(s.Bytes()))
	//校验需要三个东西  签名 公钥 数据
	res := ecdsa.Verify(&pubKey,hash[:],r,s)
	fmt.Printf("校验结果 : %v\n",res)



}
