package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type inputSubmitHandler func(message string)

func textInput(submitHandler inputSubmitHandler) *tview.TextArea {
	input := tview.NewTextArea().
		SetLabel("Message: ").
		SetLabelStyle(tcell.StyleDefault.Background(ColorBackground))
	input.SetTextStyle(tcell.StyleDefault.Background(ColorBackground))
	input.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		if key.Key() == tcell.KeyEnter {
			message := input.GetText()
			go submitHandler(message)
			input.SetText("", false)
			return nil
		}

		return key
	})
	input.SetBackgroundColor(ColorBackground).
		SetBorder(true)
	return input
}
