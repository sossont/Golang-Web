package models

type Users struct {
	ID           uint64 `json:"id"`            // Id ( PK )
	Password     string `json:"password"`      // 비밀번호
	UserId       string `json:"userid"`        // 유저 아이디
	Username     string `json:"username"`      // 닉네임
	Email        string `json:"email"`         // 이메일
	PhoneNumber  string `json:"phoneNumber"`   // 폰번호
	RefreshToken string `json:"refresh_token"` // 인증을 위한 Refresh Token
}
