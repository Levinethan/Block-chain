package main
import (
	"fmt"
	"os"
)
type CLI struct {
	bc *BlockChain
}
const Usage = `
	addBlock --data DATA     "添加区块"
	printChain               "正向打印区块链"
	printChainR              "反向打印区块链"
	getBalance --address  ADDRESS   "获取指定地址的余额"
`
func (cli *CLI)Run()  {
	args := os.Args
	if len(args) <2{
		fmt.Printf(Usage)
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
		}else {
			fmt.Printf("添加数据使用不当")
			fmt.Printf(Usage)
		}
		fmt.Printf("添加区块\n")
	case "printChain":
		fmt.Printf("正向打印区块\n")
		cli.PrintBlockChain()
	case "printChainR":
		fmt.Printf("反向打印区块\n")
		cli.PrintBlockChainReverse()
	case "getBalance":
		fmt.Printf("获取余额\n")
		if len(args) ==4 && args[2] =="--address"{
			address := args[3]
			cli.GetBalance(address)
		}
	default:
		fmt.Printf("无效的命令\n")
		fmt.Printf(Usage)
	}
}
