package main

import (
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const Usage  =`addBlock  --data DATA "add data to blockchain"
printChain  "print all blockchain data"`


func (cli *CLI)Run()  {
	args := os.Args
	//fmt.Printf("args : ",args)
	if len(args) <2{
		fmt.Printf(Usage)
		return
	}
	cmd := args[1]
	switch cmd {
	case "addBlock":
		fmt.Printf("添加区块")
	case "printChain":
		fmt.Printf("输出区块")
	default:
		fmt.Printf("无效的命令")
		fmt.Printf(Usage)
	}
}
