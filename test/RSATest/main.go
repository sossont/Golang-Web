package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func ImportPEM(file *os.File) (*pem.Block, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := fileInfo.Size()
	buf := make([]byte, size)
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(buf)
	keyPem, _ := pem.Decode([]byte(buf))
	return keyPem, nil
}

func PEMtoPrivateKey(PEM []byte) *rsa.PrivateKey {
	privateKey, err := x509.ParsePKCS1PrivateKey(PEM)
	if err != nil {
		log.Fatal("PK 변환 에러")
	}
	return privateKey
}

func PEMtoPublicKey(PEM []byte) *rsa.PublicKey {
	publicKey, err := x509.ParsePKIXPublicKey(PEM)
	if err != nil {
		log.Fatal("PK 변환 에러")
	}

	return publicKey.(*rsa.PublicKey)
}

// EncryptPk 암호화 하려는 문장과 공개키를 넣으면 공개키로 암호화 해서 반환한다.
func EncryptPk(text string, publicKey *rsa.PublicKey, label []byte) string {
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(text), label)
	// cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(text))
	// 공개키 암호화 이외에는 EncryptOAEP 함수를 사용하는 게 더 안전하다고 한다. 출처 : 공식문서 https://pkg.go.dev/crypto/rsa
	if err != nil {
		log.Fatal("Error : ", err)
	}

	// cipherTextStr := base64.StdEncoding.EncodeToString(cipherText)
	return hex.EncodeToString(cipherText)
}

// DecryptPk 복호화 하려는 문장과 PK를 넣으면 복호화 해서 반환해준다.
func DecryptPk(text string, privateKey *rsa.PrivateKey, label []byte) string {
	cipherText, err := hex.DecodeString(text)
	// decText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, DecodeStr)
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherText, label)
	if err != nil {
		log.Fatal("Error : ", err)
	}

	return string(plainText)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}
	label := []byte(os.Getenv("LABEL"))

	// import private_key.pem
	privateKeyFile, err := os.Open("/Users/hwanu/Desktop/Hongik/APL/keypair/private_key.pem")
	if err != nil {
		log.Fatal("File Open Error")
	}
	privateKeyPem, err := ImportPEM(privateKeyFile)
	if err != nil {
		log.Fatal("Import PEM Error")
	}
	privateKeyFile.Close()

	// import public_key.pem
	publicKeyFile, err := os.Open("/Users/hwanu/Desktop/Hongik/APL/keypair/public_key.pem")
	if err != nil {
		log.Fatal("File Open Error")
	}
	publicKeyPem, err := ImportPEM(publicKeyFile)
	if err != nil {
		log.Fatal("Import PEM Error")
	}
	publicKeyFile.Close()

	// EnCrypt and DeCrypt
	text := "Test"
	privateKey := PEMtoPrivateKey(privateKeyPem.Bytes) // 소유자 개인키
	publicKey := PEMtoPublicKey(publicKeyPem.Bytes)    // 소유자 공개키

	cipherText := EncryptPk(text, publicKey, label)       // 암호문
	plainText := DecryptPk(cipherText, privateKey, label) // 평문

	fmt.Println("text : ", text)
	fmt.Println("cipherText : ", cipherText)
	fmt.Println("plainText : ", plainText)
	fmt.Println("is same : ", text == plainText)
}
