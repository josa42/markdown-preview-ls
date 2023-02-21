package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var output io.Writer = os.Stderr

func Init(logFile string) func() {
	f, err := os.OpenFile(fmt.Sprintf(logFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	output = f

	log.SetOutput(f)

	return func() { f.Close() }
}

func New() *log.Logger {
	return log.New(output, "", log.LstdFlags)
}
