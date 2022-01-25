package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	newServer().Run()
}

func newServer() *gin.Engine {
	r := gin.Default()
	r.GET("", Handler)
	r.GET("/name", UserHandler)
	r.POST("/add", addHandler)
	return r
}

type Account struct {
	Id   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func Handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"responseData": "hello world",
	})
}

func UserHandler(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"greetings": fmt.Sprintf("hello %v", name),
	})
}

func addHandler(c *gin.Context) {
	var data Account
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("err: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dataReceived": data,
	})
}
