package api

type VideoResponse struct {
	Items []VideoItem `json:"items"`
	Meta  Meta        `json:"meta"`
}

type VideoItem struct {
	VideoID     string     `json:"videoId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	PublishedAt string     `json:"publishedAt"`
	ChannelID   string     `json:"channelId"`
	ChannelName string     `json:"channelName"`
	URL         string     `json:"url"`
	IsNew       bool       `json:"isNew"`
	Thumbnails  Thumbnails `json:"thumbnails"`
}

type Meta struct {
	Count           int    `json:"count"`
	LastPublishedAt string `json:"lastPublishedAt"`
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
