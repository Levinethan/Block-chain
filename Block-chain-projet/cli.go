package main
import (
	"fmt"
	"os"
	"strconv"
)
type CLI struct {
	bc *BlockChain
}
const Usage = `
	
	printChain               "正向打印区块链"
	printChainR              "反向打印区块链"
	getBalance --address  ADDRESS   "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA  "由FROM转钱amount给TO，有miner挖矿，同时写入data"
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
	case "send":
		fmt.Printf("转账中...\n")
		if len(args)!=7{
			fmt.Printf("参数个数错误")
			fmt.Printf(Usage)
			return
		}

		from := args[2]
		to := args[3]
		amount ,_ := strconv.ParseFloat(args[4],64)
		miner := args[5]
		data := args[6]
		cli.Send(from,to,amount,miner,data)

	default:
		fmt.Printf("无效的命令\n")
		fmt.Printf(Usage)
	}
}
