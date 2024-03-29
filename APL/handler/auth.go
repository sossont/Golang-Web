package handler

import (
	"Golang-Web/APL/db"
	"Golang-Web/APL/helper"
	"Golang-Web/APL/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func SignUp(c *gin.Context) {
	user := new(models.Users)
	// 양식에 맞지 않은 요청
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

	result := db.Find(&user, "user_id=?", user.UserId)

	// userid 이 이미 존재하는 경우
	if result.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, "이미 존재하는 유저 아이디 입니다.")
		return
	}

	// 넘어온 비밀번호 값 해싱해서 저장
	hashPw, err := helper.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "해쉬 값 오류")
		return
	}
	user.Password = hashPw

	// DB에 저장. 저장 실패 시 오류 반환
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
	result := db.Find(&user, "user_id=?", user.UserId)

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
	td, err := helper.CreateToken(user.UserId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "JWT 토큰 생성 실패")
		return
	}

	// 쿠키 저장.
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access-token",
		Value:    td.AccessToken,
		Expires:  time.Unix(td.AtExpires, 0),
		HttpOnly: true,
	})

	response := map[string]string{
		"Access token":  td.AccessToken,
		"Refresh token": td.RefreshToken,
	}

	c.JSON(http.StatusOK, response)

	// 로그인 성공하면 Refresh Token 을 DB 에 저장.
	db.Model(&user).Update("refresh_token", td.RefreshToken)
	return
}

func Logout(c *gin.Context) {
	err := helper.AccessTokenValid(c.Request)
	if err != nil {
		log.Println("토큰 유효 에러")
		c.JSON(
			http.StatusBadRequest,
			"Token 이 유효하지 않습니다.")
		return
	}

	// Cookie 에서 Access Token 을 지워 버린다.
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access-token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	c.JSON(http.StatusCreated, "로그아웃 완료")
}

// Check : Access Token 유효한지 확인하기 위한 API
func Check(c *gin.Context) {
	err := helper.AccessTokenValid(c.Request)
	if err == nil {
		log.Println("토큰 유효함!")
		c.JSON(http.StatusOK,
			"Token 유효합니다.")
		return
	}

	log.Println("엑세스 토큰 만료")
	userId, err := helper.RefreshAccessToken(c.Request)
	if err != nil {
		log.Println("리프레시 토큰 에러")
		log.Println(err)
		c.JSON(http.StatusBadRequest, "Refresh Token이 만료되었습니다.")
		return
	}

	td, err := helper.CreateToken(userId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "JWT 토큰 생성 실패")
		return
	}

	// 쿠키 저장.
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access-token",
		Value:    td.AccessToken,
		Expires:  time.Unix(td.AtExpires, 0),
		HttpOnly: true,
	})

	response := map[string]string{
		"Message":      "액세스 토큰 재발급 성공",
		"Access token": td.AccessToken,
	}
	c.JSON(http.StatusCreated, response)
}

/*
func Refresh(c *gin.Context) {
	user := new(models.Users)
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	// DB 연결
	db := db.Connect()
	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	defer sqlDB.Close()

	// 이 Refresh Token 을 갖고 있는 유저를 DB 에서 찾는다.

	refreshToken := mapToken["refresh_token"]

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, "리프레시 토큰 만료")
		return
	}

	//토큰 유효한지 확인
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {

		// 토큰 복호화해서 user_id 추출.
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["userid"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token
		result := db.Find(&user, "id=?", userId)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusUnprocessableEntity, "존재하지 않는 유저입니다.")
		}
		db.Model(&user).Update("refresh_token", "") // refresh_token 항목 공백으로 만들기.

		//Create new pairs of refresh and access tokens
		td, err := helper.CreateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "JWT 토큰 생성 실패")
		}

		// 쿠키 저장.
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "access-token",
			Value:    td.AccessToken,
			Expires:  time.Unix(td.AtExpires, 0),
			HttpOnly: true,
		})

		response := map[string]string{
			"Access token":  td.AccessToken,
			"Refresh token": td.RefreshToken,
		}

		c.JSON(http.StatusOK, response)

		// 로그인 성공하면 Refresh Token 을 DB 에 저장.
		db.Model(&user).Update("refresh_token", td.RefreshToken)
		tokens := map[string]string{
			"access_token":  td.AccessToken,
			"refresh_token": td.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
*/

/*
참고 자료 : https://learn.vonage.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr/
https://covenant.tistory.com/203
*/
