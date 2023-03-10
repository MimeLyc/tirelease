package fileserver

import (
	"fmt"
	"testing"
	"tirelease/commons/ifile"

	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	dir := "tmp"
	name := "test.txt"
	qualifiedName := fmt.Sprintf("%s/%s", dir, name)
	err := ifile.CreateFileRecursively(dir, name)
	assert.Nil(t, err)
	targetUrl := "tirelease/tmp/test_111.txt"
	downloadUrl, err := UploadFile(qualifiedName, targetUrl)
	assert.Nil(t, err)
	assert.NotEqual(t, "", downloadUrl)
	err = ifile.RmAllFile(dir)
	assert.Nil(t, err)
}
