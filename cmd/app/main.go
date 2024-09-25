package main

import (
	"emporium/internal/app"
	"log"
)

func main() {
	srv := new(app.Server)

	if err := srv.Run(); err != nil {
		log.Fatal("Error occurred while starting server:", err.Error())
	}
}
