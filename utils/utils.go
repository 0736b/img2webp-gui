package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	OutputDirPath string = "./output/"
)

const (
	_        = iota
	KB int64 = 1 << (10 * iota)
	MB
	GB
)

func FormatFileSize(size int64) string {

	switch {
	case size == -1:
		return ""
	case size < int64(KB):
		return fmt.Sprintf("%d B", size)
	case size < int64(MB):
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	case size < int64(GB):
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	}
}

func CreateOutputDir() error {

	err := os.MkdirAll(OutputDirPath, os.ModeDir)
	if err != nil {
		log.Println("CreateOutputDir", err.Error())
		return err
	}

	return nil
}

func ExtractFileName(path string) string {

	return filepath.Base(path)
}
