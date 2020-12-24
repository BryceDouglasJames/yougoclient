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
	CurrentUsers []*modules.Users
	Refresh      sync.RWMutex
	UserCache    []*modules.Users
	returnobj    []*modules.Respond
)

func main() {

	go RefreshSearch()

	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		Refresh.RLock()
		defer Refresh.RUnlock()

		response, err := json.Marshal(UserCache)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			fmt.Println()
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(response))

	})

	ThisClient := &modules.Users{}
	http.HandleFunc("/videos", ThisClient.ServeArray)

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

//TODO user integration
/*func updateClient() {}
func stopClient()   {}
func createUser() {}
func deleteUser() {}
func updateSuer() {}*/
