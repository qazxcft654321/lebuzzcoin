package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFizzbuzzHashData(t *testing.T) {
	tests := map[string]struct {
		fizzbuzz *Fizzbuzz
		expected string
	}{
		"case1": {
			fizzbuzz: &Fizzbuzz{},
			expected: "709e80c88487a2411e1ee4dfb9f22a861492d20c4765150c0c794abd70f8147c",
		},

		"case2": {
			fizzbuzz: &Fizzbuzz{
				ModA:     2,
				ModB:     10,
				Limit:    10000,
				ReplaceA: "testA",
				ReplaceB: "testB",
			},
			expected: "cc32648928a141cfb49173b8649f3a760afb525c7863604b1692249b3e2d7614",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hash := tc.fizzbuzz.HashData()
			assert.Equal(t, tc.expected, hash)
		})
	}

}
