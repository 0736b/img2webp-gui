package models

import (
	"fmt"
	"img2webp/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ImageItem struct {
	Path              string
	FileName          string
	OriginalFileSize  int64
	ConvertedFileSize int64
	IsConverting      bool
	Thumbnail         *canvas.Image
}

func NewImageItemWidget(item *ImageItem, update func()) *fyne.Container {

	fileNameLabel := widget.NewLabel(item.FileName)
	fileNameLabel.TextStyle = fyne.TextStyle{Bold: true}

	originalSizeLabel := widget.NewLabel(utils.FormatFileSize(item.OriginalFileSize))

	leftSection := container.NewGridWithRows(2, fileNameLabel, originalSizeLabel)

	percentageSizeLabel := widget.NewLabel("")
	percentageSizeLabel.TextStyle = fyne.TextStyle{Bold: true}

	convertedSizeLabel := widget.NewLabel(utils.FormatFileSize(item.ConvertedFileSize))

	loading := widget.NewActivity()
	loading.Start()

	rightSection := container.NewGridWithRows(2, percentageSizeLabel, convertedSizeLabel)
	rightSection.Hide()

	fileInfoContainer := container.NewStack(
		loading, rightSection,
	)

	if !item.IsConverting {
		loading.Stop()
		loading.Hide()
		if item.ConvertedFileSize != -1 {
			percentageSizeLabel.SetText(calcPercentage(item.OriginalFileSize, item.ConvertedFileSize))
		}
		rightSection.Show()
	}

	return container.NewGridWithRows(2, container.NewGridWithColumns(2, leftSection, container.NewBorder(nil, nil, nil, fileInfoContainer)), container.New(layout.NewCustomPaddedLayout(0, 0, 0, 0)))
}

func calcPercentage(originalSize, convertedSize int64) string {
	if convertedSize >= originalSize {
		return fmt.Sprintf("+ %.0f %%", (float32(convertedSize-originalSize)/float32(originalSize))*100)
	} else {
		return fmt.Sprintf("- %.0f %%", (float32(originalSize-convertedSize)/float32(originalSize))*100)
	}
}
