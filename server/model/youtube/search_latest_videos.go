package youtube

type SearchResponse struct {
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	ID      VideoID `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type VideoID struct {
	VideoID string `json:"videoId"`
}

type Snippet struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	PublishedAt  string     `json:"publishedAt"`
	ChannelID    string     `json:"channelId"`
	ChannelTitle string     `json:"channelTitle"`
	Thumbnails   Thumbnails `json:"thumbnails"`
}

type Thumbnails struct {
	Default Thumbnail `json:"default"`
	Medium  Thumbnail `json:"medium"`
	High    Thumbnail `json:"high"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
