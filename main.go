package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

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

	for {
		timeline, _ := api.GetHomeTimeline(options)

		timeline = filter(timeline)

		// Add some randomness
		r := random(2, 3)
		for i, tweet := range timeline {
			if i%r == 0 {
				_, err := api.Favorite(tweet.Id)
				if err != nil {
					fmt.Println("Error", err)
				}
			}
		}
		time.Sleep(15 * time.Minute)
	}

}

func random(min, max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return r.Intn(max+1) + min
}

func filter(t []a.Tweet) []a.Tweet {
	// Words, tags etc. to follow
	tags := getList("tags.txt")
	// Words, tags etc. to filter out
	xTags := getList("tags-exclude.txt")
	// Users to filter out
	xUser := getList("exclude.txt")

	filtered := []a.Tweet{}
	for _, v := range t {
		if hasString(tags, v.Text) &&
			!hasString(xTags, v.Text) &&
			!hasString(xUser, v.User.ScreenName) &&
			len(strings.Split(v.Text, "#")) <= 5 && !v.Favorited {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func getList(fn string) []string {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return []string{}
	}

	res := strings.Split(string(data), "\n")

	var resF []string
	for _, v := range res {
		if len(v) != 0 {
			resF = append(resF, v)
		}
	}

	return resF
}

func hasString(ss []string, s string) bool {
	for _, v := range ss {
		if strings.Index(strings.ToLower(s), strings.ToLower(v)) != -1 {
			return true
		}
	}
	return false
}
