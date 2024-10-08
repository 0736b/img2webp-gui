package services

/*
#include "../webp/include/webp/encode.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"img2webp/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

type WebpService interface {
	GetFileSize(path string) (int64, error)
	ConvertToWebp(path string, outputPath string) (string, error)
}

type WebpServiceImpl struct{}

func NewWebpService() *WebpServiceImpl {
	return &WebpServiceImpl{}
}

func (s *WebpServiceImpl) GetFileSize(path string) (int64, error) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println("GetFileSize failed to read file stat", err.Error())
		return -1, err
	}

	return fileInfo.Size(), nil
}

func (s *WebpServiceImpl) ConvertToWebp(path string, outputPath string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		log.Println("ConvertToWebp failed to opening file:", err.Error())
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("ConvertToWebp failed decoding image:", err.Error())
		return "", err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	rgba := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	var output *C.uint8_t
	outputSize := C.size_t(0)
	outputSize = C.WebPEncodeRGBA(
		(*C.uint8_t)(unsafe.Pointer(&rgba.Pix[0])),
		C.int(width),
		C.int(height),
		C.int(rgba.Stride),
		C.float(75),
		&output,
	)

	if outputSize == 0 {
		log.Println("ConvertToWebp failed to encoding to WebP")
		return "", fmt.Errorf("error encoding to webp")
	}

	webpOutput := C.GoBytes(unsafe.Pointer(output), C.int(outputSize))

	defer C.free(unsafe.Pointer(output))

	fileName := utils.ExtractFileName(path)
	ext := filepath.Ext(fileName)
	writePath := outputPath + strings.TrimSuffix(fileName, ext) + ".webp"

	err = os.WriteFile(writePath, webpOutput, 0644)
	if err != nil {
		log.Println("ConvertToWebp failed to writing WebP file:", err.Error())
		return "", fmt.Errorf("error failed to wrting webp file")
	}

	// log.Println("ConvertToWebp successfully", writePath)

	return writePath, nil
}
