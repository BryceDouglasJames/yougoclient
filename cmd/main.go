package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	modules "github.com/brycedouglasjames/yougoclient"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var (
	Search string
)

func main() {

	//create context for the client to run on
	ctx := context.Background()

	//grab API Key from JSON file
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	//grab service config
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	//create API client
	client := modules.CreateClient(ctx, config)

	//run youtube service
	service, err := youtube.New(client)

	//catch if an error occurs
	modules.HandleError(err, "Error creating YouTube client")

	//instantiate API query
	modules.SearchQuery(service)

	/*************HANDLE*************/

	fs := http.FileServer(http.Dir("./site"))
	http.Handle("/", fs)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Hello world!")
		fmt.Fprint(w)

	})

	//handler for search results
	search := &modules.SearchRequest{ID: "nil"}
	http.HandleFunc("/query", search.SearchHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//TODO user integration
/*func updateClient() {}
func stopClient()   {}
func createUser() {}
func deleteUser() {}
func updateSuer() {}*/
