package yougoclient

//User struct for user presets
type Users struct {
	UserName []string
	Searches []*Respond
}

func (h *Users) AddUser() {}

func (h *Users) AddVideo(data *Respond) *Users {
	h.Searches = append(h.Searches, data)
	return h
}

//Respond used for storing video information
type Respond struct {
	VideoID      string
	ThumbnailURL string
	VideoTitle   string
}

func (h *Respond) SetResponse(vid string, pic string, title string) {
	h.VideoID = vid
	h.ThumbnailURL = pic
	h.VideoTitle = title
}

//SearchRequest used for user search
type SearchRequest struct {
	ID string
}

func (h *SearchRequest) SetRequest(name string) *SearchRequest {
	h.ID = name
	return h
}

func (h *SearchRequest) GetRequest() string {
	return h.ID
}
