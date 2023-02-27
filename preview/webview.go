package preview

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/josa42/markdown-preview-ls/control"
	"github.com/josa42/markdown-preview-ls/render"
	"github.com/webview/webview"
)

func runWebView(ch control.PreviewChannels, port int, initialFile control.File) {
	w := webview.New(true)
	defer w.Destroy()
	w.SetSize(480, 320, webview.HintNone)

	currentFile := control.File{}

	update := func(file control.File) {
		if currentFile.FilePath != file.FilePath {
			currentFile = file
			relPath := file.RelFilePath()
			w.Dispatch(func() {
				w.SetTitle(relPath)
				w.Navigate(fmt.Sprintf("http://localhost:%d/%s.html", port, relPath))
			})
		} else {
			currentFile = file
			w.Eval("__update()")
		}
	}

	scroll := func(position control.ScrollPosition) {
		w.Eval(fmt.Sprintf("__scroll(%f)", position.Position))
	}

	update(initialFile)

	go func() {
		for {
			file := <-ch.Update
			update(file)
		}
	}()
	go func() {
		for {
			position := <-ch.Scroll
			scroll(position)
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
		return render.Body(currentFile.Source)
	})

	w.Run()
}
