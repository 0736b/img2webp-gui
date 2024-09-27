package services

import (
	"log"
	"os"
)

type WebpService interface {
	GetFileSizeString(path string) int64
}

type WebpServiceImpl struct{}

func NewWebpService() *WebpServiceImpl {
	return &WebpServiceImpl{}
}

func (s *WebpServiceImpl) GetFileSizeString(path string) int64 {

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println("GetFileSizeString", err.Error())
		return -1
	}

	fileSize := fileInfo.Size()

	return fileSize
}
