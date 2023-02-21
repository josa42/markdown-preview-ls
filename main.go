package main

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"

	"github.com/josa42/markdown-preview-ls/control"
	"github.com/josa42/markdown-preview-ls/preview"
	"github.com/josa42/markdown-preview-ls/previewserver"
	"github.com/josa42/markdown-preview-ls/server"
)

func main() {
	cmd := "server"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "preview":
		text := ""
		if len(os.Args) > 2 {
			text = os.Args[2]
		}
		runPreview(text)
	default:
		runServer()
	}
}

func runServer() {
	ch := control.Channels{
		Open:   make(chan string),
		Close:  make(chan bool),
		Update: make(chan string),
	}

	previewIsOpen := false

	go func() {
		for {
			text := <-ch.Open
			if !previewIsOpen {
				previewIsOpen = true
				cmd := exec.Command(os.Args[0], "preview", text)
				go func() {
					cmd.Run()
					previewIsOpen = false
				}()
			} else {
				http.Post("http://localhost:3333/update", "text/plain", bytes.NewBufferString(text))

			}
		}
	}()
	go func() {
		for {
			<-ch.Close
			if previewIsOpen {
				http.Post("http://localhost:3333/close", "text/plain", nil)
			}
		}
	}()

	go func() {
		for {
			text := <-ch.Update
			http.Post("http://localhost:3333/update", "text/plain", bytes.NewBufferString(text))
		}
	}()

	server.Run(ch)
}

func runPreview(text string) {
	ch := control.PreviewChannels{
		Close:  make(chan bool),
		Update: make(chan string),
	}

	go previewserver.Run(ch)

	preview.Run(ch, text)
}
