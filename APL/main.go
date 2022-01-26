package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type targets struct {
	Coldchain []Coldchain `json:"target"`
}

type Coldchain struct {
	Id       string `json:"id"`
	Datetime string `json:"datetime"`
	Temp     string `json:"temperature"`
}

func Search(c *gin.Context) {
	var target targets
	if err := c.BindJSON(&target); err != nil {
		fmt.Println(err.Error())
	}
	c.IndentedJSON(http.StatusOK, target)
}

func Insert(c *gin.Context) {

	var coldchain Coldchain
	if err := c.BindJSON(&coldchain); err != nil {
		fmt.Println(err.Error())
	}
	c.IndentedJSON(http.StatusOK, coldchain)
	/*
		body := c.Request.Body
			value, err := ioutil.ReadAll(body)
			if err != nil {
				fmt.Println(err.Error())
			}
		var data map[string]interface{}
		json.Unmarshal([]byte(value), &data) // JSON을 Golang 자료형으로
		c.JSON(http.StatusOK, gin.H{
			"id":          data["id"],
			"datetime":    data["datetime"],
			"temperature": data["temperature"],
		})

		doc, _ := json.Marshal(data) // Go 자료형을 JSON으로
		c.String(http.StatusOK, string(doc))
	*/
}

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "홈 화면 입니다",
	})
}

func main() {
	r := gin.Default()
	r.GET("/", Home)
	r.POST("/insert", Insert)
	r.POST("/search", Search)
	r.Run(":8080")
}
