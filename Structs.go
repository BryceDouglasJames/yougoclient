package yougoclient

var (
	CurrentSearch      string
	CurrentSearchIndex = 0
	ClientList         []*Users
	PASSFLAG           = 0
	UserSearch         []*Respond
	UserRequest        string
	AmountOfUsers      = 0
	FinishedSearch     = false
)

//User struct for user presets
type Users struct {
	UserIndex   int        `json:"UserIndex"`
	UserName    string     `json:"UserName"`
	Searches    []*Respond `json:"Searches"`
	SessionTime int        `json:"SessionTime"`
}

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

//SearchRequest used for user search
type SearchRequest struct {
	ID   string
	User string
}

type UserAdd struct {
	ID string
}
