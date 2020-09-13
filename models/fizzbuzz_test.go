package models

import (
	"errors"
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

func TestComputeResult(t *testing.T) {
	tests := map[string]struct {
		result   *Result
		expected []string
		error    error
	}{
		"case1": {
			result: &Result{
				Fizzbuzz: &Fizzbuzz{},
			},
			error: errors.New(""),
		},

		"case2": {
			result: &Result{
				Fizzbuzz: &Fizzbuzz{
					ModA:     0,
					ModB:     5,
					Limit:    100,
					ReplaceA: "fizz",
					ReplaceB: "buzz",
				},
			},
			error: errors.New(""),
		},

		"case3": {
			result: &Result{
				Fizzbuzz: &Fizzbuzz{
					ModA:     3,
					ModB:     5,
					Limit:    100,
					ReplaceA: "fizz",
					ReplaceB: "buzz",
				},
			},
			expected: []string{"1", "2", "fizz", "4", "buzz", "fizz", "7", "fizzbuzz", "fizz", "buzz", "11", "fizz", "13", "14", "fizz", "fizzbuzz", "17", "fizz", "19", "buzz", "fizz", "22", "23", "fizzbuzz", "buzz", "26", "fizz", "28", "29", "fizz", "31", "fizzbuzz", "fizz", "34", "buzz", "fizz", "37", "38", "fizz", "fizzbuzz", "41", "fizz", "43", "44", "fizz", "46", "47", "fizzbuzz", "49", "buzz", "fizz", "52", "53", "fizz", "buzz", "fizzbuzz", "fizz", "58", "59", "fizz", "61", "62", "fizz", "fizzbuzz", "buzz", "fizz", "67", "68", "fizz", "buzz", "71", "fizzbuzz", "73", "74", "fizz", "76", "77", "fizz", "79", "fizzbuzz", "fizz", "82", "83", "fizz", "buzz", "86", "fizz", "fizzbuzz", "89", "fizz", "91", "92", "fizz", "94", "buzz", "fizzbuzz", "97", "98", "fizz"},
			error:    nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.result.ComputeResult()
			assert.Equal(t, tc.expected, tc.result.Result)
			assert.IsType(t, tc.error, err)
		})
	}
}
