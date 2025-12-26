package main

import (
	"net/http"
	handler "go-lang-web-api/server/handler"
)
func main(){
	http.HandleFunc("/api/v1/analytics/genres", handler.GetTrends)
	http.ListenAndServe(":8000", nil)
}
