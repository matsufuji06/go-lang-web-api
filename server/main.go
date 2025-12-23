package main

import (
	"go-lang-web-api/server/handler"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR loading .env:", err)
	}

	http.HandleFunc("/api/v1/videos", handler.GetVideos)
	http.ListenAndServe(":8000", nil)

}
