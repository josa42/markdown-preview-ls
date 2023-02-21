package server

import (
	"log"

	"github.com/josa42/go-ls"
	"github.com/josa42/md-ls/logger"
	"github.com/sourcegraph/go-lsp"
)

func Run(preview func(txt string)) {
	defer logger.Init("/tmp/md-ls.log")()

	s := ls.New()
	s.VerboseLogging = true

	s.Root.Initialize(Initialize)

	s.TextDocument.DidOpen(func(ctx ls.RequestContext, p lsp.DidOpenTextDocumentParams) error {
		ctx.Server.State.SetDocument(p.TextDocument)

		text, _ := ctx.Server.State.GetText(p.TextDocument.URI)
		preview(text)

		return nil
	})

	s.TextDocument.DidChange(func(ctx ls.RequestContext, p lsp.DidChangeTextDocumentParams) error {
		ctx.Server.State.ApplyCanges(p.TextDocument.URI, p.ContentChanges)

		text, _ := ctx.Server.State.GetText(p.TextDocument.URI)
		preview(text)

		return nil
	})

	s.TextDocument.DidClose(func(ctx ls.RequestContext, p lsp.DidCloseTextDocumentParams) error {
		ctx.Server.State.Remove(p.TextDocument.URI)
		return nil
	})

	if err := s.StartAndWait(); err != nil {
		log.Printf("Server exited: %v", err)
	}
}
