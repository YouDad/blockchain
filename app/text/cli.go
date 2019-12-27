package text

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/YouDad/blockchain/core"
)

var (
	blockchain *core.Blockchain
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain - create a new blockchain")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func printChain() {
	iter := blockchain.Begin()

	for {
		block := iter.Next()
		if block == nil {
			break
		}
		pow := core.NewProofOfWork(block)
		fmt.Printf("Prev:   %x\n", block.PrevBlockHash)
		fmt.Printf("Hash:   %x\n", block.Hash)
		fmt.Printf("String: %s\n", block.App.ToString())
		fmt.Printf("PoW:    %s\n\n", strconv.FormatBool(pow.Validate()))
	}
}

func validateArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func Main() {
	validateArgs()

	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")
	switch os.Args[1] {
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if createBlockChainCmd.Parsed() {
		blockchain = core.CreateBlockchain()
	} else {
		blockchain = core.NewBlockchain()
	}
	defer blockchain.Close()

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		blockchain.MineBlock(GetAppString(*addBlockData))
	}

	if printChainCmd.Parsed() {
		printChain()
	}
}
