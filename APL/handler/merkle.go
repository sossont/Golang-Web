package handler

import (
	"Golang-Web/APL/helper"
	"Golang-Web/APL/models"
	"encoding/hex"
	"fmt"
	"github.com/cbergoon/merkletree"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var memoryRepo = map[string]string{} // 메모리 저장소.

// createContent 머클 트리에 들어갈 항목 리스트를 request 를 넘기면 자동으로 만들어서 반환.
func createContent(requests []models.MerkleReq) []merkletree.Content {
	var list []merkletree.Content
	vector := requests[0].ImgVector
	if vector == "" {
		log.Println("Vector : ", vector)
		return nil
	}

	if _, exists := memoryRepo[vector]; exists {
		return nil
	}

	// 요청이 몇개 들어올 지 모르므로 배열로 만든다.
	for _, request := range requests {

		if request.ImgVector != vector {
			log.Println("Vector : ", vector)

			return nil
		}

		list = append(list, models.Content{
			Id:        request.Id,
			Percent:   request.Percent,
			TradeDate: request.TradeDate,
			ImgVector: request.ImgVector,
		})
	}
	return list
}

func CreateMerkle(c *gin.Context) {
	var requests []models.MerkleReq

	if err := c.BindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, "머클 트리 Request 오류")
		fmt.Println(err.Error())
		return
	}

	list := createContent(requests)
	if list == nil {
		c.JSON(http.StatusBadRequest, "작품 벡터 오류입니다.")
		return
	}

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
	// privateKey := helper.GeneratePk()
	if err := godotenv.Load(); err != nil {
		log.Println(err)
		log.Fatal("파일 로딩 에러 (.env) ")
	}
	label := []byte(os.Getenv("LABEL"))

	privateKeyFile, err := os.Open("/Users/hwanu/Desktop/Hongik/APL/keypair/private_key.pem")
	if err != nil {
		log.Fatal("File Open Error")
	}
	privateKeyPem, err := helper.ImportPEM(privateKeyFile)
	if err != nil {
		log.Fatal("Import PEM Error")
	}
	privateKeyFile.Close()

	// import public_key.pem
	publicKeyFile, err := os.Open("/Users/hwanu/Desktop/Hongik/APL/keypair/public_key.pem")
	if err != nil {
		log.Fatal("File Open Error")
	}
	publicKeyPem, err := helper.ImportPEM(publicKeyFile)
	if err != nil {
		log.Fatal("Import PEM Error")
	}
	publicKeyFile.Close()

	// EnCrypt and DeCrypt
	privateKey := helper.PEMtoPrivateKey(privateKeyPem.Bytes) // 소유자 개인키
	publicKey := helper.PEMtoPublicKey(publicKeyPem.Bytes)    // 소유자 공개키

	cipherText := helper.Encrypt(hashStr, publicKey, label)    // 소유자 공개키로 암호화한 암호문
	plainText := helper.Decrypt(cipherText, privateKey, label) // 소유자 개인키로 복호화한 평문

	log.Println("암호화된 문장 : ", cipherText)
	log.Println("복호화된 문장 : ", plainText)

	memoryRepo[requests[0].ImgVector] = cipherText // 소유자 공개키로 암호화된 암호문 저장

	c.JSON(http.StatusOK, gin.H{
		"Message":    "머클 트리 생성 완료",
		"암호화 된 루트 값": cipherText,
	})
	return
}
