package handler

import (
	"Golang-Web/APL/db"
	"Golang-Web/APL/helper"
	"Golang-Web/APL/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

// 예시 유저.
var user = models.Users{
	ID:       1,
	Username: "username",
	Password: "password",
}

var ACCESS_SECRET = viper.GetString(`token.ACCESS_SECRET`)

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

func Login(c *gin.Context) {
	user := new(models.Users)
	if err := c.Bind(user); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	// DB 연결
	db := db.Connect()
	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	defer sqlDB.Close()

	inputPassword := user.Password
	result := db.Find(&user, "username=?", user.Username)

	// username 존재하지 않는 경우.
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, "존재하지 않는 유저 아이디 입니다.")
		return
	}

	// 비밀번호가 틀린 경우
	checkHash := helper.CheckPasswordHash(user.Password, inputPassword)
	if checkHash == false {
		c.JSON(http.StatusBadRequest, "비밀번호가 틀렸습니다.")
		return
	}

	// 저 두가지 경우가 아니면 성공
	td, err := helper.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "JWT 토큰 생성 실패")
	}

	// 쿠키 저장.
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access-token",
		Value:    td.AccessToken,
		Expires:  td.AtExpires,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, gin.H{
		"Message":       "로그인 성공",
		"Access token":  td.AccessToken,
		"Refresh token": td.RefreshToken,
	})

	// 로그인 성공하면 Refresh Token 을 DB 에 저장.
	db.Model(&user).Update("refresh_token", td.RefreshToken)
	return
}

func VerifyAccessToken(c *gin.Context) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}
	accessToken := c.GetHeader("access-token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, "엑세스 토큰이 없습니다.")
		return
	}

	claims := jwt.MapClaims{}

	// 토큰 decode. claims 에 복호화한 정보 저장.
	_, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, "토큰 인증 실패. 재발급 받으세요.")
		return
	}

	c.JSON(http.StatusOK, "엑세스 토큰 인증 완료")
	return
}

func RecreateToken(c *gin.Context) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}

	at := c.GetHeader("access-token")
	if at == "" {
		c.JSON(http.StatusUnauthorized, "엑세스 토큰이 없습니다.")
		return
	}

	claims := jwt.MapClaims{}

	// 토큰 decode. claims 에 복호화한 정보 저장.
	_, err := jwt.ParseWithClaims(at, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, "토큰 인증 실패. 재발급 받으세요.")
		return
	}
	userId := claims["userid"].(uint64)
	accessToken, err := helper.CreateToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "엑세스 토큰 재생성중 에러")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "토큰 재생성 완료",
		"accessToken": accessToken,
	})
	return
}
func CheckToken(c *gin.Context) {
	user := new(models.Users)
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}
	// DB 연결
	db := db.Connect()
	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	defer sqlDB.Close()

	// 헤더에서 갖고오는 코드 c.Request.Header.Get("Authorization")
	accessToken := c.GetHeader("access-token")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, "엑세스 토큰이 없습니다.")
		return
	}

	claims := jwt.MapClaims{}

	// 토큰 decode. claims 에 복호화한 정보 저장.
	_, err = jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, "토큰 복호화 실패")
		return
	}

	userId := claims["userid"]
	result := db.Find(&user, "id=?", userId)
	// userId 존재하지 않는 경우.
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, "존재하지 않는 유저입니다.")
		return
	}
	refreshToken := user.RefreshToken
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, "Refresh Token Error")
		return
	}
	c.JSON(http.StatusOK, "토큰 검증 완료")
	return
}

/*
참고 자료 : https://learn.vonage.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr/
https://covenant.tistory.com/203
*/
