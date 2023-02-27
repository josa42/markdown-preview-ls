package main

import (
	"os"
	"strconv"

	"github.com/josa42/markdown-preview-ls/control"
	"github.com/josa42/markdown-preview-ls/logger"
	"github.com/josa42/markdown-preview-ls/preview"
	"github.com/josa42/markdown-preview-ls/server"
)

func main() {
	cmd := "server"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "preview":
		// TODO CLI based on flag options
		port := 0
		filePath := ""
		source := ""
		if len(os.Args) > 2 {
			filePath = os.Args[2]
		}
		if len(os.Args) > 3 {
			source = os.Args[3]
		}
		if len(os.Args) > 4 {
			port, _ = strconv.Atoi(os.Args[4])
		}

		runPreview(port, control.NewFile(filePath, source))
	default:
		runServer()
	}
}

func runServer() {
	defer logger.Init("/tmp/markdown-preview-ls.log")()
	ch := control.NewChannels()

	client := preview.Client{}

	go func() {
		for {
			file := <-ch.Open
			if !client.IsRunning() {
				client.Execute(file)
			} else {
				client.Update(file)
			}
		}
	}()
	go func() {
		for {
			<-ch.Close
			client.Close()
		}
	}()

	go func() {
		for {
			position := <-ch.Scroll
			client.Scroll(position)
		}
	}()

	go func() {
		for {
			file := <-ch.Update
			client.Update(file)
		}
	}()

	server.Run(ch)
}

func runPreview(port int, initialFile control.File) {
	defer logger.Init("/tmp/markdown-preview-ls.preview.log")()

	preview.Run(port, initialFile)
}
