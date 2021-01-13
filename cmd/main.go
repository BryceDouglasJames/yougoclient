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

	//START CONCURRENT ROUTINES
	go RefreshSearch()
	go CreateUser()

	//run file server
	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	//handler for responding to search queries
	http.HandleFunc("/userbank/1999", func(w http.ResponseWriter, r *http.Request) {
		Refresh.RLock()
		defer Refresh.RUnlock()

		switch r.Method {
		//POSTS USER UPDATES
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

		//GET ALL USERS 			MAYBE CHANGE THIS
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

	//handles new session clients
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

//Concurrent call for pushing searches to called user. Very convoluted.
//This function is flawed. Query collisions WILL occur if there are 2 requests simultaneously.
//Needs better solution
func RefreshSearch() {
	for {
		Refresh.Lock()

		//delete user if their ttl is expired
		for _, user := range modules.ClientList {
			user.SessionTime = user.SessionTime + 1000
			if user.SessionTime >= 60000 {
				modules.DeleteUser(user.UserName)
			}
		}

		//loops through and sends results back to selected user.
		//Could very well be optimized this method is overly forced.
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

			//clear global buffers
			modules.CurrentSearch = ""
			modules.UserRequest = "..."
			modules.UserSearch = nil
			modules.FinishedSearch = false
		}

		//open channel
		Refresh.Unlock()
		time.Sleep(3 * time.Second)
	}
}

//concurrent function for
func CreateUser() {
	for {
		AddUser.Lock()

		CurrentUsers = modules.ClientList
		AddUser.Unlock()
		time.Sleep(1 * time.Second)

	}
}
