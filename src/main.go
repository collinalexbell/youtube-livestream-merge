// Sample Go code for user authorization

package main

import (
	"encoding/json"
	"fmt"
	"time"
	"errors"
	"log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message + ": %v", err.Error())
	}
}

func channelsListByUsername(service *youtube.Service, part string, forUsername string) {
	call := service.Channels.List([]string{part})
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', " +
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}

func listLivestreams(service *youtube.Service, forChannel string) []string {
	call := service.Search.List([]string{"snippet"})
	call = call.ChannelId(forChannel)
	call = call.EventType("live")
	call = call.Type("video")
	response, err := call.Do()
	if err != nil {
		log.Fatalf("err")
	}
	rv := []string{}
	for _, item := range(response.Items){
		rv = append(rv, item.Id.VideoId)
	}

	return rv
}

func getSecondInTime(service *youtube.Service, livestreams []string) (*youtube.Video, error){
	call := service.Videos.List([]string{"liveStreamingDetails, snippet"})
	call = call.Id(livestreams...)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("err")
	}

	latest := ""
	var latestStream *youtube.Video

	for _, item := range(response.Items) {
		if latest == "" || latest < item.LiveStreamingDetails.ActualStartTime {
			latest = item.LiveStreamingDetails.ActualStartTime
			latestStream = item
		}
	}
	if latest != "" && len(response.Items) >= 2{
		return latestStream, nil
	} else {
		return nil, errors.New("error")
	}
}


func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	time.Sleep(5 * time.Minute)
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")
	livestreams := listLivestreams(service, "UChsxOQf3j6Jw_BbzWcsIQPg")
	latestStream, secondErr := getSecondInTime(service, livestreams)
	if secondErr != nil {
		fmt.Println("only 1 stream")
	} else {
		url := "https://www.youtube.com/embed/" + latestStream.Id + "?autoplay=1"
		err = exec.Command("open", url).Start()
	}

}
