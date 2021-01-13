package yougoclient

import (
	"flag"
	"fmt"
	"log"

	"google.golang.org/api/youtube/v3"
)

var (
	seq        = 1
	query      = flag.String("query", "java", "Search term")
	maxResults = flag.Int64("max-results", 1, "Max YouTube results")
)

/*
*YOUTUBE QUERY FUNCTION
 */
func SearchQuery(service *youtube.Service, search string) map[string]string {
	call := service.Search.List([]string{"id,snippet"}).Q(CurrentSearch).MaxResults(int64(5))
	response, err := call.Do()
	HandleError(err, "")

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
			fmt.Printf("%s", item.Snippet.Title)
			fmt.Println(item.Snippet.Title)
		}
	}

	return videos
}

/*Function loops through queried videos and makes
 *another request for more related videos.
 */
func RelatedVideoGenerate(service *youtube.Service, videoPass map[string]string) *Users {
	user := &Users{}
	for key := range videoPass {
		call2 := service.Search.List([]string{"id, snippet"}).RelatedToVideoId(key).Type("video").MaxResults(int64(1))
		response, err := call2.Do()
		HandleError(err, "")
		for _, item := range response.Items {

			//time.Sleep(1 * time.Second)

			if item.Snippet != nil {
				fmt.Println("+++" + item.Id.VideoId + "+++")
				fmt.Println(item.Snippet.Title)

				data := &Respond{
					VideoID:      item.Id.VideoId,
					ThumbnailURL: item.Id.VideoId,
					VideoTitle:   item.Snippet.Title,
				}
				UserSearch = append(UserSearch, data)
			}
		}
	}
	return user
}

/*
*ERROR HANDLER
 */
func HandleError(err error, message string) string {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}

	return string(message)
}
