package yougoclient

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

var (
	//users      = map[int]*User{}
	//CurrentVideos = []string{}
	seq        = 1
	query      = flag.String("query", "pewdiepie", "Search term")
	maxResults = flag.Int64("max-results", 5, "Max YouTube results")
)

func CreateClient(ctx context.Context, config *oauth2.Config) *http.Client {
	TokenFile, err := TokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}

	token, err := GetTokenFromFile(TokenFile)

	if err != nil {
		token = GetTokenFromWeb(config)
		SaveToken(TokenFile, token)
	}
	return config.Client(ctx, token)
}

func TokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

func GetTokenFromFile(filename string) (*oauth2.Token, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	retrievedToken, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return retrievedToken
}

func SaveToken(filename string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", filename)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

/*YOUTUBE QUERY FUNCTION*/
func SearchQuery(service *youtube.Service) map[string]string {
	call := service.Search.List([]string{"id,snippet"}).Q(*query).MaxResults(*maxResults)
	response, err := call.Do()
	HandleError(err, "")

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
			//fmt.Printf("%s", item.Snippet.Title)
		}
	}

	return videos
}

func RelatedVideoGenerate(service *youtube.Service, videoPass map[string]string) *Users {
	user := &Users{}
	for key := range videoPass {
		call2 := service.Search.List([]string{"id,snippet"}).RelatedToVideoId(key).Type("video").MaxResults(*maxResults)
		response, err := call2.Do()
		HandleError(err, "")
		for _, item := range response.Items {
			data := &Respond{}
			data.SetResponse(item.Id.VideoId, item.Snippet.Thumbnails.Default.Url, item.Snippet.Title)
			user.AddVideo(data)
		}
	}
	return user
}

func HandleError(err error, message string) string {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}

	return string(message)
}
