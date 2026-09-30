package main

import (
	"log"
	"os"
	"server/app/command"
)

func main() {
	if err := command.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
