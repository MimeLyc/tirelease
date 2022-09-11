package ifile

import (
	"fmt"
	"os"
)

func MustOpen(name string) *os.File {
	r, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	return r
}

func CreateFileRecursively(dir, name string) error {
	qualifiedName := fmt.Sprintf("%s/%s", dir, name)
	if _, err := os.Stat(qualifiedName); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
		file, err := os.Create(qualifiedName)

		if err != nil {
			return err
		}

		file.Close()
		return err
	}

	return nil
}

func RmAllFile(qualifiedName string) error {
	return os.RemoveAll(qualifiedName)
}
