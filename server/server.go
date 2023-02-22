package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/josa42/go-ls"
	"github.com/josa42/go-ls/lsp"
	"github.com/josa42/markdown-preview-ls/control"
)

const (
	CMD_OPEN   = "mardown-preview.open"
	CMD_UPDATE = "mardown-preview.update"
	CMD_CLOSE  = "mardown-preview.close"
)

var currentURI lsp.DocumentURI

func unmarshalUpdateParams(input interface{}) UpdateParams {
	params := UpdateParams{}
	jsonData, _ := json.Marshal(input)
	json.Unmarshal(jsonData, &params)

	return params
}

type UpdateParams struct {
	TextDocument lsp.TextDocumentIdentifier `json:"textDocument"`
}

func Run(ch control.Channels) {
	s := ls.New()
	s.VerboseLogging = true

	s.State.OnChange(func(uri lsp.DocumentURI) {
		if uri == currentURI {
			text, _ := s.State.GetText(uri)
			ch.Update <- control.NewFile(string(uri), text)
		}

	})

	s.Root.Initialize(func(s ls.Server, ctx context.Context, p lsp.InitializeParams) (lsp.InitializeResult, error) {
		return lsp.InitializeResult{
			Capabilities: lsp.ServerCapabilities{
				TextDocumentSync: &lsp.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    lsp.TDSKFull,
				},
				ExecuteCommandProvider: &lsp.ExecuteCommandOptions{Commands: []string{CMD_OPEN, CMD_UPDATE, CMD_CLOSE}},
			},
		}, nil
	})

	s.Root.Shutdown(func(ctx ls.RequestContext) error {
		go func() { ch.Close <- true }()
		return nil
	})

	s.Workspace.ExecuteCommand(func(ctx ls.RequestContext, p lsp.ExecuteCommandParams) error {
		switch p.Command {
		case CMD_OPEN:
			go func() {
				if len(p.Arguments) == 1 {
					params := unmarshalUpdateParams(p.Arguments[0])
					currentURI = params.TextDocument.URI
					text, _ := ctx.Server.State.GetText(params.TextDocument.URI)

					ch.Open <- control.NewFile(string(currentURI), text)
				}
			}()

		case CMD_UPDATE:
			go func() {
				if len(p.Arguments) == 1 {
					params := unmarshalUpdateParams(p.Arguments[0])
					currentURI = params.TextDocument.URI
					text, _ := ctx.Server.State.GetText(params.TextDocument.URI)

					ch.Update <- control.NewFile(string(currentURI), text)
				}
			}()

		case CMD_CLOSE:
			go func() { ch.Close <- true }()
		}

		return nil
	})

	s.TextDocument.DidOpen(func(ctx ls.RequestContext, p lsp.DidOpenTextDocumentParams) error {
		ctx.Server.State.SetDocument(p.TextDocument)
		return nil
	})

	s.TextDocument.DidChange(func(ctx ls.RequestContext, p lsp.DidChangeTextDocumentParams) error {
		ctx.Server.State.ApplyCanges(p.TextDocument.URI, p.ContentChanges)
		return nil
	})

	s.TextDocument.DidClose(func(ctx ls.RequestContext, p lsp.DidCloseTextDocumentParams) error {
		ctx.Server.State.Remove(p.TextDocument.URI)
		return nil
	})

	s.Start()

	if err := s.Wait(); err != nil {
		log.Printf("Server exited: %v", err)
	}
}
