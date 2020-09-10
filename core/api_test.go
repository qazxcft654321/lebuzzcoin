package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIVersionFromFile(t *testing.T) {
	tests := map[string]struct {
		file         string
		stringLenght int
		expected     interface{}
	}{
		"case1": {
			file:         "BADFILE",
			stringLenght: 0,
			expected:     &os.PathError{},
		},

		"case2": {
			file:         "../VERSION",
			stringLenght: 5,
			expected:     nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			version, err := GetAPIVersionFromFile(tc.file)
			assert.IsType(t, tc.expected, err)
			assert.Equal(t, len(version), tc.stringLenght)
		})
	}
}
