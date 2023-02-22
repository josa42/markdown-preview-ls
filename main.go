package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/josa42/markdown-preview-ls/control"
	"github.com/josa42/markdown-preview-ls/logger"
	"github.com/josa42/markdown-preview-ls/ports"
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

	previewPort := 0

	url := func(path string) string {
		return fmt.Sprintf("http://localhost:%d%s", previewPort, path)
	}

	go func() {
		for {
			file := <-ch.Open
			if previewPort == 0 {
				previewPort, _ = ports.GetFreePort()
				cmd := exec.Command(os.Args[0], "preview", file.FilePath, file.Source, fmt.Sprintf("%d", previewPort))
				go func() {
					cmd.Run()
					previewPort = 0
				}()
			} else {
				res, _ := json.Marshal(file)
				http.Post(url(preview.UpdateCommand), "application/json", bytes.NewBuffer(res))

			}
		}
	}()
	go func() {
		for {
			<-ch.Close
			if previewPort > 0 {
				http.Post(url(preview.CloseCommand), "text/plain", nil)
			}
		}
	}()

	go func() {
		for {
			file := <-ch.Update
			res, _ := json.Marshal(file)
			http.Post(url(preview.UpdateCommand), "application/json", bytes.NewBuffer(res))
		}
	}()

	server.Run(ch)
}

func runPreview(port int, initialFile control.File) {
	defer logger.Init("/tmp/markdown-preview-ls.preview.log")()

	preview.Run(port, initialFile)
}
