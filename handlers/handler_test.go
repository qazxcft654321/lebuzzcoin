package handlers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	h := New(os.Stdout)

	assert.IsType(t, h, &Handler{})
	assert.NotEmpty(t, h)
}
