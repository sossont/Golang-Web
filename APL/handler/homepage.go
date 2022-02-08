package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "홈 화면 입니다",
	})
}
