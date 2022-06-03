package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Start...")

	app := NewApp(os.Stdout)
	_ = app.Run(os.Args)
}
