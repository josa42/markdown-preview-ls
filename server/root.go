package server

import (
	"context"

	"github.com/josa42/go-ls"
	"github.com/josa42/go-ls/lsp"
	golsp "github.com/sourcegraph/go-lsp"
)

func Initialize(s ls.Server, ctx context.Context, p lsp.InitializeParams) (lsp.InitializeResult, error) {

	return lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &golsp.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    golsp.TDSKFull,
			},

			CompletionProvider: &golsp.CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{" "},
			},

			DocumentFormattingProvider: true,
			HoverProvider:              true,
			DocumentSymbolProvider:     true,
			FoldingRangeProvider:       true,
			// DocumentHighlightProvider:  true,
			// SignatureHelpProvider: &lsp.SignatureHelpOptions{
			// 	TriggerCharacters: []string{},
			// },
		},
	}, nil
}
