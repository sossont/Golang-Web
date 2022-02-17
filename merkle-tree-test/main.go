package main

import (
	"crypto/sha256"
	"github.com/cbergoon/merkletree"
	"log"
)

type TestContent struct {
	x string
}

// CalculateHash TestContent 값의 Hash값 계산.
func (t TestContent) CalculateHash() ([]byte, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(t.x)); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// Equals 두 값이 같은 지 테스트
func (t TestContent) Equals(other merkletree.Content) (bool, error) {
	return t.x == other.(TestContent).x, nil
}

func main() {
	var list []merkletree.Content
	list = append(list, TestContent{x: "Hello"})
	list = append(list, TestContent{x: "Hi"})
	list = append(list, TestContent{x: "Ex1"})
	list = append(list, TestContent{x: "화누정"})

	merkleTree, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	} else {
		merkleRoot := merkleTree.MerkleRoot()
		log.Println("Merkel Root : ", merkleRoot)
	}

	verifyTree, err := merkleTree.VerifyTree()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Verify Tree : ", verifyTree)
	}

	tc := TestContent{x: "Ex1"}
	verifyContent, err := merkleTree.VerifyContent(tc)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Verify Content : ", tc, verifyContent)
	}

	log.Println(merkleTree)
	path, idx, err := merkleTree.GetMerklePath(TestContent{x: "Hi"})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Merkle Path : ", path, idx)
	}
}
