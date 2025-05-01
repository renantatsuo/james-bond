package ui

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/renantatsuo/james-bond/internal/agent"
	"github.com/rivo/tview"
)

type UI struct {
	app      *tview.Application
	pages    *tview.Pages
	errorBox *tview.Modal
	agent    *agent.Agent
}

func New(agent *agent.Agent) *UI {
	return &UI{
		agent: agent,
	}
}

func (ui *UI) Init(ctx context.Context) error {
	ui.app = tview.NewApplication()

	header := header()

	messages := tview.NewTextView().
		SetChangedFunc(func() { ui.app.Draw() }).
		SetScrollable(true).
		SetDynamicColors(true).
		SetTextColor(ColorForeground)
	messages.SetBackgroundColor(ColorBackground)
	messages.Box.
		SetTitle("Messages").
		SetBorder(true)

	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//handle scroll
		scrollRow, scrollCol := messages.GetScrollOffset()
		if event.Key() == tcell.KeyUp {
			messages.ScrollTo(max(scrollRow-1, 0), scrollCol)
		} else if event.Key() == tcell.KeyDown {
			messages.ScrollTo(scrollRow+1, scrollCol)
		} else if event.Key() == tcell.KeyHome {
			messages.ScrollToBeginning()
		} else if event.Key() == tcell.KeyEnd {
			messages.ScrollToEnd()
		}
		return event
	})

	handleMessage := ui.handleSubmit(ctx, messages)
	input := textInput(func(message string) {
		handleMessage(message)
		messages.ScrollToEnd()
	})

	ui.errorBox = tview.NewModal().
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.pages.HidePage("error")
		})
	ui.errorBox.Box.SetTitle("An error occurred")

	grid := tview.NewGrid().
		SetRows(7, 0, 4).
		SetColumns(0).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(messages, 1, 0, 1, 1, 0, 0, false).
		AddItem(input, 2, 0, 1, 1, 0, 0, true)
	grid.Box.SetBackgroundColor(ColorBackground)

	ui.pages = tview.NewPages().
		AddPage("main", grid, true, true).
		AddPage("error", ui.errorBox, true, false)

	ui.app.SetRoot(ui.pages, true)

	if err := ui.app.Run(); err != nil {
		return err
	}

	return nil
}

func (ui *UI) handleSubmit(ctx context.Context, messagesWriter io.Writer) inputSubmitHandler {
	return func(message string) {
		fmt.Fprintf(messagesWriter, "[::d]You:[::-] %s \n", message)

		msg := agent.Message{Content: message, Type: agent.MessageTypeUser}

		response, err := ui.agent.SendUserMessage(ctx, msg, "gpt-4.1-mini")
		if err != nil {
			ui.ShowError(err)
			return
		}

		if response == "" {
			ui.ShowError(fmt.Errorf("something got wrong, response is empty"))
			return
		}

		unquoted, err := strconv.Unquote(response)
		if err == nil {
			response = unquoted
		}
		response = tview.TranslateANSI(escapeANSI(response))

		fmt.Fprintf(messagesWriter, "[::d]James Bond:[::-] %s\n", response)
	}
}

func (ui *UI) ShowError(err error) {
	ui.errorBox.SetText(err.Error())
	ui.pages.ShowPage("error")
	ui.app.Draw()
}

func (ui *UI) Stop() {
	if ui.app != nil {
		ui.app.Stop()
	}
}

func escapeANSI(s string) string {
	s = strings.ReplaceAll(s, `\033`, "\033")
	s = strings.ReplaceAll(s, `\x1b`, "\x1b")

	return s
}
