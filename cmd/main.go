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

type response struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {

	//create context for the client to run on
	ctx := context.Background()

	//grab API Key from JSON file
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	//create API client
	//client := createClient(ctx, config)
	client := modules.CreateClient(ctx, config)

	service, err := youtube.New(client)

	modules.HandleError(err, "Error creating YouTube client")

	modules.SearchQuery(service)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		payload := response{ID: "1234", Name: "Bryce"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
		fmt.Fprint(w)
	})
	log.Fatal(http.ListenAndServe(":8081", nil))

}

/*func updateClient() {}
func stopClient()   {}
func createUser() {}
func deleteUser() {}
func updateSuer() {}*/
