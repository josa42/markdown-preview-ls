package render

import (
	"bytes"
	"fmt"

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
				margin: 10px 20px;
			}
		</style>
		<script>
		  function __preventNavigation() {
				[...document.querySelectorAll('a')].forEach((a) => {
					a.onclick = ((evt) => {
						evt.preventDefault();
						window.__handleNavigation(a.href ?? '')
					})
				});
			}
		async function __update() {
			document.querySelector('body').innerHTML = await __getText()
			__preventNavigation()
		}
		window.onload = () => __preventNavigation()
		</script>
	</head>
	<body class="markdown-body">%s</body>
</html>
`

func Body(source string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote),
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
