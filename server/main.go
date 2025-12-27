package main

import (
<<<<<<< HEAD
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
=======
	handler "go-lang-web-api/server/handler"
	"net/http"
	"golang.org/x/time/rate"
		utils "go-lang-web-api/server/service/utils"
)
var limiter = rate.NewLimiter(1, 1) 
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			utils.ResponseErrorJson(w, http.StatusTooManyRequests,"API rate limit exceeded. Please try again later" )
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// メインの mux（HTML / static 用）
	mux := http.NewServeMux()

	// 静的ファイル
	mux.Handle(
>>>>>>> 23e4720dc6159ebb9ccb7de2dad5870bdfef61d7
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

<<<<<<< HEAD
	// index.htmlの配信設定
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	fmt.Println("Server is running!")
	http.HandleFunc("/api/v1/channels", handler.GetChannelInfo)
	http.ListenAndServe(":8000", nil)

=======
	// index.html
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	// --- API 用 mux ---
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/v1/analytics/genres", handler.GetTrends)

	// API にだけレート制限を適用
	mux.Handle("/api/", rateLimitMiddleware(apiMux))

	http.ListenAndServe(":8000", mux)
>>>>>>> 23e4720dc6159ebb9ccb7de2dad5870bdfef61d7
}
