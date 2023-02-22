package preview

import (
	"github.com/josa42/markdown-preview-ls/control"
)

func Run(port int, initialFile control.File) {
	ch := control.NewPreviewChannels()

	go runServer(ch, port, initialFile)
	runWebView(ch, port, initialFile)
}
