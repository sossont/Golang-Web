package helper

import (
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    time.Time // 엑세스 토큰 유효기간
	RtExpires    time.Time // Refresh Token
}

func CreateToken(userId uint64) (*TokenDetails, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15)   // 15분
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7) // 7일
	var err error
	// Access Token 만들기
	// ENV 에 ACCESS_SECRET 에 담긴 값을 이용하여 JWT 서명

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userid"] = userId
	atClaims["exp"] = td.AtExpires.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["userid"] = userId
	rtClaims["exp"] = td.RtExpires.Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_TOKEN")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func CreateAccessToken(userId uint64) (*TokenDetails, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15) // 15분
	var err error
	// Access Token 만들기
	// ENV 에 ACCESS_SECRET 에 담긴 값을 이용하여 JWT 서명

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userid"] = userId
	atClaims["exp"] = td.AtExpires.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}
