package api

type Response struct {
	Country   string `json:"country"`
	Genres []Genre  `json:"genres"`
}

type Genre struct {
	Genre string `json:"genre"`
	Ratio float64 `json:"ratio"`
	TopVideos []VideoURL `json:"topVideos"`
}


type ErrorResponse struct {
	Message string `json:"message"`
}

type GenreRatio struct {
	CategoryID string  `json:"categoryID"`
	Ratio      float64 `json:"ratio"`
	VideoURL      []VideoURL `json:"video"`
}

type VideoURL struct {
	VideoURL   string `json:"url"`
}