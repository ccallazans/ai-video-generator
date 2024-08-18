package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ccallazans/ai-video-generator/internal/server"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/gommon/log"
)

// Just in case want to use python venv
func run() {
	if os.Getenv("ENV") == "local" {
		cmd := exec.Command("bash", "-c", "source ./venv/bin/activate && echo 'Virtual environment activated'")

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Error activating virtual env: %s", err)
			return
		}

		log.Printf("Output: %s", output)
		return
	}
}

func main() {
	server := server.NewServer()

	log.Printf("Server running on port %s", os.Getenv("PORT"))
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
