package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	modules "github.com/brycedouglasjames/yougoclient"
)

var (
	Refresh      sync.RWMutex
	AddUser      sync.RWMutex
	CurrentUsers []*modules.Users
	UserCache    []*modules.Users
	returnobj    []*modules.Respond
)

type RetrieveUserProfile struct {
	User string
}

func main() {

	go RefreshSearch()
	go CreateUser()

	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		Refresh.RLock()
		defer Refresh.RUnlock()

		switch r.Method {
		case "POST":
			//request buffer 100 KB
			r.Body = http.MaxBytesReader(w, r.Body, 100000)

			//create decoder
			decoder := json.NewDecoder(r.Body)

			//STRICT request scope
			decoder.DisallowUnknownFields()

			//init decoder
			query := &RetrieveUserProfile{}
			err := decoder.Decode(query)
			if err != nil {
				w.WriteHeader(401)
				w.Write([]byte(err.Error()))
				fmt.Println()
				return
			}

			//Makes sure there is only ONE json object
			err = decoder.Decode(&struct{}{})
			if err != io.EOF {
				msg := "Request body must only contain a single JSON object"
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			tempUser := &modules.Users{}
			for _, user := range CurrentUsers {
				if query.User == user.UserName {
					tempUser = user
				}
			}

			response, err := json.Marshal(tempUser)
			if err != nil {
				w.WriteHeader(401)
				w.Write([]byte(err.Error()))
				fmt.Println()
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		case "GET":
			response, err := json.Marshal(CurrentUsers)
			if err != nil {
				w.WriteHeader(401)
				w.Write([]byte(err.Error()))
				fmt.Println()
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(response))
		}
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		AddUser.RLock()
		defer AddUser.RUnlock()

		response, err := json.Marshal(CurrentUsers)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			fmt.Println()
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(response))
	})

	user := &modules.Users{}
	http.HandleFunc("/userpass", user.AddUser)

	videos := &modules.Users{}
	http.HandleFunc("/videos", videos.ServeArray)

	//handler for search results
	search := &modules.Users{}
	http.HandleFunc("/query", search.SearchHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func RefreshSearch() {
	for {
		Refresh.Lock()
		for _, user := range modules.ClientList {
			user.SessionTime = user.SessionTime + 1000
			if user.SessionTime >= 60000 {
				modules.DeleteUser(user.UserName)
			}
		}

		if modules.FinishedSearch == true {
			for k, user := range modules.ClientList {
				if user.UserName == modules.UserRequest {
					for _, item := range modules.UserSearch {
						var pass = true
						for _, video := range user.Searches {
							if video.VideoID == item.VideoID {
								pass = false
							}
						}

						if pass {
							data := &modules.Respond{
								VideoID:      item.VideoID,
								ThumbnailURL: item.ThumbnailURL,
								VideoTitle:   item.VideoTitle,
							}
							CurrentUsers[k].Searches = append(CurrentUsers[k].Searches, data)
						}
					}
				}
			}

			modules.CurrentSearch = ""
			modules.UserRequest = "..."
			modules.UserSearch = nil

			modules.FinishedSearch = false
		}

		Refresh.Unlock()
		time.Sleep(3 * time.Second)
	}
}

func CreateUser() {
	for {
		AddUser.Lock()

		CurrentUsers = modules.ClientList
		/*tempCurrentUsers := CurrentUsers
		CurrentUsers = nil

		for _, user := range modules.ClientList {
			var pass = true
			thisName := user.UserName
			for _, client := range tempCurrentUsers {
				if thisName == client.UserName {
					pass = false
				}
			}

			if pass {
				CurrentUsers = append(CurrentUsers, user)
			}
		}
		*/
		AddUser.Unlock()
		time.Sleep(1 * time.Second)

	}
}

//TODO user integration
/*func updateClient() {}
func stopClient()   {}


func updateSuer() {}*/
