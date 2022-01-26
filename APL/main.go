package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Req struct {
	Id      string `json:"id"`
	Invoice string `json:"invoiceNumber"`
	Data    []Data `json:"data"`
}

type Data struct {
	Datetime string  `json:"datetime"`
	Temp     float32 `json:"temperature"`
}

type targets struct {
	Coldchain []Coldchain `json:"target"`
}

type Coldchain struct {
	Id       string `json:"id"`
	Datetime string `json:"datetime"`
	Temp     string `json:"temperature"`
}

func SearchEx(c *gin.Context) {
	var target targets
	if err := c.BindJSON(&target); err != nil {
		fmt.Println(err.Error())
	}
	c.IndentedJSON(http.StatusOK, target)
}

func InsertEx(c *gin.Context) {

	var coldchain Coldchain
	if err := c.BindJSON(&coldchain); err != nil {
		fmt.Println(err.Error())
	}
	c.IndentedJSON(http.StatusOK, coldchain)
}

var id = make(map[string]string) // 데이터 저장소
var invoice = make(map[string]bool)

func Insert(c *gin.Context) {
	var request Req
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err.Error())
	}
	var found bool
	_, foundid := id[request.Id]
	foundinvoicde := invoice[request.Invoice]
	if !foundid && !foundinvoicde {
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
	var request Req
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

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "홈 화면 입니다",
	})
}

func main() {
	r := gin.Default()
	r.GET("/", Home)
	r.POST("/insertex", InsertEx)
	r.POST("/searchex", SearchEx)
	r.POST("/insert", Insert)
	r.POST("/search", Search)
	r.Run(":8080")
}

/*
참고 자료 : https://github.com/gin-gonic/gin/issues/715
https://velog.io/@soosungp33/golang-Gin
*/
