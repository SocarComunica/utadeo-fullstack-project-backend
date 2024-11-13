package main

import (
	"log"

	"backend/src"
)

func main() {
	if err := src.Run(); err != nil {
		log.Fatalf("error while starting application: %v", err)
	}
}
