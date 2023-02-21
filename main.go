package main

import (
	"log"

	"github.com/josa42/md-ls/preview"
	"github.com/josa42/md-ls/server"
)

func main() {
	run := make(chan bool)
	textUpdate := make(chan string)

	go server.Run(func(text string) {
		log.Println("write")
		run <- true
		textUpdate <- text
	})

	preview.Run(textUpdate, run)
}
