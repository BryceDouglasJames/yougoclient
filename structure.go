package yougoclient

//User struct for user presets
type User struct {
	UserName string      `json:"username"`
	Searches []*Response `json:"searches"`
}

//Response used for storing video information
type Response struct {
	VideoID      string `json:"id"`
	ThumbnailURL string `json:"url"`
	VideoTitle   string `json:"title"`
}

//SearchRequest used for user search
type SearchRequest struct {
	ID string `json:"id"`
}
