module github.com/josa42/md-ls

go 1.19

require (
	github.com/josa42/go-ls v0.0.0-20220605165229-063fde83b2a6
	github.com/sourcegraph/go-lsp v0.0.0-20181119182933-0c7d621186c1
	github.com/webview/webview v0.0.0-20230210061304-7b40e46d97e9
	github.com/yuin/goldmark v1.5.4
)

require (
	bitbucket.org/creachadair/stringset v0.0.8 // indirect
	github.com/creachadair/jrpc2 v0.4.5 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7 // indirect
)

replace (
  github.com/josa42/go-ls => ../go-ls
)
