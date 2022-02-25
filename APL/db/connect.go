package db

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func Connect() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("파일 로딩 에러 (.env) From db_connect.go")
	}
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	PROTOCOL := "tcp(localhost:3306)" // 로컬 환경은 로컬호스트의 3306 포트
	DBNAME := os.Getenv("DBNAME")
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	print(CONNECT)
	db, err := gorm.Open(mysql.Open(CONNECT), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	return db
}
