package yougoclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func NewClient() {
	//get call context
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
	client := CreateClient(ctx, config)

	//run youtube service
	service, err := youtube.New(client)
	HandleError(err, "Error creating YouTube client")

	//API query
	videos := SearchQuery(service, CurrentSearch)

	//package response
	RelatedVideoGenerate(service, videos)

	//FOR TESTING
	/*data := &Respond{
		VideoID:      "PLZ",
		ThumbnailURL: "HELP",
		VideoTitle:   "HELLO",
	}
	UserSearch = append(UserSearch, data)
	//fmt.Println(UserSearch)*/
}

/*
*Creates user struct and appends to session clients slice
 */
func NewUser(name string) {
	tempUser := &Users{
		UserName:    name,
		SessionTime: 0,
	}
	ClientList = append(ClientList, tempUser)
}

/*
*Searches linearly through session clients and deletes selected username
 */
func DeleteUser(name string) {
	fmt.Println("User " + name + " has timed out of their session.")
	tempList := ClientList
	ClientList = nil
	for _, user := range tempList {
		if user.UserName != name {
			ClientList = append(ClientList, user)
		}
	}
}
