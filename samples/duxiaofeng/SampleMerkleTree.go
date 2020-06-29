package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

type MerkleTreeNode struct {
	Data  []byte
	Hash  []byte
	Left  *MerkleTreeNode
	Right *MerkleTreeNode
}

type MerkleTree struct {
	Root *MerkleTreeNode
}

type Queue struct {
	List []MerkleTreeNode
}

func NewQueue() *Queue {
	list := make([]MerkleTreeNode, 0)
	return &Queue{list}
}

func (q *Queue) Push(node *MerkleTreeNode) {
	q.List = append(q.List, *node)
}

func (q *Queue) Pop() *MerkleTreeNode {
	if q.Len() == 0 {
		panic("Empty!")
	}
	node := q.List[0]
	q.List = q.List[1:]
	return &node
}

func (q *Queue) Len() int {
	return len(q.List)
}

func NewMerkleTreeNode(left, right *MerkleTreeNode, data []byte) *MerkleTreeNode {
	var node MerkleTreeNode

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Hash = hash[:]
	} else {
		childHash := append(left.Hash, right.Hash...)
		hash := sha256.Sum256(childHash)
		node.Hash = hash[:]
	}

	node.Left = left
	node.Right = right
	node.Data = data
	return &node
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var root MerkleTree
	nodes := NewQueue()
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, i := range data {
		nodes.Push(NewMerkleTreeNode(nil, nil, i))
	}

	for nodes.Len() > 1 {
		left := nodes.Pop()
		right := nodes.Pop()
		node := NewMerkleTreeNode(left, right, []byte(""))
		nodes.Push(node)
	}
	root.Root = nodes.Pop()
	return &root
}

func PreOrderVisit(root *MerkleTreeNode) {
	if root != nil {
		fmt.Print(root.Data)
		PreOrderVisit(root.Left)
		PreOrderVisit(root.Right)
	}
}

func main() {
	data := make([][]byte, 6)
	for i := 0; i < 6; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	root := NewMerkleTree(data)
	PreOrderVisit(root.Root)
}

//              nil
//        nil           nil
// [52]   [53]  nil       nil
//              [48] [49] [50] [51]
