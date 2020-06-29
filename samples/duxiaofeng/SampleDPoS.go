package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var Share map[string]float32
var Total float32
var Candidate []string

type Coin struct {
	Time    int64
	Num     float32
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
}

type BlockChain struct {
	Blocks []Block
}

func init() {
	Candidate = make([]string, 10)
	Share = make(map[string]float32)
	Total = float32(1)
	for i := 0; i < 10; i++ {
		Candidate[i] = "192.168.1." + strconv.Itoa(i)
		Share[Candidate[i]] = float32(1)
		Total++
	}
}

func Vote() string {
	total := len(Candidate)
	rand.Seed(time.Now().UnixNano())
	winner := Candidate[rand.Intn(int(total))]
	return winner
}

func Dividend(coin float32) {
	Candidate = make([]string, 0)
	for k, _ := range Share {
		Share[k] += float32(Share[k]/Total) * coin
		for i := 0; i < int(Share[k]); i++ {
			Candidate = append(Candidate, k)
			Total++
		}
	}
}

func GenesisBlock(data string) *BlockChain {
	addr := Vote()
	var bc BlockChain
	bc.Blocks = make([]Block, 1)
	newCoin := Coin{
		Time:    time.Now().Unix(),
		Num:     float32(1 + rand.Intn(10)),
		Address: addr,
	}
	hash := sha256.Sum256([]byte(data))
	bc.Blocks[0] = Block{
		PrevHash:  []byte(""),
		Hash:      hash[:],
		Data:      data,
		Height:    1,
		Timestamp: time.Now().Unix(),
		Coin:      newCoin,
	}
	Dividend(newCoin.Num)
	CoinPool = append(CoinPool, newCoin)
	return &bc
}

func GenerateBlock(bc *BlockChain, data string) {
	addr := Vote()
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newCoin := Coin{
		Time:    time.Now().Unix(),
		Num:     float32(1 + rand.Intn(10)),
		Address: addr,
	}
	hash := sha256.Sum256([]byte(data))
	b := Block{
		PrevHash:  prevBlock.Hash,
		Hash:      hash[:],
		Data:      data,
		Height:    prevBlock.Height + 1,
		Timestamp: time.Now().Unix(),
		Coin:      newCoin,
	}
	bc.Blocks = append(bc.Blocks, b)
	Dividend(newCoin.Num)
	CoinPool = append(CoinPool, newCoin)
}

func Print(bc *BlockChain) {
	for _, i := range bc.Blocks {
		fmt.Printf("PrevHash: %x\n", i.PrevHash)
		fmt.Printf("Hash: %x\n", i.Hash)
		fmt.Println("Block's Data: ", i.Data)
		fmt.Println("Current Height: ", i.Height)
		fmt.Println("Timestamp: ", i.Timestamp)
		fmt.Println("Address: ", i.Coin.Address)
	}
}

func PrintCoinPool() {
	for _, i := range CoinPool {
		fmt.Println("Coin's Num: ", i.Num)
		fmt.Println("Coin's Time: ", i.Time)
		fmt.Println("Coin's Owner: ", i.Address)
	}
}

func PrintShare() {
	for k, v := range Share {
		fmt.Println(k, ": ", v)
	}
}

func main() {
	bc := GenesisBlock("genesis block")
	GenerateBlock(bc, "send 1$ to alice")
	GenerateBlock(bc, "send 2$ to alice")
	GenerateBlock(bc, "send 3$ to alice")
	Print(bc)
	PrintShare()
	fmt.Println(Candidate)
	fmt.Println(len(Candidate))
}
