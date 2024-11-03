//go:build !windows
package config

import (
	"os"
	"syscall"
)

func writeLock(filename string,b []byte) error  {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	err = syscall.Flock(int(f.Fd()),syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return err
	}
	return nil
}