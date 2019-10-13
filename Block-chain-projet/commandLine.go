package main

import "fmt"

func (cli *CLI)AddBlock (data string) {
	//cli.bc.AddBlock(data) TODO
	fmt.Printf("添加区块成功 \n")
}
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
		fmt.Printf("时间戳 : %d\n",block.TimeStamp)
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
		total += utxo.value
	}
	fmt.Printf("%s的余额为: %f\n ",address,total)


}
//逻辑 cli调用   实现commandline 实现