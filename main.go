package main

/*
https://tutorialedge.net/projects/hacker-news-clone-vuejs/part-4-hitting-an-api/
https://medium.com/chris-opperwall/using-the-hacker-news-api-9904e9ab2bc1
https://medium.com/@inhereat/terminal-color-rendering-tool-library-support-8-16-colors-256-colors-by-golang-a68fb8deee86
https://github.com/manifoldco/promptui
https://libs.garden/go/terminal
https://github.com/HackerNews/API
https://www.google.com/search?q=commandline+interface+golang&oq=commandline+interface+golang&aqs=chrome..69i57j0l7.9567j0j1&sourceid=chrome&ie=UTF-8
*/
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

func getTopID10(numstories int) []int {
	var storiesID []int
	response, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	bytes, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(bytes, &storiesID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return storiesID[:numstories]
}

func convertUnixtime(unixTime string) string {
	// convert unixTime to int64

	unixTimeInt64, err := strconv.ParseInt(unixTime, 10, 64)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return time.Unix(unixTimeInt64, 0).String()

}

func getStoriesData(numstories int) []Story {
	var story Story
	var stories []Story
	storiesID := getTopID10(numstories)
	time.Sleep(2 * time.Second)

	for _, value := range storiesID {
		stringStoriID := strconv.Itoa(value)

		response, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json", stringStoriID))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer response.Body.Close()

		bytes, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(bytes, &story)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stories = append(stories, story)
		time.Sleep(1 * time.Second)

	}
	return stories
}
func displayStories(storiesSlice []Story) {

	for _, value := range storiesSlice {
		stringUnixtime := strconv.Itoa(value.Time)

		fmt.Println(value.Title)
		fmt.Println("By : " + value.By)
		fmt.Println(value.URL)
		fmt.Println("Time : " + convertUnixtime(stringUnixtime) + "\n")
	}

}

func main() {
	color.Info.Tips("Fetching Stories \n")
	numberOfStories := flag.Int("stories", 10, "how many stories do you want")
	flag.Parse()

	if *numberOfStories > 50 {
		fmt.Println("Maximum stories to fetch is 50")
		os.Exit(1)
	}
	stories := getStoriesData(*numberOfStories)
	displayStories(stories)
}
