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