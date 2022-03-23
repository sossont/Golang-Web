package helper

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"log"
	"os"
)

// 이 helper 를 사용하는 이유, import 를 간결하게 하고 함수 하나만 사용함으로써 가독성을 높이기 위함.

// GeneratePk Private Key 를 생성해서 반환해준다. 여기 안에 공개 키도 들어있다.
func GeneratePk() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Error : ", err)
	}
	return privateKey
}

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

// Encrypt 암호화 하려는 문장과 공개키를 넣으면 공개키로 암호화 해서 반환한다.
func Encrypt(text string, publicKey *rsa.PublicKey, label []byte) string {
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(text), label)
	// cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(text))
	// 공개키 암호화 이외에는 EncryptOAEP 함수를 사용하는 게 더 안전하다고 한다. 출처 : 공식문서 https://pkg.go.dev/crypto/rsa
	if err != nil {
		log.Fatal("Error : ", err)
	}

	// cipherTextStr := base64.StdEncoding.EncodeToString(cipherText)
	return hex.EncodeToString(cipherText)
}

// Decrypt 복호화 하려는 문장과 PK를 넣으면 복호화 해서 반환해준다.
func Decrypt(text string, privateKey *rsa.PrivateKey, label []byte) string {
	cipherText, err := hex.DecodeString(text)
	// decText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, DecodeStr)
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherText, label)
	if err != nil {
		log.Fatal("Error : ", err)
	}

	return string(plainText)
}

/*
// Encrypt 암호화 하려는 문장과 공개키를 넣으면 공개키로 암호화 해서 반환한다.
func Encrypt(text string, key *rsa.PrivateKey) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}

	publicKey := &key.PublicKey
	label := []byte(os.Getenv("LABEL"))
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(text), label)
	// cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(text))
	// 공개키 암호화 이외에는 EncryptOAEP 함수를 사용하는 게 더 안전하다고 한다. 출처 : 공식문서 https://pkg.go.dev/crypto/rsa
	if err != nil {
		log.Fatal("Error : ", err)
	}

	cipherTextStr := base64.StdEncoding.EncodeToString(cipherText)
	return cipherTextStr
}

// Decrypt 복호화 하려는 문장과 PK를 넣으면 복호화 해서 반환해준다.
func Decrypt(text string, privateKey *rsa.PrivateKey) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("파일 로딩 에러 (.env) ")
	}
	label := []byte(os.Getenv("LABEL"))
	cipherText, err := hex.DecodeString(text)
	// decText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, DecodeStr)
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherText, label)
	if err != nil {
		log.Fatal("Error : ", err)
	}

	return string(plainText)
}

*/
