package main

import (
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const Usage  = "addBlock --data"
const PChain  = "printChain"



func (cli *CLI)Run()  {
	args := os.Args
	if len(args) <2{
		fmt.Printf(Usage)
		fmt.Printf(PChain)
		//cli.PrintBlockChain()
		return
	}
	cmd := args[1]
	switch cmd {
	case "addBlock":
		//获取数据  使用bc添加addBlock
		if len(args) ==4 &&args[2]== "--data" {
			data := args[3]
			cli.AddBlock(data)
			cli.PrintBlockChain()
		}else {
			fmt.Printf("添加数据使用不当")
		}
		fmt.Printf("添加区块\n")
	case "printChain":
		cli.PrintBlockChain()
		fmt.Printf("输出区块\n")
	default:
		fmt.Printf("无效的命令\n")

	}
}
