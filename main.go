package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type TweeBird struct {
	Id                  int64
	Username, CreatedAt string
}

type WebHook struct {
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	Content   string `json:"content"`
}

var flags = struct {
	consumerKey    string
	consumerSecret string
	username       string
	displayName    string
	avatarUrl      string
	webhook        string
	dataFolder     string
}{}

func main() {
	flag.StringVar(&flags.consumerKey, "key", "", "Twitter Consumer Key")
	flag.StringVar(&flags.consumerSecret, "secret", "", "Twitter Consumer Secret")
	flag.StringVar(&flags.username, "username", "", "Twitter Username")
	flag.StringVar(&flags.displayName, "displayName", "", "Discord Webhook Name")
	flag.StringVar(&flags.avatarUrl, "avatarUrl", "", "Discord Webhook Avatar Profile Photo Url")
	flag.StringVar(&flags.webhook, "webhook", "", "Discord Webhook Url")
	flag.StringVar(&flags.dataFolder, "dataFolder", ".", "Data folder")

	flag.Parse()
	flagutil.SetFlagsFromEnv(flag.CommandLine, "TWITTER")

	if flags.consumerKey == "" || flags.consumerSecret == "" {
		log.Fatal("Application Access Token required")
	}

	config := &clientcredentials.Config{
		ClientID:     flags.consumerKey,
		ClientSecret: flags.consumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)

	userTimelineParams := &twitter.UserTimelineParams{ScreenName: flags.username, Count: 1}
	tweet, _, _ := client.Timelines.UserTimeline(userTimelineParams)

	var lastCachedTweet = readLastCachedTweet()
	var fetchedTime, _ = time.Parse(time.RubyDate, tweet[0].CreatedAt)
	var cachedTime, _ = time.Parse(time.RubyDate, lastCachedTweet.CreatedAt)
	var tweetUrl = "https://twitter.com/" + flags.username + "/status/" + tweet[0].IDStr

	if fetchedTime.After(cachedTime) {
		sendMessage(WebHook{
			Username:  flags.displayName,
			AvatarUrl: flags.avatarUrl,
			Content:   tweet[0].User.Name + "\n" + tweetUrl,
		})
	} else {
		fmt.Println("No new tweets found.")
	}

	saveLastTweetInfo(TweeBird{
		Id:        tweet[0].ID,
		Username:  flags.username,
		CreatedAt: tweet[0].CreatedAt,
	})
}

func saveLastTweetInfo(data TweeBird) bool {
	file, _ := json.MarshalIndent(data, "", "    ")

	return ioutil.WriteFile(flags.dataFolder+"/tweebird.json", file, 0644) == nil
}

func readLastCachedTweet() TweeBird {
	file, _ := ioutil.ReadFile(flags.dataFolder + "/tweebird.json")
	data := TweeBird{}

	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func sendMessage(hook WebHook) {
	url := flags.webhook

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(hook)

	req, err := http.NewRequest("POST", url, payloadBuf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
