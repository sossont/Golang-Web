package handler

import (
	"Golang-Web/APL/models"
	"fmt"
	"github.com/cbergoon/merkletree"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateMerkle(c *gin.Context) {
	var request models.MerkleReq

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "머클 트리 Request 오류")
		fmt.Println(err.Error())
		return
	}
	var list []merkletree.Content
	list = append(list, models.Content{Value: request.PrevId})
	list = append(list, models.Content{Value: request.PrevTradeDate})
	list = append(list, models.Content{Value: request.ImgVector1})
	list = append(list, models.Content{Value: request.ImgVector2})
	list = append(list, models.Content{Value: request.Id})
	list = append(list, models.Content{Value: request.TradeDate})
	list = append(list, models.Content{Value: request.ImgVector1})
	list = append(list, models.Content{Value: request.ImgVector2})

	mt, err := merkletree.NewTree(list)
	log.Println(mt)

	if err != nil {
		log.Fatal(err)
	} else {
		merkleRoot := mt.MerkleRoot()
		log.Println("Merkle Root : ", merkleRoot)
	}

	vt, err := mt.VerifyTree()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Verify Tree : ", vt)
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "머클 트리 생성 완료",
	})
	return
}

// VerifyMerkle 이 부분은 머클 트리를 DB에 저장 해야 구현 가능할 듯.
func VerifyMerkle(c *gin.Context) {
	var request models.VerifyMerkleReq

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "머클 트리 Verify Request 오류")
		fmt.Println(err.Error())
		return
	}
}
