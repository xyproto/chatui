package chat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	historyPaneWidth = 30
	OpenAIAPIURL     = "https://api.openai.com/v1/engines/davinci-codex/completions"
)

type ChaTUI struct {
	app              *tview.Application
	conversationPane *tview.TextView
	conversationList *tview.List
	chatHistory      []*tview.TextView
	apiKey           string
}

func NewChaTUI(apiKey string) *ChaTUI {
	app := tview.NewApplication()

	chaTUI := &ChaTUI{
		app:              app,
		apiKey:           apiKey,
		conversationPane: tview.NewTextView(),
		conversationList: tview.NewList(),
		chatHistory:      []*tview.TextView{},
	}

	chaTUI.setupUI()
	chaTUI.setupKeybindings()

	return chaTUI
}

func (c *ChaTUI) Run() error {
	return c.app.Run()
}

func (c *ChaTUI) setupUI() {
	flex := tview.NewFlex().
		AddItem(c.conversationPane, 0, 1, true).
		AddItem(c.conversationList, historyPaneWidth, 0, false)

	c.app.SetRoot(flex, true)
}

func (c *ChaTUI) setupKeybindings() {
	c.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Implement keybindings and functionalities
		return event
	})
}

func (c *ChaTUI) requestCompletion(prompt string) (string, error) {
	client := &http.Client{}

	payload := map[string]interface{}{
		"prompt":   prompt,
		"max_tokens": 50,
		"n": 1,
		"stop": []string{"\n"},
	}

	payloadJSON, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", OpenAIAPIURL, strings.NewReader(string(payloadJSON)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	choices := result["choices"].([]interface{})
	choice := choices[0].(map[string]interface{})
	text := choice["text"].(string)

	return text, nil
}
