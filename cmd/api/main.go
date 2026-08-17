package main

import (
	"finances/internal/router"
	"log"
	"net/http"
	"time"
)

func main() {

	router := router.GenerateRouter()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Iniciando a API na porta: %s", "8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
