package helper

import (
	"Golang-Web/APL/db"
	"Golang-Web/APL/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64 // 엑세스 토큰 유효기간
	RtExpires    int64 // Refresh Token 유효기간
}

func CreateToken(userId uint64) (*TokenDetails, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 1).Unix()    // 원래는 15분인데, 실험용으로 1분으로 설정.
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // 7일
	var err error
	// Access Token 만들기
	// ENV 에 ACCESS_SECRET 에 담긴 값을 이용하여 JWT 서명

	td.AccessToken, err = createAccessToken(userId, td, err)
	if err != nil {
		return nil, err
	}

	td.RefreshToken, err = createRefreshToken(userId, td, err)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func createRefreshToken(userId uint64, td *TokenDetails, err error) (string, error) {
	rtClaims := jwt.MapClaims{}
	rtClaims["userid"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	return rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
}

func createAccessToken(userId uint64, td *TokenDetails, err error) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userid"] = userId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
}

// ExtractToken 헤더에서 토큰 추출하는 함수.
func ExtractToken(r *http.Request) string {
	// bearToken := r.Header.Get("access-token")
	cookie, err := r.Cookie("access-token")
	if err != nil {
		return ""
	}
	Token := cookie.Value
	//normally Authorization the_token_xxx
	log.Print(Token) // 그냥 로그용
	if len(Token) == 0 {
		return ""
	}
	return Token
}

// VerifyAccessToken 토큰을 가져와서 Signing Method 검증.
func VerifyAccessToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Verify Access Token 1번 에러")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		log.Println("Verify Access Token 2번 에러")
		return nil, err
	}
	return token, nil
}

// AccessTokenValid 유효한 토큰인지 검증.
func AccessTokenValid(r *http.Request) error {
	token, err := VerifyAccessToken(r)
	if err != nil {
		log.Println("Access Token Valid 1번 에러")
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		log.Println("Access Token Valid 2번 에러")

		return err
	}
	return nil
}

// CheckToken 어떤 작업 할때 마다 이 코드를 넣어줘야 한다.
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
		c.JSON(http.StatusUnprocessableEntity, "엑세스 토큰이 없습니다.")
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
