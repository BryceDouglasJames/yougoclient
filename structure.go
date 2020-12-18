package yougoclient

type Response struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//SearchRequest used for user search
type SearchRequest struct {
	ID string `json:"id"`
}
