package main

import (
	"fmt"
	"go-lang-web-api/server/handler"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR loading .env:", err)
	}
	fmt.Println("API KEY =", os.Getenv("YOUTUBE_API_KEY"))

	http.HandleFunc("/api/v1/videos", handler.GetVideos)
	http.ListenAndServe(":8000", nil)

}
