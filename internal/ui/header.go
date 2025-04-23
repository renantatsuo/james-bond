package ui

import (
	"bytes"
	"image/png"
	"os"

	"github.com/rivo/tview"
)

func header() *tview.Image {
	imgSrc, err := os.ReadFile("assets/007.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(bytes.NewReader(imgSrc))
	if err != nil {
		panic(err)
	}
	image := tview.NewImage().
		SetImage(img).
		SetColors(tview.TrueColor).
		SetAlign(tview.AlignCenter, tview.AlignLeft)

	image.SetBackgroundColor(ColorBackground)

	return image
}
