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
	r.POST("/insert", handler.Insert)
	r.POST("/search", handler.Search)
	r.POST("/insertEtcd", handler.InsertEtcd)
	r.POST("/searchEtcd", handler.SearchEtcd)
	r.POST("/signup", handler.SignUp)
	r.POST("/login", handler.Login)
	r.Run(":8080")

}

/*
참고 자료 :
https://github.com/gin-gonic/gin/issues/715
https://velog.io/@soosungp33/golang-Gin
https://github.com/etcd-io/etcd/blob/main/tests/integration/clientv3/examples/example_kv_test.go
*/
