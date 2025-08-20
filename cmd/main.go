package main

import (
	"os"
	"sca/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../config/.env")

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	server.Init(host, port)
}
