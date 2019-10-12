package main

//import "fmt"

func main()  {
	bc := NewBlockChain()
	cli :=CLI{bc:bc}
	cli.Run()


	//bc.AddBlock("111111111111111111")
	/*bc.AddBlock("222222222222222222")
	it := bc.NewInterator()
	for {
		block := it.Next() //返回区块 然后左移
		fmt.Printf("当前区块HASH： %x\n",block.Hash)
		fmt.Printf("前区块HASH： %x\n",block.PrevHash)
		fmt.Printf("区块HASH数据： %s\n",block.Data)
		if len(block.PrevHash)==0 {
			fmt.Printf("遍历结束")
			break
		}
	}
	/*for i, block := range bc.blocks{
		fmt.Printf("当前区块高度：===== %d==== \n",i)
		fmt.Printf("当前区块HASH： %x\n",block.Hash)
		fmt.Printf("前区块HASH： %x\n",block.PrevHash)
		fmt.Printf("区块HASH数据： %s\n",block.Data)

	}
*/

}

