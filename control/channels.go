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

type ScrollPosition struct {
	FilePath string
	Position float32
}

func NewScrollPosition(filePath string, position float32) ScrollPosition {
	return ScrollPosition{
		FilePath: fileProtoExpr.ReplaceAllString(filePath, ""),
		Position: position,
	}
}

type Channels struct {
	Open   chan File
	Update chan File
	Scroll chan ScrollPosition
	Close  chan bool
}

func NewChannels() Channels {
	return Channels{
		Open:   make(chan File),
		Update: make(chan File),
		Scroll: make(chan ScrollPosition),
		Close:  make(chan bool),
	}
}

type PreviewChannels struct {
	Close  chan bool
	Update chan File
	Scroll chan ScrollPosition
}

func NewPreviewChannels() PreviewChannels {
	return PreviewChannels{
		Close:  make(chan bool),
		Update: make(chan File),
		Scroll: make(chan ScrollPosition),
	}
}
