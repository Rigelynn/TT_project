package service

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// MerkleNode 定义默克尔树的节点结构体
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// MerkleTree 定义默克尔树结构（树形结构）
type MerkleTree struct {
	RootNode *MerkleNode //只记录一个根节点，任意一个根节点都可以追溯全部节点。

}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}
	//如果为左右节点空，那么就说明他是原始数据节点（叶子节点）。
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		//将[32]byte转换为byte[]
		mNode.Data = hash[:]
	}
	//赋值
	mNode.Left = left
	mNode.Right = right
	return &mNode
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	//定义一个结构体类型的切片
	var nodes []MerkleNode
	//确保节点为2的整数倍
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
	//循环嵌套完成节点树形构造
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)

		}
		nodes = newLevel

	}
	//构造默克尔树
	mTree := MerkleTree{&nodes[0]}
	return &mTree
}

func showMerkleTree(root *MerkleNode) {
	if root == nil {
		return
	} else {
		PrintNode(root)
	}
	showMerkleTree(root.Left)
	showMerkleTree(root.Right)
}

func PrintNode(node *MerkleNode) {
	fmt.Printf("%p\n", node)
	if node != nil {
		fmt.Printf("left[%p],right[%p],data(%x)\n", node.Left, node.Right, node.Data)
	}

}
func check(node *MerkleNode) bool {
	if node.Left == nil {
		return true
	}
	prevHashes := append(node.Left.Data, node.Right.Data...)
	hash32 := sha256.Sum256(prevHashes)
	hash := hash32[:]
	return bytes.Compare(hash, node.Data) == 0

}
