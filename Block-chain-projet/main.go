package main

//import "fmt"



func main()  {
	bc := NewBlockChain("张三")
	cli :=CLI{bc:bc}
	cli.Run()


	//bc.AddBlock("111111111111111111")
	//bc.AddBlock("222222222222222222")
	//it := bc.NewInterator()
/*
	for {
		block := it.Next() //返回区块 然后左移
		fmt.Printf("版本号 %d\n",block.Version)
		fmt.Printf("前区块HASH： %x\n",block.PrevHash)
		fmt.Printf("merkel根 %x\n",block.MerkelRoot)
		fmt.Printf("时间戳 : %d\n",block.TimeStamp)
		fmt.Printf("随机数 :%d\n",block.Nonce)
		fmt.Printf("难度值 :%d\n",block.Difficulty)
		fmt.Printf("当前区块HASH： %x\n",block.Hash)
		fmt.Printf("区块HASH数据： %s\n",block.Data)
		if len(block.PrevHash)==0 {
			fmt.Printf("遍历结束")
			break
		}
	}

*/
}

