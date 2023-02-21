package previewserver

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/josa42/markdown-preview-ls/control"
)

func Run(ch control.PreviewChannels) {

	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		ch.Close <- true
		io.WriteString(w, "ok")
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		ch.Update <- string(body)

		io.WriteString(w, "ok")
	})

	http.ListenAndServe("localhost:3333", nil)
}
