package gui

import (
	"fmt"
	"image/color"
	"img2webp/gui/models"
	"img2webp/services"
	"img2webp/utils"
	"log"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AppState struct {
	win            fyne.Window
	service        services.WebpService
	fileList       []*models.ImageItem
	listWidget     *widget.List
	statusLabel    *widget.Label
	convertedCount int32
	mutex          sync.RWMutex
}

func NewAppState(w fyne.Window, service services.WebpService) *AppState {

	return &AppState{
		win:      w,
		service:  service,
		fileList: []*models.ImageItem{},
	}
}

func (ui *AppState) SetupUI() {

	ui.listWidget = ui.createListWidget()
	dropLabel := widget.NewLabel("Drag and drop your image files")
	ui.statusLabel = widget.NewLabel("Waiting for files...")
	clearBtn := widget.NewButtonWithIcon("Clear", theme.DeleteIcon(), ui.onClearList)

	bg := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 60})
	scrollContainer := container.NewVScroll(ui.listWidget)

	content := container.NewBorder(
		container.NewCenter(dropLabel),
		container.NewBorder(nil, nil, ui.statusLabel, clearBtn), nil, nil,
		container.New(layout.CustomPaddedLayout{TopPadding: 8, BottomPadding: 8, LeftPadding: 0, RightPadding: 0}, container.NewStack(bg, scrollContainer)),
	)

	ui.win.SetOnDropped(ui.onDropFiles)
	ui.win.SetContent(container.New(layout.CustomPaddedLayout{TopPadding: 0, BottomPadding: 8, LeftPadding: 12, RightPadding: 12}, content))

}

func (ui *AppState) createListWidget() *widget.List {

	return widget.NewList(
		func() int {
			ui.mutex.RLock()
			defer ui.mutex.RUnlock()
			return len(ui.fileList)
		},
		func() fyne.CanvasObject {
			return models.NewImageItemWidget(&models.ImageItem{}, ui.forceRefreshList)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			ui.mutex.RLock()
			item := ui.fileList[i]
			ui.mutex.RUnlock()
			widget := o.(*fyne.Container)
			widget.Objects = models.NewImageItemWidget(item, ui.forceRefreshList).Objects
			ui.listWidget.SetItemHeight(i, 50)
		})
}

func (ui *AppState) forceRefreshList() {

	ui.listWidget.Refresh()
}

func (ui *AppState) onClearList() {

	ui.mutex.Lock()
	oldFileListLen := len(ui.fileList)
	newFileList := make([]*models.ImageItem, 0, len(ui.fileList))
	for _, item := range ui.fileList {
		if item.IsConverting {
			newFileList = append(newFileList, item)
		}
	}
	ui.fileList = newFileList
	atomic.StoreInt32(&ui.convertedCount, 0)
	ui.mutex.Unlock()

	if len(ui.fileList) == 0 {
		ui.statusLabel.SetText("Waiting for files...")
	}

	if oldFileListLen != len(ui.fileList) {
		ui.listWidget.Refresh()
	}
}

func (ui *AppState) onDropFiles(pos fyne.Position, uris []fyne.URI) {

	for _, uri := range uris {
		item := &models.ImageItem{
			Path:              uri.Path(),
			FileName:          utils.ExtractFileName(uri.Path()),
			OriginalFileSize:  ui.service.GetFileSize(uri.Path()),
			ConvertedFileSize: -1,
			IsConverting:      true,
		}
		ui.mutex.Lock()
		ui.fileList = append(ui.fileList, item)
		ui.mutex.Unlock()
		ui.listWidget.Refresh()
		go ui.convertFile(item)
	}

	ui.statusLabel.SetText("Converting...")
}

func (ui *AppState) convertFile(item *models.ImageItem) {

	convertedPath, err := ui.service.ConvertToWebp(item.Path)
	if err != nil {
		log.Println("convertFile failed", err.Error())
		return
	}

	if convertedPath != "" {
		ui.mutex.Lock()
		item.ConvertedFileSize = ui.service.GetFileSize(convertedPath)
		item.IsConverting = false
		ui.mutex.Unlock()
		ui.addConvertedCount()
		ui.forceRefreshList()
	}
}

func (ui *AppState) addConvertedCount() {

	count := atomic.AddInt32(&ui.convertedCount, 1)
	ui.mutex.RLock()
	totalFiles := len(ui.fileList)
	ui.mutex.RUnlock()
	ui.statusLabel.SetText(fmt.Sprintf("Converted %d/%d files", count, totalFiles))
}
