package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	OUTPUT_PATH string = "./output/"
)

func FormatFileSize(size int64) string {

	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
}

func CreateOutputDir() error {

	err := os.MkdirAll(OUTPUT_PATH, os.ModeDir)
	if err != nil {
		log.Println("CreateOutputDir", err.Error())
		return err
	}

	return nil
}

func ExtractFileName(path string) string {

	return filepath.Base(path)
}
