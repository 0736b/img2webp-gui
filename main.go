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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {

	// ensure the output dir is exists
	utils.CreateOutputDir()

	service := services.NewWebpService()

	a := app.New()
	w := a.NewWindow("Img2Webp Converter")
	w.Resize(fyne.NewSize(648, 324))
	w.SetFixedSize(true)

	appState := gui.NewAppState(w, service)
	appState.SetupUI()

	w.ShowAndRun()

}
