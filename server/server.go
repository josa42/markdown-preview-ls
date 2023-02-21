package server

import (
	"log"

	"github.com/josa42/go-ls"
	"github.com/josa42/md-ls/control"
	"github.com/josa42/md-ls/logger"
	"github.com/sourcegraph/go-lsp"
)

func Run(ch control.Channels) {
	defer logger.Init("/tmp/md-ls.log")()

	s := ls.New()
	s.VerboseLogging = true

	s.Root.Initialize(Initialize)

	s.Workspace.ExecuteCommand(func(ctx ls.RequestContext, p lsp.ExecuteCommandParams) error {
		switch p.Command {
		case "openPreview":
			ch.Open <- true
		case "closePreview":
			ch.Open <- false
		}
		return nil
	})

	s.TextDocument.DidOpen(func(ctx ls.RequestContext, p lsp.DidOpenTextDocumentParams) error {
		ctx.Server.State.SetDocument(p.TextDocument)

		text, _ := ctx.Server.State.GetText(p.TextDocument.URI)
		ch.Update <- text

		return nil
	})

	s.TextDocument.DidChange(func(ctx ls.RequestContext, p lsp.DidChangeTextDocumentParams) error {
		ctx.Server.State.ApplyCanges(p.TextDocument.URI, p.ContentChanges)

		text, _ := ctx.Server.State.GetText(p.TextDocument.URI)
		ch.Update <- text

		return nil
	})

	s.TextDocument.DidClose(func(ctx ls.RequestContext, p lsp.DidCloseTextDocumentParams) error {
		ctx.Server.State.Remove(p.TextDocument.URI)
		return nil
	})

	s.Start()

	ch.Started <- true

	if err := s.Wait(); err != nil {
		log.Printf("Server exited: %v", err)
	}
}
