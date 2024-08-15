package main

import (
	"fmt"

	"github.com/ccallazans/ai-video-generator/internal/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
