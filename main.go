package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gookit/color"
)

// Story struct
type Story struct {
	By          string
	Descendants int
	Kids        []int
	Score       int
	Time        int
	Title       string
	Type        string
	URL         string
}

const baseurl = "https://hacker-news.firebaseio.com/v0"

// parse time
func convertUnixtime(unixTime string) string {
	// convert unixTime to int64
	unixTimeInt64, err := strconv.ParseInt(unixTime, 10, 64)

	if err != nil {
		color.Error.Printf("%s \n", err)

		os.Exit(1)
	}
	tm := time.Unix(unixTimeInt64, 0)
	timeNow := time.Now()
	hnhour := time.Time.Hour(tm)
	nowHour := time.Time.Hour(timeNow)
	// the minus part is because of diffrence in time zones
	hoursPastInterger := (nowHour - hnhour) - 2
	hoursPast := strconv.Itoa(hoursPastInterger)

	return hoursPast

}

/*
 TODO:
- storeids in redis and compare each time if the ids has changed to avoid
fetching from the server each time
- view stories as pages and not 10 stories each time

*/

// All the id
func getID(numstories int) []int {
	var storiesID []int
	response, err := http.Get(fmt.Sprintf("%s/topstories.json", baseurl))
	if err != nil {
		color.Error.Printf("%s \n", err)

		os.Exit(1)
	}
	defer response.Body.Close()

	bytes, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(bytes, &storiesID)

	if err != nil {
		color.Error.Printf("%s \n", err)

		os.Exit(1)
	}
	return storiesID[numstories : numstories+10]
}

func getStoriesData(numstories int) []Story {
	var story Story
	var stories []Story
	storiesID := getID(numstories)
	time.Sleep(2 * time.Second)

	for _, value := range storiesID {
		stringStoriID := strconv.Itoa(value)

		response, err := http.Get(fmt.Sprintf("%s/item/%s.json", baseurl, stringStoriID))

		if err != nil {
			color.Error.Printf("%s \n", err)
			os.Exit(1)
		}
		defer response.Body.Close()

		bytes, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(bytes, &story)

		if err != nil {
			color.Error.Printf("%s \n", err)
			os.Exit(1)
		}
		stories = append(stories, story)
		time.Sleep(1 * time.Second)

	}
	return stories
}
func displayStories(storiesSlice []Story) {

	for index, value := range storiesSlice {
		indexString := strconv.Itoa(index + 1)

		color.Error.Println(fmt.Sprintf("%s %s", value.Type, indexString))

		stringUnixtime := strconv.Itoa(value.Time)

		fmt.Println("\n-------------------------------------")
		fmt.Println(value.Title + "|")
		fmt.Println("By : " + value.By + "|")
		fmt.Println(value.URL + "|")
		fmt.Println("Written " + convertUnixtime(stringUnixtime) + " hours ago |")
		fmt.Println("------------------------------------- ")
	}

}

func main() {
	// change to where to start
	numberOfStories := flag.Int("s", 0, "Show stories from")
	flag.Parse()

	color.Info.Tips("Fetching Stories \n")
	stories := getStoriesData(*numberOfStories)
	displayStories(stories)
}
