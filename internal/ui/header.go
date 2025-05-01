package ui

import (
	"github.com/rivo/tview"
)

func header() *tview.TextView {
	text := `      █████████       ████████   █████████████ ██████████████████  
   ██████   ████   ██████  █████ ███     ██████ ██████████         
  ██████   █████ ██████    █████       ██████ ███  ██              
 ██████   ██████ █████    ██████     ██████      ███               
 ████    ██████ █████   ███████   ███████                          
 ████  ██████   ████   ██████    ██████                            
  █████████       █████████    ██████`
	image := tview.NewTextView().
		SetText(text)

	image.SetBackgroundColor(ColorBackground)

	return image
}
