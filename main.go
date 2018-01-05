package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	a "github.com/ChimeraCoder/anaconda"
)

type auth struct {
	ConsumerKey, ConsumerSecret, AccessToken, AccessTokenSecret string
}

var keys auth

// Get authentication keys from file
func init() {
	data, err := ioutil.ReadFile("auth.json")
	if err != nil {
		log.Panic("Error", err)
	}

	err = json.Unmarshal(data[:len(data)-1], &keys)
	if err != nil {
		log.Panic("Error", err)
	}
}

func main() {
	a.SetConsumerKey(keys.ConsumerKey)
	a.SetConsumerSecret(keys.ConsumerSecret)
	api := a.NewTwitterApi(keys.AccessToken, keys.AccessTokenSecret)

	options := url.Values{
		"count": {"200"},
	}

	timeline, _ := api.GetHomeTimeline(options)

	timeline = filter(timeline)
	fmt.Println(len(timeline))
	for i, tweet := range timeline {
		fmt.Printf("====== %v =======\n", i)
		fmt.Println(tweet.User.Name)
		fmt.Println(tweet.Text)
		fmt.Println(tweet.Favorited)
		// if i%6 == 0 {
		//  // Like the tweet
		// 	fav, _ := api.Favorite(tweet.Id)
		// 	fmt.Println(fav.Text)
		// }
	}
}

func filter(t []a.Tweet) []a.Tweet {
	// Words, tags etc. to follow
	tags := []string{"#100DaysOfCode"}
	// Users not to follow
	xUser := []string{"jetrubyagency"}

	filtered := []a.Tweet{}
	for _, v := range t {
		if hasString(tags, v.Text) && !hasString(xUser, v.User.Name) && len(strings.Split(v.Text, "#")) <= 5 && !v.Favorited {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func hasString(ss []string, s string) bool {
	for _, v := range ss {
		if strings.Index(strings.ToLower(s), strings.ToLower(v)) != -1 {
			return true
		}
	}
	return false
}
