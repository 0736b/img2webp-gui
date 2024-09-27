package main

/*
#cgo CFLAGS: -I${SRCDIR}/webp/include
#cgo LDFLAGS: -L${SRCDIR}/webp/lib -lwebp
#include "webp/encode.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"unsafe"
)

func main() {
	// Open an image file
	file, err := os.Open("input.jpg")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Get image dimensions
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Convert image to RGBA
	rgba := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	// Encode to WebP
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
		fmt.Println("Error encoding to WebP")
		return
	}

	// Convert C buffer to Go slice
	goOutput := C.GoBytes(unsafe.Pointer(output), C.int(outputSize))

	// Free the C buffer
	defer C.free(unsafe.Pointer(output))

	// Save the WebP image
	err = os.WriteFile("output.webp", goOutput, 0644)
	if err != nil {
		fmt.Println("Error writing WebP file:", err)
		return
	}

	fmt.Println("Image converted successfully")
}
