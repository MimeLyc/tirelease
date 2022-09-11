package ifile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFile(t *testing.T) {
	dir := "tmp"
	name := "test.txt"
	qualifiedName := fmt.Sprintf("%s/%s", dir, name)
	err := CreateFileRecursively(dir, name)
	assert.Nil(t, err)
	file := MustOpen(qualifiedName)
	assert.NotNil(t, file)

	RmAllFile(dir)
}
