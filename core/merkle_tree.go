package core

import (
	"crypto/sha256"

	"github.com/YouDad/blockchain/types"
)

// MerkleTree represent a Merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode represent a Merkle tree node
type MerkleNode struct {
	l     int
	r     int
	Left  *MerkleNode
	Right *MerkleNode
	Data  types.HashValue
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
func NewMerkleTree(dataSeq [][]byte) *MerkleTree {
	var nodes []MerkleNode

	for index, data := range dataSeq {
		node := NewMerkleNode(index, index, nil, nil, data)
		nodes = append(nodes, *node)
	}

	for len(nodes) != 1 {
		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			if j == len(nodes)-1 {
				nodes = append(nodes, nodes[j])
				last := nodes[len(nodes)-1]
				delta := last.r - last.l + 1
				last.l += delta
				last.r += delta
			}
			node := NewMerkleNode(nodes[j].l, nodes[j+1].r, &nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		nodes = newLevel
	}

	mTree := MerkleTree{&nodes[0]}

	return &mTree
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(l, r int, left, right *MerkleNode, data []byte) *MerkleNode {
	if data == nil {
		data = append(left.Data, right.Data...)
	}
	hash := sha256.Sum256(data)
	return &MerkleNode{l, r, left, right, hash[:]}
}

// Index start from 0
func (tree *MerkleTree) FindPath(index int) []types.HashValue {
	node := tree.RootNode
	var nodes []*MerkleNode

	for node.l != node.r {
		nodes = append(nodes, node)
		mid := node.l + node.r
		if index <= mid/2 {
			node = node.Left
		} else {
			node = node.Right
		}
	}

	var ret []types.HashValue
	for i := len(nodes) - 1; i > 0; i-- {
		ret = append(ret, nodes[i].Data)
	}
	return ret
}
