package core

import (
	"crypto/sha256"

	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
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
func (tree *MerkleTree) FindPath(index int) []types.MerklePath {
	node := tree.RootNode
	var nodes []*MerkleNode

	for node.l != node.r {
		nodes = append(nodes, node)
		if index < (node.l+node.r)/2 {
			node = node.Left
		} else {
			node = node.Right
		}
	}

	var ret []types.MerklePath
	for i := len(nodes) - 1; i >= 0; i-- {
		if index > nodes[i].Left.r {
			ret = append(ret, types.MerklePath{
				Left:      true,
				HashValue: nodes[i].Left.Data,
			})
		}
		if index < nodes[i].Right.l {
			ret = append(ret, types.MerklePath{
				Left:      false,
				HashValue: nodes[i].Right.Data,
			})
		}
	}
	return ret
}

func NewTxnMerkleTree(txns []*types.Transaction) *MerkleTree {
	var txnsBytes [][]byte
	for _, txn := range txns {
		txnsBytes = append(txnsBytes, utils.Encode(txn))
	}
	return NewMerkleTree(txnsBytes)
}

func NewBlockMerkleTree(blocks []*types.Block) *MerkleTree {
	var blocksBytes [][]byte
	for _, block := range blocks {
		blocksBytes = append(blocksBytes, block.BlockHeader.Hash())
	}
	return NewMerkleTree(blocksBytes)
}
