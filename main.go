package main

import (
	"log"

	"github.com/MOZGIII/external-dns-endpoint-server/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
