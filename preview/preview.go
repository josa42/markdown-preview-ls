package preview

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
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

func Run() {
	currentFileName := "README.md"

	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Basic Example")
	w.SetSize(480, 320, webview.HintNone)

	update := func() {
		c, _ := ioutil.ReadFile(currentFileName)
		w.SetHtml(fmt.Sprintf(page, currentFileName, render(c)))
	}

	go watch(func(fileName string) {
		if fileName == currentFileName {
			update()
		}
	})

	w.Bind("__handleNavigation", func(link string) {
		fmt.Printf("navigation: '%s'\n", link)
		if strings.HasPrefix(link, "https://") || strings.HasPrefix(link, "http://") {
			go exec.Command("open", link).Run()
		} else if _, err := os.Stat(link); err == nil {
			currentFileName = filepath.Clean(link)
			update()
		}
	})

	update()

	w.Run()
}

func watch(fn func(fileName string)) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				fn(event.Name)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func render(source []byte) string {
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
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}

	return string(buf.Bytes())
}
