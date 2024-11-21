package fileutil

import (
	"errors"
	"os"
	"path/filepath"
)

func WriteFile(name string, data []byte) error  {
	dir:=filepath.Dir(name)
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir,os.ModePerm)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(name,data,os.ModePerm)
}
