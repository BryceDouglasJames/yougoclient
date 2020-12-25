package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	modules "github.com/brycedouglasjames/yougoclient"
)

var (
	Refresh sync.RWMutex
	AddUser sync.RWMutex

	CurrentUsers []*modules.Users
	UserCache    []*modules.Users
	returnobj    []*modules.Respond
)

func main() {

	go RefreshSearch()
	go CreateUser()

	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		Refresh.RLock()
		defer Refresh.RUnlock()

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

		returnobj = nil
		UserCache = nil

		temp := &modules.Users{}
		temp.UserName = "Bryce"

		for _, item := range modules.UserSearch {
			data := &modules.Respond{
				VideoID:      item.VideoID,
				ThumbnailURL: item.ThumbnailURL,
				VideoTitle:   item.VideoTitle,
			}
			returnobj = append(returnobj, data)
			temp.Searches = append(temp.Searches, data)
		}

		UserCache = append(UserCache, temp)

		Refresh.Unlock()
		time.Sleep(1 * time.Second)
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

func deleteUser() {}
func updateSuer() {}*/
