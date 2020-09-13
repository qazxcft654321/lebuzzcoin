package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
)

// TODO: all nums are uint32 so range = 0 > 2^32 | max = (str 20B * limit 1000) < 512B redis max entry
type Fizzbuzz struct {
	ModA     uint16 `json:"mod_a" validate:"required,numeric,min=1,max=1000"` // range 1 to 1000
	ModB     uint16 `json:"mod_b" validate:"required,numeric,min=1,max=1000"` // range 1 to 1000
	Limit    uint16 `json:"limit" validate:"required,numeric,min=1,max=1000"` // range 1 to 1000
	ReplaceA string `json:"replace_a" validate:"required,alphanum,max=20"`    // max 20Bytes long
	ReplaceB string `json:"replace_b" validate:"required,alphanum,max=20"`    // max 20Bytes long
}

type Result struct {
	Hash     string    `json:"hash"`
	Fizzbuzz *Fizzbuzz `json:"fizzbuzz"`
	Result   []string  `json:"result"`
	Flag     string    `json:"state"`
}

func (fb *Fizzbuzz) HashData() string {
	record := string(fb.ModA) + string(fb.ModB) + string(fb.Limit) + fb.ReplaceA + fb.ReplaceB
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)

	return hex.EncodeToString(hashed)
}

func (r *Result) ComputeResult() error {
	if r.Fizzbuzz.ModA == 0 || r.Fizzbuzz.ModB == 0 || r.Fizzbuzz.Limit == 0 {
		return errors.New("Cannot modulo by 0")
	}

	chain := []string{}
	var i uint16
	for i = 1; i < r.Fizzbuzz.Limit; i++ {
		if i%(r.Fizzbuzz.ModA+r.Fizzbuzz.ModB) == 0 {
			chain = append(chain, r.Fizzbuzz.ReplaceA+r.Fizzbuzz.ReplaceB)
		} else if i%r.Fizzbuzz.ModA == 0 {
			chain = append(chain, r.Fizzbuzz.ReplaceA)
		} else if i%r.Fizzbuzz.ModB == 0 {
			chain = append(chain, r.Fizzbuzz.ReplaceB)
		} else {
			chain = append(chain, strconv.Itoa(int(i)))
		}
	}

	r.Result = chain
	return nil
}
