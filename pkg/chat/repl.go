package chat

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type REPL struct {
	inputField *tview.InputField
	history    []string
	currentIdx int
	chaTUI     *ChaTUI
}

func NewREPL(chaTUI *ChaTUI) *REPL {
	inputField := tview.NewInputField()
	repl := &REPL{
		inputField: inputField,
		history:    []string{},
		currentIdx: -1,
		chaTUI:     chaTUI,
	}

	repl.setupInputField()

	return repl
}

func (r *REPL) setupInputField() {
	r.inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			r.handleUserInput()
		case tcell.KeyUp:
			r.previousHistory()
		case tcell.KeyDown:
			r.nextHistory()
		}
		return event
	})
}

func (r *REPL) handleUserInput() {
	input := r.inputField.GetText()
	if strings.TrimSpace(input) == "" {
		return
	}

	r.addToHistory(input)

	response, err := r.chaTUI.requestCompletion(input)
	if err != nil {
		// If there is an error, clear the input field and set an error message
		r.inputField.SetText("")
		r.chaTUI.conversationPane.SetText("Error: Failed to get a response from the ChatGPT API.")
		return
	}

	// Update the conversation pane with the user input and ChatGPT response
	r.chaTUI.conversationPane.SetText(r.chaTUI.conversationPane.GetText(false) + "\n\nUser: " + input + "\nChatGPT: " + response)
	r.inputField.SetText("")
}

func (r *REPL) addToHistory(command string) {
	r.history = append(r.history, command)
	r.currentIdx = len(r.history)
}

func (r *REPL) previousHistory() {
	if r.currentIdx > 0 {
		r.currentIdx--
		r.inputField.SetText(r.history[r.currentIdx])
	}
}

func (r *REPL) nextHistory() {
	if r.currentIdx < len(r.history)-1 {
		r.currentIdx++
		r.inputField.SetText(r.history[r.currentIdx])
	} else {
		r.currentIdx = len(r.history)
		r.inputField.SetText("")
	}
}
