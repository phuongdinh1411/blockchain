package core

import (
  "bytes"
  "crypto/sha256"
  "encoding/binary"
  "fmt"
  "math"
  "math/big"
)

// Use to validate a block
type ProofOfWork struct {
  block  *Block
  target *big.Int
}

const targetBits = 10

func NewProofOfWork(b *Block) *ProofOfWork {
  target := big.NewInt(1)
  target.Lsh(target, uint(256-targetBits))
  pow := &ProofOfWork{b, target}
  return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
  data := bytes.Join(
    [][]byte{
      pow.block.PrevBlockHash,
      pow.block.Data,
      IntToHex(pow.block.TimeStamp),
      IntToHex(int64(targetBits)),
      IntToHex(int64(nonce)),
    },
    []byte{},
  )
  return data
}

func IntToHex(num int64) []byte {
  buff := new(bytes.Buffer)
  _ = binary.Write(buff, binary.BigEndian, num)
  return buff.Bytes()
}

func (pow *ProofOfWork) Run() (int, []byte) {
  maxNonce := math.MaxInt64
  nonce := 0
  var hashInt big.Int
  var hash [32]byte

  fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
  for nonce < maxNonce {
    data := pow.prepareData(nonce)
    hash = sha256.Sum256(data)
    fmt.Printf("\r%x", hash)
    hashInt.SetBytes(hash[:])

    if hashInt.Cmp(pow.target) == -1 {
      break
    } else {
      nonce += 1
    }
  }
  fmt.Print("\n\n")
  return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
  data := pow.prepareData(pow.block.Nonce)
  hash := sha256.Sum256(data)
  var hashInt big.Int
  hashInt.SetBytes(hash[:])

  isValid := hashInt.Cmp(pow.target) == -1
  return isValid
}
