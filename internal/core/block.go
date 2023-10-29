package core

import (
  "bytes"
  "crypto/sha256"
  "encoding/gob"
  "log"
  "strconv"
  "time"
)

type Block struct {
  TimeStamp     int64 // when this block is created
  PrevBlockHash []byte
  Data          []byte
  Hash          []byte
  Nonce         int // use to validate a block
}

func (b *Block) SetHash() {
  timestamp := []byte(strconv.FormatInt(b.TimeStamp, 10))
  headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
  hash := sha256.Sum256(headers)
  b.Hash = hash[:]
}

func NewBlock(data string, preHashBlock []byte) *Block {
  block := &Block{
    TimeStamp:     time.Now().Unix(),
    PrevBlockHash: preHashBlock,
    Data:          []byte(data),
    Hash:          []byte{},
    Nonce:         0,
  }
  pow := NewProofOfWork(block)
  nonce, hash := pow.Run()

  block.Hash = hash[:]
  block.Nonce = nonce
  return block
}

func NewGenesisBlock() *Block {
  return NewBlock("Genesis Block", []byte{})
}

func (b *Block) Serialize() []byte {
  var result bytes.Buffer
  encoder := gob.NewEncoder(&result)
  if err := encoder.Encode(b); err != nil {
    log.Fatal(err)
  }
  return result.Bytes()
}

func Deserialize(d []byte) *Block {
  var block Block
  decoder := gob.NewDecoder(bytes.NewReader(d))
  if err := decoder.Decode(&block); err != nil {
    log.Fatal(err)
  }
  return &block
}
