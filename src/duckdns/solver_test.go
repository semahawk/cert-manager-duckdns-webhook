package duckdns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolver_Name(t *testing.T) {
	solver := NewSolver(nil)
	assert.Equal(t, "duckdns", solver.Name())
}
