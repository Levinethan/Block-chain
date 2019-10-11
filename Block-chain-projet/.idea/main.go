package main




func main()  {
	bc := NewBlockChain()
	bc.AddBlock("向小U转账了100万")
	bc.AddBlock("我又向小U转账了100万")
	/*for i, block := range bc.blocks{
		fmt.Printf("当前区块高度：===== %d==== \n",i)
		fmt.Printf("当前区块HASH： %x\n",block.Hash)
		fmt.Printf("前区块HASH： %x\n",block.PrevHash)
		fmt.Printf("区块HASH数据： %s\n",block.Data)

	}
*/

}

