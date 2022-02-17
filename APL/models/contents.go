package models

import (
	"crypto/sha256"
	"github.com/cbergoon/merkletree"
)

// Content 머클 트리 안에 들어갈 값
type Content struct {
	Value string
}

func (c Content) CalculateHash() ([]byte, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(c.Value)); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func (c Content) Equals(other merkletree.Content) (bool, error) {
	return c.Value == other.(Content).Value, nil
}
