package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	type testType string
	test1 := testType("1")
	test2 := testType("2")

	assert.Equal(t, Join([]testType{test1, test2}, ",", "\""), "\"1\",\"2\"")
}
