package main

import (
	"log"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/app"
)

func main() {
	srv := new(app.Server)

	if err := srv.Run(); err != nil {
		log.Fatal("Error occurred while starting server:", err.Error())
	}
}
