package handler

import (
	"Golang-Web/APL/helper"
	"Golang-Web/APL/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var id = make(map[string]string) // 데이터 저장소
var invoice = make(map[string]bool)

func Insert(c *gin.Context) {
	var request models.Req
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err.Error())
	}
	var found bool
	_, foundid := id[request.Id]
	foundInvoice := invoice[request.Invoice]
	if !foundid && !foundInvoice {
		found = false
	} else {
		found = true
	}
	// 데이터가 존재하지 않으면 생성
	if found == false {
		id[request.Id] = request.Invoice
		invoice[request.Invoice] = true
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            request.Id,
		"invoiceNumber": request.Invoice,
		"result":        !found,
	})

}

func Search(c *gin.Context) {
	var request models.Req
	c.JSON(http.StatusOK, request.Invoice)
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err.Error())
	}

	found := false
	// 현재 위변조 확인은 안되므로 id, invoice 만 따진다.
	if id[request.Id] == request.Invoice {
		found = true
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            request.Id,
		"invoiceNumber": request.Invoice,
		"validate":      found,
	})
}

func InsertEtcd(c *gin.Context) {
	var request models.Req
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err.Error())
	}
	var found bool
	foundid := helper.GetKV(request.Id)
	foundInvoice := helper.GetKV(request.Invoice)
	if !foundid && !foundInvoice {
		found = false
	} else {
		found = true
	}

	// 데이터가 존재하지 않으면 생성
	if found == false {
		id[request.Id] = request.Invoice
		invoice[request.Invoice] = true
		helper.PutKV(request.Id, request.Invoice) // id, Invoice 각각 저장.
		helper.PutKV(request.Invoice, request.Id)

		// Data 는 json 배열인데, 이것을 string 변환시켜서 저장 해놓는다. 나중에 json 으로 인코딩 할 수 있는 형식으로 저장.
		var data string
		for _, d := range request.Data {
			inp := fmt.Sprintf("{datetime : %s, temperature : %f},", d.Datetime, d.Temp)
			data += inp
		}

		// 저장은 key + /data. (ex : 123123/data , BL2020-OR04R-02/data)
		key1 := request.Id + "/data"
		helper.PutKV(key1, data)
		key2 := request.Invoice + "/data"
		helper.PutKV(key2, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            request.Id,
		"invoiceNumber": request.Invoice,
		"result":        !found,
	})

}

func SearchEtcd(c *gin.Context) {
	var request models.Req
	c.JSON(http.StatusOK, request.Invoice)
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err.Error())
	}

	// 기본값은 false.
	found := false

	var data string
	for _, d := range request.Data {
		inp := fmt.Sprintf("{datetime : %s, temperature : %f},", d.Datetime, d.Temp)
		data += inp
	}
	key1 := request.Id + "/data"
	key2 := request.Invoice + "/data"

	// 검증되었으면 true.
	if helper.GetKV(key1) && helper.GetKV(key2) {
		found = true
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            request.Id,
		"invoiceNumber": request.Invoice,
		"validate":      found,
	})
}
