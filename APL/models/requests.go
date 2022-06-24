package models

type Req struct {
	Id      string `json:"id"`
	Invoice string `json:"invoiceNumber"`
	Data    []Data `json:"data"`
}

type Data struct {
	Datetime string  `json:"datetime"`
	Temp     float32 `json:"temperature"`
}

// MerkleReq Requests For Make Merkle Tree
type MerkleReq struct {
	Id        string `json:"id"`        // 사용자 아이디
	Percent   string `json:"percent"`   // 지분
	TradeDate string `json:"Datetime"`  // 거래 날짜
	ImgVector string `json:"ImgVector"` // 이미지 벡터 값(이미지 아이디 값으로 바꿔도 될 듯)
}

type VerifyMerkleReq struct {
	Value string `json:"value"`
}

/*
type targets struct {
	Coldchain []Coldchain `json:"target"`
}

type Coldchain struct {
	Id       string `json:"id"`
	Datetime string `json:"datetime"`
	Temp     string `json:"temperature"`
}
*/
