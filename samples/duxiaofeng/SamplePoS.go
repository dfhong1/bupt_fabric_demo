package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"
)

const (
	dif         = 2
	INT64_MAX   = math.MaxInt64
	MaxProbably = 255
	MinProbably = 235
	MaxCoinAge  = 10
	Minute      = 60
)

type Coin struct {
	Time    int64
	Num     int
	Address string
}

var CoinPool []Coin

type Block struct {
	PrevHash  []byte
	Hash      []byte
	Data      string
	Height    int64
	Timestamp int64
	Coin      Coin
	Nonce     int
	Dif       int64
}

type BlockChain struct {
	Blocks []Block
}

func init() {
	rand.Seed(time.Now().UnixNano())
	CoinPool = make([]Coin, 0)
}

func GenesisBlock(data string, addr string) *BlockChain {
	var bc BlockChain
	bc.Blocks = make([]Block, 1)
	newCoin := Coin{
		Time:    time.Now().Unix(),
		Num:     1 + rand.Intn(5),
		Address: addr,
	}
	bc.Blocks[0] = Block{
		PrevHash:  []byte(""),
		Data:      data,
		Height:    1,
		Timestamp: time.Now().Unix(),
		Coin:      newCoin,
		Nonce:     0,
	}
	bc.Blocks[0].Hash, bc.Blocks[0].Nonce, bc.Blocks[0].Dif = ProofOfStake(dif, addr, bc.Blocks[0])
	CoinPool = append(CoinPool, newCoin)
	return &bc
}

func GenerateBlock(bc *BlockChain, data string, addr string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newCoin := Coin{
		Time:    time.Now().Unix(),
		Num:     1 + rand.Intn(5),
		Address: addr,
	}
	b := Block{
		PrevHash:  prevBlock.Hash,
		Data:      data,
		Height:    prevBlock.Height + 1,
		Timestamp: time.Now().Unix(),
	}
	b.Hash, b.Nonce, b.Dif = ProofOfStake(dif, addr, b)
	b.Coin = newCoin
	bc.Blocks = append(bc.Blocks, b)
	CoinPool = append(CoinPool, newCoin)
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

func ProofOfStake(dif int, addr string, b Block) ([]byte, int, int64) {
	var coinAge int64
	var realDif int64
	realDif = int64(MinProbably)
	curTime := time.Now().Unix()

	for k, i := range CoinPool {
		if i.Address == addr && i.Time+MaxCoinAge < curTime {
			//币龄增加, 并设置上限
			var curCoinAge int64
			if curTime-i.Time < 3*MaxCoinAge {
				curCoinAge = curTime - i.Time
			} else {
				curCoinAge = 3 * MaxCoinAge
			}
			coinAge += int64(i.Num) * curCoinAge
			//参与挖矿的币龄置为0
			CoinPool[k].Time = curTime
		}
	}

	if realDif+int64(dif)*coinAge/Minute > int64(MaxProbably) {
		realDif = MaxProbably
	} else {
		realDif += int64(dif) * coinAge / Minute
	}

	target := big.NewInt(1)
	target.Lsh(target, uint(realDif))
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
			return hash[:], nonce, 255 - realDif
		}
	}

	return []byte(""), -1, 255 - realDif
}

func Print(bc *BlockChain) {
	for _, i := range bc.Blocks {
		fmt.Printf("PrevHash: %x\n", i.PrevHash)
		fmt.Printf("Hash: %x\n", i.Hash)
		fmt.Println("Block's Data: ", i.Data)
		fmt.Println("Current Height: ", i.Height)
		fmt.Println("Timestamp: ", i.Timestamp)
		fmt.Println("Nonce: ", i.Nonce)
		fmt.Println("Dif: ", i.Dif)
	}
}

func PrintCoinPool() {
	for _, i := range CoinPool {
		fmt.Println("Coin's Num: ", i.Num)
		fmt.Println("Coin's Time: ", i.Time)
		fmt.Println("Coin's Owner: ", i.Address)
	}
}

func main() {
	addr1 := "192.168.1.1"
	addr2 := "192.168.1.2"
	bc := GenesisBlock("reigns", addr1)
	GenerateBlock(bc, "send 1$ to alice", addr1)
	GenerateBlock(bc, "send 1$ to bob", addr1)
	GenerateBlock(bc, "send 2$ to alice", addr1)
	time.Sleep(11 * time.Second)
	GenerateBlock(bc, "send 3$ to alice", addr1)
	GenerateBlock(bc, "send 4$ to alice", addr2)
	Print(bc)
	PrintCoinPool()
}
