package control

import (
	"os"
	"path/filepath"
	"regexp"
)

var fileProtoExpr = regexp.MustCompile(`^file://`)

type File struct {
	FilePath string
	Source   string
}

func NewFile(filePath string, source string) File {
	return File{
		FilePath: fileProtoExpr.ReplaceAllString(filePath, ""),
		Source:   source,
	}
}

func (f File) RelFilePath() string {
	cwd, _ := os.Getwd()
	rel, _ := filepath.Rel(cwd, f.FilePath)

	return rel
}

type Channels struct {
	Open   chan File
	Close  chan bool
	Update chan File
}

func NewChannels() Channels {
	return Channels{
		Open:   make(chan File),
		Close:  make(chan bool),
		Update: make(chan File),
	}
}

type PreviewChannels struct {
	Close  chan bool
	Update chan File
}

func NewPreviewChannels() PreviewChannels {
	return PreviewChannels{
		Close:  make(chan bool),
		Update: make(chan File),
	}
}
