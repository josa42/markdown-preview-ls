package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/josa42/markdown-preview-ls/control"
	"github.com/josa42/markdown-preview-ls/ports"
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
		source := ""
		port := 0
		if len(os.Args) > 2 {
			source = os.Args[2]
		}
		if len(os.Args) > 3 {
			port, _ = strconv.Atoi(os.Args[3])
		}

		runPreview(port, source)
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

	previewPort := 0

	url := func(path string) string {
		return fmt.Sprintf("http://localhost:%d/%s", previewPort, path)
	}

	go func() {
		for {
			text := <-ch.Open
			if previewPort == 0 {
				previewPort, _ = ports.GetFreePort()
				cmd := exec.Command(os.Args[0], "preview", text, fmt.Sprintf("%d", previewPort))
				go func() {
					cmd.Run()
					previewPort = 0
				}()
			} else {
				http.Post(url("update"), "text/plain", bytes.NewBufferString(text))

			}
		}
	}()
	go func() {
		for {
			<-ch.Close
			if previewPort > 0 {
				http.Post(url("close"), "text/plain", nil)
			}
		}
	}()

	go func() {
		for {
			text := <-ch.Update
			http.Post(url("update"), "text/plain", bytes.NewBufferString(text))
		}
	}()

	server.Run(ch)
}

func runPreview(port int, initialSource string) {
	ch := control.PreviewChannels{
		Close:  make(chan bool),
		Update: make(chan string),
	}

	go previewserver.Run(port, ch)

	preview.Run(ch, initialSource)
}
