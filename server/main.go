package main

import (
	"fmt"
	"go-lang-web-api/server/handler"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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

	fmt.Println("Server is running!")
	http.HandleFunc("/api/v1/channels", handler.GetChannelInfo)
	http.ListenAndServe(":8000", nil)

}
