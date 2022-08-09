package utils

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var TimeStampFormat = "2006-01-02T15:04:05.000Z07:00"

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
