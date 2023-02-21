package main

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"

	"github.com/josa42/md-ls/control"
	"github.com/josa42/md-ls/preview"
	"github.com/josa42/md-ls/previewserver"
	"github.com/josa42/md-ls/server"
)

func main() {
	cmd := "server"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "preview":
		runPreview()
	default:
		runServer()
	}
}

func runServer() {
	ch := control.Channels{
		Open:   make(chan bool),
		Close:  make(chan bool),
		Update: make(chan string),
	}

	previewIsOpen := false

	go func() {
		for {
			<-ch.Open
			if !previewIsOpen {
				previewIsOpen = true
				cmd := exec.Command(os.Args[0], "preview")
				go func() {
					cmd.Run()
					previewIsOpen = false
				}()
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

func runPreview() {
	ch := control.PreviewChannels{
		Close:  make(chan bool),
		Update: make(chan string),
	}

	go previewserver.Run(ch)

	preview.Run(ch)

}
