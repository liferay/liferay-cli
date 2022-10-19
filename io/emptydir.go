package io

import (
	"io"
	"os"
)

func Exists(name string) bool {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return true
	}

	return false
}

func IsDirEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}
