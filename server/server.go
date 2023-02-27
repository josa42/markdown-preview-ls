package server

import (
	"errors"
	"log"

	"github.com/josa42/go-ls"
	"github.com/josa42/go-ls/lsp"
	"github.com/josa42/go-ls/utils"
	"github.com/josa42/markdown-preview-ls/control"
)

const (
	CMD_OPEN   = "mardown-preview.open"
	CMD_UPDATE = "mardown-preview.update"
	CMD_CLOSE  = "mardown-preview.close"
	CMD_SCROLL = "mardown-preview.scroll"
)

var currentURI lsp.DocumentURI

type UpdateParams struct {
	TextDocument lsp.TextDocumentIdentifier `json:"textDocument"`
}

type ScrollParams struct {
	TextDocument lsp.TextDocumentIdentifier `json:"textDocument"`
	Position     float32                    `json:"position"`
}

func Run(ch control.Channels) {
	s := ls.New()
	s.VerboseLogging = true

	s.TextDocument.RegisterChangesListener()
	s.State.OnChange(func(uri lsp.DocumentURI) {
		if uri == currentURI {
			text, _ := s.State.GetText(uri)
			ch.Update <- control.NewFile(string(uri), text)
		}

	})

	s.Root.Shutdown(func(ctx ls.RequestContext) error {
		go func() { ch.Close <- true }()
		return nil
	})

	s.Workspace.RegisterCommand(CMD_OPEN, func(ctx ls.RequestContext, args []interface{}) error {
		if len(args) != 1 {
			return errors.New("argument is required")
		}

		params, _ := utils.Unmarkshal[UpdateParams](args[0])
		currentURI = params.TextDocument.URI
		text, _ := ctx.Server.State.GetText(params.TextDocument.URI)

		ch.Open <- control.NewFile(string(currentURI), text)

		return nil
	})

	s.Workspace.RegisterCommand(CMD_UPDATE, func(ctx ls.RequestContext, args []interface{}) error {
		if len(args) != 1 {
			return errors.New("argument is required")
		}

		params, _ := utils.Unmarkshal[UpdateParams](args[0])
		currentURI = params.TextDocument.URI
		text, _ := ctx.Server.State.GetText(params.TextDocument.URI)

		ch.Update <- control.NewFile(string(currentURI), text)

		return nil
	})

	s.Workspace.RegisterCommand(CMD_CLOSE, func(ctx ls.RequestContext, args []interface{}) error {
		go func() { ch.Close <- true }()

		return nil
	})

	s.Workspace.RegisterCommand(CMD_SCROLL, func(ctx ls.RequestContext, args []interface{}) error {
		if len(args) != 1 {
			return errors.New("argument is required")
		}

		params, _ := utils.Unmarkshal[ScrollParams](args[0])

		if currentURI == params.TextDocument.URI {
			ch.Scroll <- control.NewScrollPosition(string(currentURI), params.Position)
		}

		return nil
	})

	s.Start()

	if err := s.Wait(); err != nil {
		log.Printf("Server exited: %v", err)
	}
}
