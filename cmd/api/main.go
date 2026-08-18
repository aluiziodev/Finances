package main

import (
	"finances/internal/connection"
	"finances/internal/router"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if err := connection.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer connection.DB.Close()

	router := router.GenerateRouter()

	port := os.Getenv("API_PORT")

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Iniciando a API na porta: %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
