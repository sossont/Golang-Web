package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"log"
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

// EncryptPk 암호화 하려는 문장과 PK를 넣으면 암호화 해서 반환한다.
func EncryptPk(text string, key *rsa.PrivateKey) string {
	publicKey := &key.PublicKey

	encText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(text))
	if err != nil {
		log.Fatal("Error : ", err)
	}

	encTextStr := base64.StdEncoding.EncodeToString(encText)
	return encTextStr
}

// DecryptPk 복호화 하려는 문장과 PK를 넣으면 복호화 해서 반환해준다.
func DecryptPk(encTextStr string, privateKey *rsa.PrivateKey) string {
	DecodeStr, _ := base64.StdEncoding.DecodeString(encTextStr)
	decText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, DecodeStr)
	if err != nil {
		log.Fatal("Error : ", err)
	}

	return string(decText)
}
