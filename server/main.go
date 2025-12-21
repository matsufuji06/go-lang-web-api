package main

import (
	"go-lang-web-api/server/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/videos", handler.GetVideos)
	http.ListenAndServe(":8000", nil)

}
