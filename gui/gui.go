package gui

import (
	"fmt"
	"image/color"
	"img2webp/gui/models"
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

	fileList   []*models.ImageItem
	listWidget *widget.List

	statusLabel *widget.Label

	mutex sync.Mutex
}

func Run() {

	a := app.New()
	w := a.NewWindow("Img2Webp")
	w.Resize(fyne.NewSize(720, 360))
	w.SetFixedSize(true)

	ui := &AppState{
		win: w,

		fileList: []*models.ImageItem{},
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
		})

	ui.listWidget = _list

	_dropLabel := widget.NewLabel("Drag and drop your image files")

	_statusLabel := widget.NewLabel("Waiting for files...")
	ui.statusLabel = _statusLabel

	_clearBtn := widget.NewButtonWithIcon("Clear", theme.DeleteIcon(), ui.onClearList)

	_bg := canvas.NewRectangle(color.RGBA{R: 17, G: 17, B: 18, A: 255})

	_scrollContainer := container.NewVScroll(ui.listWidget)

	_content := container.NewBorder(
		container.NewCenter(_dropLabel),
		container.NewBorder(nil, nil, ui.statusLabel, _clearBtn), nil, nil,
		container.New(layout.CustomPaddedLayout{TopPadding: 8, BottomPadding: 8, LeftPadding: 0, RightPadding: 0}, container.NewStack(_bg, _scrollContainer)),
	)

	ui.win.SetOnDropped(ui.onDropFiles)

	ui.win.SetContent(container.New(layout.CustomPaddedLayout{TopPadding: 0, BottomPadding: 8, LeftPadding: 12, RightPadding: 12}, _content))

	ui.win.ShowAndRun()

}

func (ui *AppState) forceRefreshList() {

	ui.listWidget.Refresh()
}

func (ui *AppState) onClearList() {

	ui.mutex.Lock()
	ui.fileList = make([]*models.ImageItem, 0)
	ui.mutex.Unlock()

	ui.listWidget.Refresh()
}

func (ui *AppState) onDropFiles(pos fyne.Position, uris []fyne.URI) {

	for _, uri := range uris {

		item := &models.ImageItem{
			Path:         uri.Path(),
			FileSize:     "",
			IsConverting: true,
		}

		ui.mutex.Lock()
		ui.fileList = append(ui.fileList, item)
		ui.mutex.Unlock()

		ui.listWidget.Refresh()

		ui.statusLabel.SetText(fmt.Sprintf("Converting %d files", len(ui.fileList)))

		go ui.convertFile(item, ui.forceRefreshList)
	}
}

func (ui *AppState) convertFile(item *models.ImageItem, update func()) {

	time.Sleep(2 * time.Second)
	item.IsConverting = false
	update()
}
