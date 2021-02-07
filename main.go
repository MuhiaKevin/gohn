package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

func getTopID10(numstories int) []int {
	var storiesID []int
	response, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")

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

	return storiesID[:numstories]
}

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

func getStoriesData(numstories int) []Story {
	var story Story
	var stories []Story
	storiesID := getTopID10(numstories)
	// time.Sleep(2 * time.Second)

	for _, value := range storiesID {
		stringStoriID := strconv.Itoa(value)

		response, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json", stringStoriID))

		if err != nil {
			// fmt.Println(err)
			color.Error.Printf("%s \n", err)
			os.Exit(1)
		}
		defer response.Body.Close()

		bytes, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(bytes, &story)

		if err != nil {
			color.Error.Printf("%s \n", err)
			// fmt.Println(err)
			os.Exit(1)
		}
		stories = append(stories, story)
		// time.Sleep(1 * time.Second)

	}
	return stories
}

func displayStories(storiesSlice []Story) {
	var output string
	var color = "\033[32m"

	for index, value := range storiesSlice {
		indexString := strconv.Itoa(index + 1)
		stringUnixtime := strconv.Itoa(value.Time)
		currentTime := convertUnixtime(stringUnixtime)

		output += fmt.Sprintf(`
%s %s %s
-------------------------------------
%s |
By %s |
%s |
Written %s hours ago |
		 
		 `, color, value.Type, indexString, value.Title, value.By, value.URL, currentTime)
	}

	cmd := exec.Command("/usr/bin/less")
	cmd.Stdin = strings.NewReader(output)
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	numberOfStories := flag.Int("stories", 10, "how many stories do you want read")
	flag.Parse()

	if *numberOfStories > 50 {
		color.Error.Println("Maximum stories to fetch is 50")
		os.Exit(1)
	}
	color.Info.Tips("Fetching Stories \n")

	stories := getStoriesData(*numberOfStories)
	displayStories(stories)

}
