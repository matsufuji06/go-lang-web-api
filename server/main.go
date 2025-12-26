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

	// 静的ファイル(JSファイル)の配信設定
	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

	// index.htmlの配信設定
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	http.HandleFunc("/api/v1/videos", handler.GetVideos)
	http.ListenAndServe(":8000", nil)

}
