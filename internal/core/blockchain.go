package core

import (
  "log"

  "github.com/boltdb/bolt"
)

type BlockChain struct {
  tip []byte
  DB  *bolt.DB
}

type BlockChainIterartor struct {
  currentHash []byte
  db          *bolt.DB
}

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

func (bc *BlockChain) AddBlock(data string) {
  var lastHash []byte

  err := bc.DB.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    lastHash = b.Get([]byte("l"))

    return nil
  })
  if err != nil {
    log.Fatal(err)
  }

  newBlock := NewBlock(data, lastHash)

  err = bc.DB.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    err = b.Put(newBlock.Hash, newBlock.Serialize())
    err = b.Put([]byte("l"), newBlock.Hash)
    bc.tip = newBlock.Hash

    return nil
  })
}

func NewBlockChain() *BlockChain {
  db, err := bolt.Open(dbFile, 0600, nil)
  if err != nil {
    log.Fatal(err)
  }
  var tip []byte
  err = db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))

    if b == nil {
      genesis := NewGenesisBlock()
      b, err = tx.CreateBucket([]byte(blocksBucket))
      err = b.Put([]byte("l"), genesis.Hash)
      err = b.Put(genesis.Hash, genesis.Serialize())
      tip = genesis.Hash
    } else {
      tip = b.Get([]byte("l"))
    }
    return nil // successful transaction
  })
  bc := BlockChain{tip, db}
  return &bc
}

func (bc *BlockChain) Iterator() *BlockChainIterartor {
  bci := &BlockChainIterartor{bc.tip, bc.DB}

  return bci
}

//from tail to head in blockchain
func (i *BlockChainIterartor) Next() *Block {
  var block *Block
  i.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    encodedBlock := b.Get(i.currentHash)
    block = Deserialize(encodedBlock)
    return nil
  })
  i.currentHash = block.PrevBlockHash
  return block
}
