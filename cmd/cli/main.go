package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "strconv"

  "github.com/phuongdinh1411/blockchain/core"
)

type CLI struct {
  bc *core.BlockChain
}

func (cli *CLI) Run() {
  addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
  printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

  addBlockData := addBlockCmd.String("data", "", "Block data")

  switch os.Args[1] {
  case "addblock":
    err := addBlockCmd.Parse(os.Args[2:])
    if err != nil {
      log.Fatal(err)
    }
  case "printchain":
    err := printChainCmd.Parse(os.Args[2:])
    if err != nil {
      log.Fatal(err)
    }
  default:
    cli.printUsage()
    os.Exit(1)
  }
  if addBlockCmd.Parsed() {
    if *addBlockData == "" {
      addBlockCmd.Usage()
      os.Exit(1)
    }
    cli.addBlock(*addBlockData)
  }

  if printChainCmd.Parsed() {
    cli.printChain()
  }
}

func (cli *CLI) addBlock(data string) {
  cli.bc.AddBlock(data)
}
func (cli *CLI) printChain() {
  bci := cli.bc.Iterator()
  for {
    block := bci.Next()
    fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
    fmt.Printf("Data: %s\n", block.Data)
    fmt.Printf("Hash: %x\n", block.Hash)
    pow := core.NewProofOfWork(block)
    fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
    fmt.Println()

    if len(block.PrevBlockHash) == 0 { //hit the genesis block
      break
    }
  }
}
func (cli *CLI) printUsage() {}
func main() {
  bc := core.NewBlockChain()
  defer bc.DB.Close()
  cli := CLI{bc}
  cli.Run()
}
