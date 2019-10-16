package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
	"log"
)

//定义一个wallat 钱包   用小写   保存所有的钱包地址
//创建方法   其实就是把  地址全部加载进来
//那么就是把新建的wallet 添加进来
//保存有了   那么还有读取
//

type Wallets struct {
	WalletsMap map[string] *Wallet
}

func NewWallets() *Wallets  {
	var ws Wallets
	//ws.WalletsMap = make(map[string]*Wallet)
	ws.loadFile()
	return &ws


}

func (ws *Wallets)CreateWallet()string  {
	wallet := NewWallet()    //创建wallet 的两个元素
	address := wallet.NewAddress()
	//var wallets  Wallets   //返回一个 *Wallets 指针 就要创建一个 对吧。。然后
	//wallets.WalletsMap = make(map[string]*Wallet)  //创建 对象不然的话是空  对吧  类型就是map string
	ws.WalletsMap[address]=wallet           //   然后往里面塞数据地址对吧  类型就是map string
	//创建完了 保存到本地

	ws.saveToFile()
	return address
}

func (ws *Wallets)saveToFile()  {
	var buffer bytes.Buffer
	//gob  如果Encode Decode 的类型是interface
	//那么要注册一下
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err :=encoder.Encode(ws)
	//一定一定要校验
	if err != nil{
		log.Panic(err)
	}

	ioutil.WriteFile("wallet.dat",buffer.Bytes(),0644)
}

func (ws *Wallets)loadFile()  {
	content,err := ioutil.ReadFile("wallet.dat")
	if err!= nil{
		log.Panic(err)
	}
	//解码
	decoder := gob.NewDecoder(bytes.NewReader(content))
	gob.Register(elliptic.P256())
	var  wsLocal Wallets
	err= decoder.Decode(&wsLocal)
	if err!= nil{
		log.Panic(err)
	}
	//ws = &wslocal  对于结构来说 里面有map的 要指定赋值 不要在最外层直接赋值
	ws.WalletsMap = wsLocal.WalletsMap
}
func (ws *Wallets)ListAllAddresses()[]string  {
	var addresses  []string
	for address := range ws.WalletsMap{
		addresses = append(addresses,address)
	}
	return  addresses
}
func GetPubKeyFromAddress(address string)[]byte  {
	addressByte := base58.Decode(address)  //25byte
	len1 := len(addressByte)
	pubKeyHash := addressByte[1:len1-4]
	 return pubKeyHash
}