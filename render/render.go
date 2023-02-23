package render

import (
	"bytes"
	"fmt"

	chromehtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/mermaid"
)

func Body(source string) string {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			&mermaid.Extender{},
			emoji.Emoji,
			meta.New(meta.WithTable()),
			highlighting.NewHighlighting(
				// highlighting.WithStyle("github-dark"),
				highlighting.WithStyle("onedark"),
				highlighting.WithFormatOptions(
					chromehtml.WithLineNumbers(true),
				),
			),
		),

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

	return string(buf.Bytes())
}

func Page(source string) string {
	return fmt.Sprintf(page, "", Body(source))
}
