package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersects(t *testing.T) {
	s := []string{"foo", "bar"}
	m := []string{"foo"}

	result := Intersects(s, m)
	assert.Equal(t, result, []string{"foo"})

	s = []string{"foo", "bar"}
	m = []string{}

	result = Intersects(s, m)
	assert.Equal(t, result, []string{})

	s = []string{"foo", "bar"}
	m = []string{"outer"}

	result = Intersects(s, m)
	assert.Equal(t, result, []string{})

}
