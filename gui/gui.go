package gui

import (
	"fmt"
	"image/color"
	"img2webp/gui/models"
	"img2webp/services"
	"img2webp/utils"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AppState struct {
	win fyne.Window

	service services.WebpService

	fileList   []*models.ImageItem
	listWidget *widget.List

	statusLabel    *widget.Label
	convertedCount int
	countChan      chan struct{}

	mutex sync.Mutex
}

func Run(service services.WebpService) {

	a := app.New()
	w := a.NewWindow("Img2Webp Converter")
	w.Resize(fyne.NewSize(648, 324))
	w.SetFixedSize(true)

	ui := &AppState{
		win:     w,
		service: service,

		fileList:       []*models.ImageItem{},
		convertedCount: 0,
		countChan:      make(chan struct{}, 1),
	}

	_list := widget.NewList(

		func() int {
			return len(ui.fileList)
		},

		func() fyne.CanvasObject {
			return models.NewImageItemWidget(&models.ImageItem{}, ui.forceRefreshList)
		},

		func(i widget.ListItemID, o fyne.CanvasObject) {
			ui.mutex.Lock()
			item := ui.fileList[i]
			widget := o.(*fyne.Container)
			widget.Objects = models.NewImageItemWidget(item, ui.forceRefreshList).Objects
			ui.mutex.Unlock()
			ui.listWidget.SetItemHeight(i, 50)
		})

	ui.listWidget = _list

	_dropLabel := widget.NewLabel("Drag and drop your image files")

	_statusLabel := widget.NewLabel("Waiting for files...")
	ui.statusLabel = _statusLabel

	_clearBtn := widget.NewButtonWithIcon("Clear", theme.DeleteIcon(), ui.onClearList)

	_bg := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 60})

	_scrollContainer := container.NewVScroll(ui.listWidget)

	_content := container.NewBorder(
		container.NewCenter(_dropLabel),
		container.NewBorder(nil, nil, ui.statusLabel, _clearBtn), nil, nil,
		container.New(layout.CustomPaddedLayout{TopPadding: 8, BottomPadding: 8, LeftPadding: 0, RightPadding: 0}, container.NewStack(_bg, _scrollContainer)),
	)

	ui.win.SetOnDropped(ui.onDropFiles)

	ui.win.SetContent(container.New(layout.CustomPaddedLayout{TopPadding: 0, BottomPadding: 8, LeftPadding: 12, RightPadding: 12}, _content))

	go ui.convertedCounter()

	ui.win.ShowAndRun()

}

func (ui *AppState) forceRefreshList() {

	ui.listWidget.Refresh()
}

func (ui *AppState) onClearList() {

	ui.mutex.Lock()

	oldLen := len(ui.fileList)

	n := 0
	for _, x := range ui.fileList {
		if x.IsConverting {
			ui.fileList[n] = x
			n++
		}
	}
	ui.fileList = ui.fileList[:n]
	ui.convertedCount = 0

	ui.mutex.Unlock()

	if len(ui.fileList) == 0 {
		ui.statusLabel.SetText("Waiting for files...")
	}

	if oldLen != len(ui.fileList) {
		ui.listWidget.Refresh()
	}
}

func (ui *AppState) onDropFiles(pos fyne.Position, uris []fyne.URI) {

	for _, uri := range uris {

		item := &models.ImageItem{
			Path:              uri.Path(),
			FileName:          utils.ExtractFileName(uri.Path()),
			OriginalFileSize:  ui.service.GetFileSizeString(uri.Path()),
			ConvertedFileSize: -1,
			IsConverting:      true,
		}

		ui.mutex.Lock()
		ui.fileList = append(ui.fileList, item)
		ui.mutex.Unlock()

		ui.listWidget.Refresh()

		go ui.convertFile(item, ui.forceRefreshList)
	}

	ui.mutex.Lock()
	ui.statusLabel.SetText("Converting...")
	ui.mutex.Unlock()
}

func (ui *AppState) convertFile(item *models.ImageItem, update func()) {

	time.Sleep(2 * time.Second)

	ui.mutex.Lock()
	item.ConvertedFileSize = 10000
	item.IsConverting = false
	ui.mutex.Unlock()

	ui.countChan <- struct{}{}

	update()
}

func (ui *AppState) convertedCounter() {

	for range ui.countChan {
		ui.convertedCount++
		ui.statusLabel.SetText(fmt.Sprintf("Converted %d/%d files", ui.convertedCount, len(ui.fileList)))
	}
}
