package models

import (
	"crypto/sha256"
	"github.com/cbergoon/merkletree"
)

// Content 머클 트리 안에 들어갈 값
type Content struct {
	Id        string // 사용자 아이디
	Percent   string // 지분
	TradeDate string // 거래 날짜
	ImgVector string
}

func (c Content) CalculateHash() ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(c.Id))
	_, err = hash.Write([]byte(c.Percent))
	_, err = hash.Write([]byte(c.TradeDate))
	_, err = hash.Write([]byte(c.ImgVector))
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func (c Content) Equals(other merkletree.Content) (bool, error) {
	return c.Id == other.(Content).Id && c.Percent == other.(Content).Percent && c.TradeDate == other.(Content).TradeDate && c.ImgVector == other.(Content).ImgVector, nil
}
