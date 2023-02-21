package preview

import (
	"os/exec"
	"strings"

	"github.com/josa42/markdown-preview-ls/control"
	"github.com/josa42/markdown-preview-ls/render"
	"github.com/webview/webview"
)

func Run(ch control.PreviewChannels, initialSource string) {

	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Markdown Preview")
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml(render.Page(initialSource))

	currentSource := initialSource

	go func() {
		for {
			currentSource = <-ch.Update
			w.Eval("__update()")
		}
	}()
	go func() {
		for {
			<-ch.Close
			w.Terminate()
		}
	}()

	w.Bind("__handleNavigation", func(link string) {
		if strings.HasPrefix(link, "https://") || strings.HasPrefix(link, "http://") {
			go exec.Command("open", link).Run()
		}
	})

	w.Bind("__getText", func() string {
		return render.Body(currentSource)
	})

	w.Run()
}
