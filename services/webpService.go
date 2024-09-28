package services

/*
#cgo CFLAGS: -I${SRCDIR}/webp/include
#cgo LDFLAGS: -L${SRCDIR}/webp/lib -lwebp
#include "../webp/include/webp/encode.h"
#include <stdlib.h>
*/
import "C"

import (
	"image"
	"img2webp/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

type WebpService interface {
	GetFileSizeString(path string) int64
	ConvertToWebp(path string) int64
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

func (s *WebpServiceImpl) ConvertToWebp(path string) int64 {

	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening file:", err)
		return -1
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("Error decoding image:", err)
		return -1
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
		log.Println("Error encoding to WebP")
		return -1
	}

	goOutput := C.GoBytes(unsafe.Pointer(output), C.int(outputSize))

	defer C.free(unsafe.Pointer(output))

	outputFileName := utils.ExtractFileName(path)
	ext := filepath.Ext(outputFileName)
	outputPath := utils.OUTPUT_PATH + strings.TrimSuffix(outputFileName, ext) + ".webp"

	err = os.WriteFile(outputPath, goOutput, 0644)
	if err != nil {
		log.Println("Error writing WebP file:", err)
		return -1
	}

	log.Println("Image converted successfully")

	return s.GetFileSizeString(outputPath)
}
