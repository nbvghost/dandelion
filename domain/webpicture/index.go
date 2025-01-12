package webpicture

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var cwebpPath string = ""
var gif2webpPath string = ""

// Supported input formats:
// WebP, JPEG, PNG, PNM (PGM, PPM, PAM), TIFF
func init() {
	dir, dirErr := os.UserCacheDir()
	if dirErr == nil {
		dir = filepath.Join(dir, "dandelion", "webp")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Printf("can't create user cache dir: %v", err)
	}
	cwebpPath = filepath.Join(dir, "cwebp")
	err := os.WriteFile(cwebpPath, cwebpBytes, os.ModePerm)
	if err != nil {
		panic(err)
	}

	gif2webpPath = filepath.Join(dir, "gif2webp")
	err = os.WriteFile(gif2webpPath, gif2webpBytes, os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("use cwebp path:%s", cwebpPath))
	log.Println(fmt.Sprintf("use gif2webp path:%s", gif2webpPath))
}
func EncodeGIF(fromFileName string, saveFileName string) error {
	cmd := exec.Command(gif2webpPath, fromFileName, "-q", "90", "-o", saveFileName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
func Encode(fromFileName string, saveFileName string) error {
	cmd := exec.Command(cwebpPath, fromFileName, "-q", "90", "-o", saveFileName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

/*
func Encode(imgBytes []byte, saveFileName string) error {
	cmd := exec.Command(
		cwebpPath, // get it from: https://storage.googleapis.com/downloads.webmproject.org/releases/webp/index.html
		"-q", "90",
		// https://developers.google.com/speed/webp/docs/dwebp
		"-o", saveFileName, // stdout
		"--", "-", // stdin
	)
	var output bytes.Buffer
	cmd.Stdout = &output

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}

	go func() {
		_, _ = stdin.Write(imgBytes)
		_ = stdin.Close()
	}()

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}

*/
