package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"time"
)

const (
	dif       = 20
	INT64_MAX = math.MaxInt64
)

type Block struct {
	PrevHash  []byte
	Hash      []byte
	Data      string
	Height    int64
	Timestamp int64
	Nonce     int
}

type BlockChain struct {
	Blocks []Block
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

func ProofOfWork(b Block, dif int) ([]byte, int) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-dif))
	nonce := 0
	for ; nonce < INT64_MAX; nonce++ {
		check := bytes.Join(
			[][]byte{
				b.PrevHash,
				[]byte(b.Data),
				IntToHex(b.Height),
				IntToHex(b.Timestamp),
				IntToHex(int64(nonce)),
			},
			[]byte{})
		hash := sha256.Sum256(check)
		var hashInt big.Int
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(target) == -1 {
			return hash[:], nonce
		}
	}
	return []byte(""), nonce
}

func GenesisBlock(data string) *BlockChain {
	var bc BlockChain
	bc.Blocks = make([]Block, 1)
	bc.Blocks[0] = Block{
		PrevHash:  []byte(""),
		Data:      data,
		Height:    1,
		Timestamp: time.Now().Unix(),
	}
	bc.Blocks[0].Hash, bc.Blocks[0].Nonce = ProofOfWork(bc.Blocks[0], dif)
	return &bc
}

func GenerateBlock(bc *BlockChain, data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	block := Block{
		PrevHash:  prevBlock.Hash,
		Data:      data,
		Height:    prevBlock.Height + 1,
		Timestamp: time.Now().Unix(),
	}
	block.Hash, block.Nonce = ProofOfWork(block, dif)
	bc.Blocks = append(bc.Blocks, block)
}

func Print(bc *BlockChain) {
	for _, i := range bc.Blocks {
		fmt.Printf("PrevHash: %x\n", i.PrevHash)
		fmt.Printf("Hash: %x\n", i.Hash)
		fmt.Println("Block's Data: ", i.Data)
		fmt.Println("Current Height: ", i.Height)
		fmt.Println("Timestamp: ", i.Timestamp)
		fmt.Println("Nonce: ", i.Nonce)
	}
}

func main() {
	blockchain := GenesisBlock("i am reigns")
	GenerateBlock(blockchain, "send 2$ to alice")
	GenerateBlock(blockchain, "send 3$ to alice")
	Print(blockchain)
}
