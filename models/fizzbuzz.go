package models

import (
	"crypto/sha256"
	"encoding/hex"
)

type Fizzbuzz struct {
	ModA     uint32 `json:"mod_a" validate:"required,numeric"`
	ModB     uint32 `json:"mod_b" validate:"required,numeric"`
	Limit    uint32 `json:"limit" validate:"required,numeric"`
	ReplaceA string `json:"replace_a" validate:"required,alphanum"`
	ReplaceB string `json:"replace_b" validate:"required,alphanum"`
}

func (fb *Fizzbuzz) HashData() string {
	record := string(fb.ModA) + string(fb.ModB) + string(fb.Limit) + fb.ReplaceA + fb.ReplaceB
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)

	return hex.EncodeToString(hashed)
}
