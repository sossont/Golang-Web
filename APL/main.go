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
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err.Error())
	}
	var found bool
	foundid := GetKV(request.Id)
	foundInvoice := GetKV(request.Invoice)
	if !foundid && !foundInvoice {
		found = false
	} else {
		found = true
	}

	// 데이터가 존재하지 않으면 생성
	if found == false {
		id[request.Id] = request.Invoice
		invoice[request.Invoice] = true
		PutKV(request.Id, request.Invoice) // id, Invoice 각각 저장.
		PutKV(request.Invoice, request.Id)

		// Data 는 json 배열인데, 이것을 string 변환시켜서 저장 해놓는다. 나중에 json 으로 인코딩 할 수 있는 형식으로 저장.
		var data string
		for _, d := range request.Data {
			inp := fmt.Sprintf("{datetime : %s, temperature : %f},", d.Datetime, d.Temp)
			data += inp
		}

		// 저장은 key + /data. (ex : 123123/data , BL2020-OR04R-02/data)
		key1 := request.Id + "/data"
		PutKV(key1, data)
		key2 := request.Invoice + "/data"
		PutKV(key2, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            request.Id,
		"invoiceNumber": request.Invoice,
		"result":        !found,
	})

}

func SearchEtcd(c *gin.Context) {
	var request Req
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
	if GetKV(key1) && GetKV(key2) {
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
	/*
		r.GET("/", Home)
		r.POST("/insert", Insert)
		r.POST("/search", Search)
		r.Run(":8080")

	*/
	r.GET("/", Home)
	r.POST("/insert", Insert)
	r.POST("/search", Search)
	r.POST("/insertEtcd", InsertEtcd)
	r.POST("/searchEtcd", SearchEtcd)
	r.Run(":8080")
}

func GetKV(key string) bool {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer cli.Close() // 어떤 에러가 발생하더라도 마지막에 Close 된다.

	ctx, cancel := context.WithTimeout(context.Background(), integration.RequestWaitTimeout)
	resp, err := cli.Get(ctx, key)
	cancel()

	if err != nil {
		log.Fatal(err)
		return false
	}
	// Get 성공하면 Key, Value 쌍을 print 하고 True 반환
	// 그 외에는 False 반환
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		return true
	}
	return false
}

func PutKV(key string, value string) bool {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // 어떤 에러가 발생하더라도 마지막에 Close 된다.
	_, err = cli.Put(context.TODO(), key, value)
	// Put 성공하면 True 반환, 그 외에는 False 반환
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Put 성공\n")
		return true
	}
	return false
}

/*
참고 자료 :
https://github.com/gin-gonic/gin/issues/715
https://velog.io/@soosungp33/golang-Gin
https://github.com/etcd-io/etcd/blob/main/tests/integration/clientv3/examples/example_kv_test.go
*/
