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
	PrevId        string `json:"prevId"`
	PrevTradeDate string `json:"prevDatetime"`
	Id            string `json:"id"`
	TradeDate     string `json:"Datetime"`
	ImgVector1    string `json:"ImgVector1"`
	ImgVector2    string `json:"ImgVector2"`
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
