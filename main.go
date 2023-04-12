package main

import (
	"os"

	"github.com/xyproto/chatui/pkg/chat"
	"github.com/xyproto/env/v2"
)

var APIKey = env.Str("OPENAI_KEY")

func main() {
	app := chat.NewChaTUI(APIKey)
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
