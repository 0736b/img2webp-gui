package models

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ImageItem struct {
	Path         string
	FileSize     string
	IsConverting bool
}

func NewImageItemWidget(item *ImageItem, update func()) *fyne.Container {

	pathLabel := widget.NewLabel(item.Path)
	sizeLabel := widget.NewLabel(item.FileSize)
	sizeLabel.Hide()

	loading := widget.NewActivity()
	loading.Start()

	fileInfoContainer := container.NewStack(
		loading, sizeLabel,
	)

	if !item.IsConverting {
		loading.Stop()
		loading.Hide()
		sizeLabel.Show()
	}

	return container.NewGridWithColumns(2, pathLabel, container.NewBorder(nil, nil, nil, fileInfoContainer))
}
