package handler

import (
	"fmt"
	api "go-lang-web-api/server/model/api"
	service "go-lang-web-api/server/service"
	util "go-lang-web-api/server/service/util"
	"net/http"
)

// 国別のジャンルトレンドを取得
func GetTrends(w http.ResponseWriter, r *http.Request) {
	// クエリを取得
	country := r.URL.Query().Get("country")

	if country == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	// YouTube対応国コード取得
	regionCodes, err := service.FetchRegionCode()

	if err != nil {
		fmt.Println("FetchRegionCode error:", err)
		util.ResponseErrorJson(
			w,
			http.StatusInternalServerError,
			"An unexpected error occurred on the server.",
		)
		return
	}

	if _, ok := regionCodes[country]; ok {
		fmt.Printf("regionCode exists. The code is %#v", country)

	} else {
		util.ResponseErrorJson(
			w,
			http.StatusForbidden,
			"YouTube is not available in the selected country.",
		)
		return
	}

	popularGenres, err := service.FetchPopularGenres(country)
	if err != nil {
		fmt.Println("FetchPopularGenres error:", err)
		util.ResponseErrorJson(
			w,
			http.StatusInternalServerError,
			"An unexpected error occurred on the server.",
		)
		return
	}

	categories, err := service.FetchCategories(country)
	if err != nil {
		fmt.Println("FetchCategories error:", err)
		util.ResponseErrorJson(
			w,
			http.StatusInternalServerError,
			"An unexpected error occurred on the server.",
		)
		return
	}

	genres := make([]api.Genre, 0, len(popularGenres))

	for i := 0; i < len(popularGenres); i++ {
		if category, ok := categories[popularGenres[i].CategoryID]; ok {
			genres = append(genres, api.Genre{
				Genre:     category,
				Ratio:     popularGenres[i].Ratio,
				TopVideos: popularGenres[i].VideoURL,
			})
		}
	}
	util.ResponseJson(
		w,
		http.StatusOK,
		api.Response{
			Country: country,
			Genres:  genres,
		},
	)
	
}
