package models

import (
	"crypto/sha256"
	"encoding/hex"
)

// TODO: all nums are uint32 so range = 0 > 2^32 | max = (str 20B * limit 1000) < 512B redis max entry
type Fizzbuzz struct {
	ModA     uint16 `json:"mod_a" validate:"required,numeric,max=10000"`
	ModB     uint16 `json:"mod_b" validate:"required,numeric,max=10000"`
	Limit    uint16 `json:"limit" validate:"required,numeric,max=10000"`
	ReplaceA string `json:"replace_a" validate:"required,alphanum,max=20"`
	ReplaceB string `json:"replace_b" validate:"required,alphanum"`
}

type Result struct {
	Hash     string    `json:"hash"`
	Fizzbuzz *Fizzbuzz `json:"fizzbuzz"`
	Result   []string  `json:"result"`
	State    string    `json:"state"`
}

func (fb *Fizzbuzz) HashData() string {
	record := string(fb.ModA) + string(fb.ModB) + string(fb.Limit) + fb.ReplaceA + fb.ReplaceB
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)

	return hex.EncodeToString(hashed)
}
