package preview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/josa42/markdown-preview-ls/control"
	"github.com/josa42/markdown-preview-ls/ports"
)

type Client struct {
	port int
}

func (p Client) IsRunning() bool {
	return p.port > 0
}

func (p *Client) Execute(file control.File) {
	port, _ := ports.GetFreePort()
	cmd := exec.Command(os.Args[0], "preview", file.FilePath, file.Source, fmt.Sprintf("%d", port))
	go func() {
		cmd.Run()
		port = 0
	}()

	p.port = port
}

func (p Client) Update(file control.File) {
	if p.IsRunning() {
		res, _ := json.Marshal(file)
		http.Post(p.url(UpdateCommand), "application/json", bytes.NewBuffer(res))
	}
}

func (p Client) Scroll(position control.ScrollPosition) {
	if p.IsRunning() {
		res, _ := json.Marshal(position)
		http.Post(p.url(ScrollCommand), "application/json", bytes.NewBuffer(res))
	}
}

func (p Client) Close() {
	if p.IsRunning() {
		http.Post(p.url(CloseCommand), "text/plain", nil)
	}
}

func (p Client) url(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", p.port, path)
}
