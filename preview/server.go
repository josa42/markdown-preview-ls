package preview

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/josa42/markdown-preview-ls/control"
	"github.com/josa42/markdown-preview-ls/render"
)

const (
	UpdateCommand = "/____api____/update"
	CloseCommand  = "/____api____/close"
)

var htmlExpr = regexp.MustCompile(`\.html$`)

func runServer(ch control.PreviewChannels, port int, initialFile control.File) {

	currentFile := initialFile
	assetsServer := http.FileServer(http.Dir("."))

	http.HandleFunc(CloseCommand, func(w http.ResponseWriter, r *http.Request) {
		ch.Close <- true
		io.WriteString(w, "ok")
	})

	http.HandleFunc(UpdateCommand, func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		file := control.File{}
		json.Unmarshal(body, &file)

		currentFile = file
		ch.Update <- file

		io.WriteString(w, "ok")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if htmlExpr.MatchString(r.URL.Path) {
			io.WriteString(w, render.Page(currentFile.Source))
			return
		}

		assetsServer.ServeHTTP(w, r)
	})

	http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}
