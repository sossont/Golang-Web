package handler

import (
	"Golang-Web/APL/db"
	"Golang-Web/APL/helper"
	"Golang-Web/APL/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64 // 엑세스 토큰 유효기간
	RtExpires    int64 // Refresh Token
}

// 예시 유저.
var user = models.Users{
	ID:       1,
	Username: "username",
	Password: "password",
}

var ACCESS_SECRET = viper.GetString(`token.ACCESS_SECRET`)

func Login(c *gin.Context) {
	// DB 연결
	db := db.Connect()
	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	defer sqlDB.Close()

	var mockUser models.Users

	if err := c.ShouldBindJSON(&mockUser); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "JSON이 잘못 되었습니다.")
		return
	}

	if user.Username != mockUser.Username || user.Password != mockUser.Password {
		c.JSON(http.StatusUnauthorized, "유저 정보가 틀렸습니다.")
		return
	}

	token, err := CreateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, token)
}

func SignUp(c *gin.Context) {
	user := new(models.Users)
	if err := c.Bind(user); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}
	db := db.Connect()
	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	defer sqlDB.Close()

	result := db.Find(&user, "username=?", user.Username)

	// username 이 이미 존재하는 경우
	if result.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, "이미 존재하는 유저 아이디 입니다.")
		return
	}

	hashPw, err := helper.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "해쉬 값 오류")
		return
	}
	user.Password = hashPw
	if err := db.Create(&user); err.Error != nil {
		c.JSON(http.StatusInternalServerError, "가입 실패")
		return
	}

	c.JSON(http.StatusOK, "유저 가입 성공")
	return
}

func CreateToken(userId uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // 7일
	var err error
	// Access Token 만들기
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	// ENV 에 ACCESS_SECRET 에 담긴 값을 이용하여 JWT 서명

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userid"] = userId
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err

	}

	os.Setenv("REFRESH_TOKEN", "mcmvasdqwer")
	rtClaims := jwt.MapClaims{}
	rtClaims["userid"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_TOKEN")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

/*
참고 자료 : https://learn.vonage.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr/
https://covenant.tistory.com/203
*/
