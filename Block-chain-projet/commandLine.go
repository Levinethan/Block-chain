package main

import (
	"fmt"
	"time"
)


//正向打印
func (cli *CLI)PrintBlockChain () {
	cli.bc.Printchain()
	fmt.Printf("打印区块完成")
}
//反向打印
func (cli *CLI)PrintBlockChainReverse () {
	bc := cli.bc
	it := bc.NewInterator()
	for {
		block := it.Next() //返回区块 然后左移
		fmt.Printf("版本号 %d\n",block.Version)
		fmt.Printf("前区块HASH： %x\n",block.PrevHash)
		fmt.Printf("merkel根 %x\n",block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeStamp),0).Format("2006-01-02 15:04:05")
		fmt.Printf("时间戳 : %s\n",timeFormat)
		fmt.Printf("随机数 :%d\n",block.Nonce)
		fmt.Printf("难度值 :%d\n",block.Difficulty)
		fmt.Printf("当前区块HASH： %x\n",block.Hash)
		fmt.Printf("区块HASH数据： %s\n",block.Transactions[0].TXInputs[0].Sig)
		if len(block.PrevHash)==0 {
			fmt.Printf("遍历结束")
			break
		}
	}
}
func (cli *CLI)GetBalance(address string)  {
	utxos :=cli.bc.FindUTXOs(address)
	total := 0.0
	for _, utxo := range utxos{
		total += utxo.Value
	}
	fmt.Printf("%s的余额为: %f\n ",address,total)


}
func (cli *CLI)Send(from,to string,amount float64,miner ,data string)  {
	fmt.Printf("from :%s\n",from)
	fmt.Printf("to :%s\n",to)
	fmt.Printf("amount :%f\n",amount)
	fmt.Printf("miner :%s\n",miner)
	fmt.Printf("data :%s\n",data)
	coinbase := NewCoinbaseTX(miner,data)
	tx := NewTransaction(from,to,amount,cli.bc)
	if tx == nil{
		fmt.Printf("无效的交易")
		return
	}
	cli.bc.AddBlock([]*Transaction{coinbase,tx})
	fmt.Printf("转账成功")
}
//逻辑 cli调用   实现commandline 实现

func (cli *CLI)NewWallet()  {
	wallet := NewWallet()
	address := wallet.NewAddress()
	fmt.Printf("私钥  %v\n",wallet.Private)
	fmt.Printf("公钥 %v\n",wallet.PubKey)
	fmt.Printf("地址 %s\n",address)
}