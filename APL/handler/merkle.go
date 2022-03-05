package handler

import (
	"Golang-Web/APL/helper"
	"Golang-Web/APL/models"
	"encoding/hex"
	"fmt"
	"github.com/cbergoon/merkletree"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var dict = map[string]string{} // 메모리 저장소.

// createContent 머클 트리에 들어갈 항목 8개를 request를 넘기면 자동으로 만들어서 반환.
func createContent(request models.MerkleReq) []merkletree.Content {
	var list []merkletree.Content
	list = append(list, models.Content{Value: request.PrevId})
	list = append(list, models.Content{Value: request.PrevTradeDate})
	list = append(list, models.Content{Value: request.ImgVector1})
	list = append(list, models.Content{Value: request.ImgVector2})
	list = append(list, models.Content{Value: request.Id})
	list = append(list, models.Content{Value: request.TradeDate})
	list = append(list, models.Content{Value: request.ImgVector1})
	list = append(list, models.Content{Value: request.ImgVector2})
	return list
}

func CreateMerkle(c *gin.Context) {
	var request models.MerkleReq

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "머클 트리 Request 오류")
		fmt.Println(err.Error())
		return
	}

	id := c.Query("id")
	// 작품 아이디가 안들어온 경우.
	if id == "" {
		c.JSON(http.StatusBadRequest, "작품 id가 필요합니다.")
		return
	}

	list := createContent(request)

	mt, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	}
	merkleRoot := mt.MerkleRoot()
	log.Println("Merkle Root : ", merkleRoot)

	// 머클 트리 해시 스트링 값
	hashStr := hex.EncodeToString(merkleRoot)
	log.Println("Hash Merkle Root : ", hashStr)

	/*
		vt, err := mt.VerifyTree()
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Verify Tree : ", vt)
		}
	*/

	// 암호화 하는 부분.
	privateKey := helper.GeneratePk()
	ret := helper.EncryptPk(hashStr, privateKey)
	log.Println("암호화된 문장 : ", ret)
	decRet := helper.DecryptPk(ret, privateKey)
	log.Println("복호화된 문장 : ", decRet)

	_, exists := dict[id]
	if !exists { // 존재하지 않는 경우
		dict[id] = ret
	} else {
		c.JSON(http.StatusBadRequest, "이미 존재하는 작품입니다.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message":    "머클 트리 생성 완료",
		"암호화 된 루트 값": ret,
	})
	return
}
