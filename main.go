package main

/*
#cgo CFLAGS: -I${SRCDIR}/webp/include
#cgo LDFLAGS: -L${SRCDIR}/webp/lib -lwebp
*/
import "C"

import (
	"img2webp/gui"
	"img2webp/services"
	"img2webp/utils"
)

func main() {

	// ensure the output dir is exists
	utils.CreateOutputDir()

	service := services.NewWebpService()

	gui.Run(service)

}
