package main

import (
	"Golang-Web/APL/handler"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	/*
		r.GET("/", Home)
		r.POST("/insert", Insert)
		r.POST("/search", Search)
		r.Run(":8080")

	*/
	r.GET("/", handler.Home)
	r.POST("/api/insert", handler.Insert)
	r.POST("/api/search", handler.Search)
	r.POST("/api/insertEtcd", handler.InsertEtcd)
	r.POST("/api/searchEtcd", handler.SearchEtcd)
	r.POST("/api/signup", handler.SignUp)
	r.POST("/api/login", handler.Login)
	r.GET("/api/check", handler.CheckToken)
	r.Run(":8080")

}

/*
참고 자료 :
https://github.com/gin-gonic/gin/issues/715
https://velog.io/@soosungp33/golang-Gin
https://github.com/etcd-io/etcd/blob/main/tests/integration/clientv3/examples/example_kv_test.go
*/
