package main

import (
	"log"
	"os"

	"github.com/iagocanalejas/buzz-client/pkg/tui"
)

func main() {
	setupFileLogger()
	if err := tui.BuildApp().App.Run(); err != nil {
		panic(err)
	}
}

func setupFileLogger() {
	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	log.SetOutput(f)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
