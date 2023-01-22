package structs

type BingNewsReachResult struct {
	Value []BingNewsRow `json:"value"`
}

type BingNewsRow struct {
	Name  string        `json:"name"`
	Url   string        `json:"url"`
	Image BingNewsImage `json:"image"`
}

type BingNewsImage struct {
	Thumbnail BingNewsThumbnail `json:"thumbnail"`
}

type BingNewsThumbnail struct {
	ContentUrl string `json:"contentUrl"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}
