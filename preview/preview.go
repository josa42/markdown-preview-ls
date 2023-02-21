package preview

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/josa42/md-ls/control"
	"github.com/webview/webview"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var page = `
<!DOCTYPE html>
<html>
	<head>
		<title>%s</title>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.2.0/github-markdown.min.css" integrity="sha512-Ya9H+OPj8NgcQk34nCrbehaA0atbzGdZCI2uCbqVRELgnlrh8vQ2INMnkadVMSniC54HChLIh5htabVuKJww8g==" crossorigin="anonymous" referrerpolicy="no-referrer" />
		<style>
			.markdown-body {
				max-width: 830px;
				margin: 10px;
			}
		</style>
		<script>
		window.onload = () => {
		  [...document.querySelectorAll('a')].forEach((a) => {
				a.onclick = ((evt) => {
					evt.preventDefault();
					window.__handleNavigation(a.href ?? '')
				})
			})
		}
		</script>
	</head>
	<body class="markdown-body">%s</body>
</html>
`

func Run(ch control.PreviewChannels) {

	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Basic Example")
	w.SetSize(480, 320, webview.HintNone)

	w.SetHtml(render("Loading..."))

	go func() {
		for {
			text := <-ch.Update
			log.Println("<- update")
			w.SetHtml(render(text))
		}
	}()
	go func() {
		for {
			<-ch.Close
			log.Printf("<- close")
			w.Terminate()
		}
	}()

	w.Bind("__handleNavigation", func(link string) {
		fmt.Printf("navigation: '%s'\n", link)
		if strings.HasPrefix(link, "https://") || strings.HasPrefix(link, "http://") {
			go exec.Command("open", link).Run()
		}
	})

	w.Run()
}

func render(source string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		panic(err)
	}

	return fmt.Sprintf(page, "", string(buf.Bytes()))
}
