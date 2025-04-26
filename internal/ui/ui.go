package ui

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/renantatsuo/james-bond/internal/agent"
	"github.com/rivo/tview"
)

type UI struct {
	app      *tview.Application
	pages    *tview.Pages
	errorBox *tview.Modal
	agent    *agent.Agent
}

type SendMessage func(string)

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
		SetBorder(true).
		SetBorderColor(ColorBackground)

	input := textInput(ui.handleSubmit(ctx, messages))

	ui.errorBox = tview.NewModal().
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.pages.HidePage("error")
		})
	ui.errorBox.Box.SetTitle("An error occurred")

	grid := tview.NewGrid().
		SetRows(6, 0, 2).
		SetColumns(0).
		SetBorders(true).
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

		response, err := ui.agent.SendUserMessage(ctx, message, "gpt-4.1-nano")
		if err != nil {
			ui.ShowError(err)
			return
		}

		if response == "" {
			ui.ShowError(fmt.Errorf("something got wrong, response is empty"))
			return
		}

		unquoted, err := strconv.Unquote(response)
		if err != nil {
			fmt.Fprintf(messagesWriter, "[::d]James Bond:[::-] %s\n", response)
			return
		}

		fmt.Fprintf(messagesWriter, "[::d]James Bond:[::-] %s\n", unquoted)
	}
}

func (ui *UI) ShowError(err error) {
	ui.errorBox.SetText(err.Error())
	ui.pages.ShowPage("error")
}

func (ui *UI) Stop() {
	if ui.app != nil {
		ui.app.Stop()
	}
}
