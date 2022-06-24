package main

import (
	"Golang-Web/APL/models"
	"crypto/sha256"
	"encoding/hex"
	"github.com/cbergoon/merkletree"
	"log"
)

type Content struct {
	Id        string // 사용자 아이디
	Percent   string // 지분
	TradeDate string // 거래 날짜
	ImgVector string
}

func (c Content) CalculateHash() ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(c.Id))
	_, err = hash.Write([]byte(c.Percent))
	_, err = hash.Write([]byte(c.TradeDate))
	_, err = hash.Write([]byte(c.ImgVector))
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func (c Content) Equals(other merkletree.Content) (bool, error) {
	return c.Id == other.(Content).Id && c.Percent == other.(Content).Percent && c.TradeDate == other.(Content).TradeDate && c.ImgVector == other.(Content).ImgVector, nil
}

// createContent 머클 트리에 들어갈 항목 리스트를 request 를 넘기면 자동으로 만들어서 반환.
func createContent(requests []models.MerkleReq) []merkletree.Content {
	var list []merkletree.Content
	// 요청이 몇개 들어올 지 모르므로 배열로 만든다.
	for _, request := range requests {
		list = append(list, models.Content{
			Id:        request.Id,
			Percent:   request.Percent,
			TradeDate: request.TradeDate,
			ImgVector: request.ImgVector,
		})
	}
	return list
}
func main() {
	var list []merkletree.Content
	// 요청이 몇개 들어올 지 모르므로 배열로 만든다.

	list = append(list, models.Content{
		Id:        "5",
		Percent:   "23.4",
		TradeDate: "20220402",
		ImgVector: "asdqwewqe1212512512321512",
	})

	list = append(list, models.Content{
		Id:        "21",
		Percent:   "23.4",
		TradeDate: "20220402",
		ImgVector: "asdqwewqe1212512512321512",
	})
	
	merkleTree, err := merkletree.NewTree(list)
	log.Println(merkleTree.Leafs[0])
	if err != nil {
		log.Fatal(err)
	} else {
		merkleRoot := merkleTree.MerkleRoot()
		log.Println("Merkel Root : ", merkleRoot)
		// MerkleRoot 값을 해시 값으로 변환.
		log.Println("Hash : ", hex.EncodeToString(merkleRoot))
	}

	verifyTree, err := merkleTree.VerifyTree()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Verify Tree : ", verifyTree)
	}

}
