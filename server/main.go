package main

import (
	handler "go-lang-web-api/server/handler"
	utils "go-lang-web-api/server/service/utils"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 1)

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			utils.ResponseErrorJson(w, http.StatusTooManyRequests, "API rate limit exceeded. Please try again later")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR loading .env:", err)
	}

	// メインの mux（HTML / static 用）
	mux := http.NewServeMux()

	// 静的ファイル
	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

	// index.html
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	// --- API 用 mux ---
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/v1/analytics/genres", handler.GetTrends)
	apiMux.HandleFunc("/api/v1/channels", handler.GetChannelInfo)

	// API にだけレート制限を適用
	mux.Handle("/api/", rateLimitMiddleware(apiMux))

	http.ListenAndServe(":8000", mux)
}
