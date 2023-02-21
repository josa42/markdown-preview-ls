package main

import (
	"github.com/josa42/md-ls/control"
	"github.com/josa42/md-ls/preview"
	"github.com/josa42/md-ls/server"
)

func main() {
	ch := control.Channels{
		Open:    make(chan bool),
		Update:  make(chan string),
		Started: make(chan bool),
	}

	// open := make(chan bool)
	// textUpdate := make(chan string)

	go server.Run(ch)

	<-ch.Started

	// go server.Run(func(chan) {
	// 	log.Println("write")
	// 	open <- true
	// 	textUpdate <- text
	// }, open)

	preview.Run(ch)
}
