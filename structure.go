package yougoclient

//User struct for user presets
type Users struct {
	UserName []string
	Searches []*Respond
}

//func (h *Users) AddUser() {}

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

func (h *Respond) ClearResponse() {
	h.VideoID = ""
	h.ThumbnailURL = ""
	h.VideoTitle = ""
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

/*//SearchQuery Something
type Query struct {
	search *string
}

func (h *Query) ChangeQuery(query *string) *Query {
	h.search = query
	return h
}

func (h *Query) GetQuery() *Query {
	return h
}

type YouService struct {
	service *youtube.Service
}

func (h *YouService) SetService(s *youtube.Service) *YouService {
	h.service = s
	return h
}

func (h *YouService) GetService() *youtube.Service {
	return h.service
}

type ThisClient struct {
	client *http.Client
}

func (h *ThisClient) SetClient(c *http.Client) *ThisClient {
	h.client = c
	return h
}

func (h *ThisClient) GetClient() *http.Client {
	return h.client
}*/
