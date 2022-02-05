package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/tests/v3/integration"
	"log"
	"net/http"
	"time"
)

// Global Variable
var id = make(map[string]string) // 데이터 저장소
var invoice = make(map[string]bool)
var url = "http://192.168.0.102:2379/v2/" // 기본 URL

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

func Insert(c *gin.Context) {
	var request Req
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

func InsertEtcd(c *gin.Context) {
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
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	_, err = cli.Put(context.TODO(), "BL2020-OR04R-02", "123123")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), integration.RequestWaitTimeout)
	resp, err := cli.Get(ctx, "BL2020-OR04R-02")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	fmt.Print(resp.Kvs)
	/*
		r.GET("/", Home)
		r.POST("/insert", Insert)
		r.POST("/search", Search)
		r.Run(":8080")

	*/
}

/*
참고 자료 :
https://github.com/gin-gonic/gin/issues/715
https://velog.io/@soosungp33/golang-Gin
*/
