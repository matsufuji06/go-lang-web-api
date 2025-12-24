package youtube

type YouTubeChannelsResponse struct {
	Items []struct {
		Snippet struct {
			Title      string `json:"title"`
			Thumbnails struct {
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		Statistics struct {
			SubscriberCount string `json:"subscriberCount"`
			ViewCount       string `json:"viewCount"`
			VideoCount      string `json:"videoCount"`
		} `json:"statistics"`
	} `json:"items"`
}
