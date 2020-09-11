package models

type Fizzbuzz struct {
	ModA     int    `json:"mod_a" validate:"required,numeric"`
	ModB     int    `json:"mod_b" validate:"required,numeric"`
	Limit    int    `json:"limit" validate:"required,numeric"`
	ReplaceA string `json:"replace_a" validate:"required,alphanum"`
	ReplaceB string `json:"replace_b" validate:"required,alphanum"`
}
