package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const fileName = "memo.txt"

func main() {
	isInit := true
	if isInit {
		//initializing file
		releases := getReleasesData()
		initializingErr := initializeFile(releases)
		if initializingErr != nil {
			log.Fatal("Error: %v\n", initializingErr)
		}
		return
	} else {
		//use existing file
		postreleases, err := readFile(fileName)
		nowreleases := getReleasesData()
		writeDataToFile(releases)
	}
}

func initializeFile(array []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Error: Creating file")
		return err
	}
	defer file.Close()
	for _, line := range array {
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func readFile(fn string) ([]string, error) {
	file, err := os.Open(fn)
	if err != nil {
		log.Fatal("Error: %v\n", err)
	}
	defer file.Close()
	returnData := make([]string, 0)
	reader := bufio.NewReader(file)
	for {
		//read file per line
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []string{}, err
		}
		returnData = append(returnData, string(line))
	}
	return returnData, nil
}

func compareAndWriteToFile(postData []string, nowData []string) error {
	postData, err := readFile(fn)
	if err != nil {
		return err
	}
	for i, line := range array {
		if line != postData[i] {
			fmt.Println("New release found: ", line)
		}
	}
}

func getReleasesData() []string {
	//initialize .env files
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error: Loading .env file")
	}
	token := os.Getenv("GITHUB_TOKEN")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	releasesList := make([]string, 0)

	for pagecount := 1; ; pagecount++ {
		//Perpage default is 30
		ListOptions := github.ListOptions{PerPage: 100, Page: pagecount}
		releases, response, err := client.Repositories.ListReleases(ctx, "yt-dlp", "yt-dlp", &ListOptions)
		//error handling
		if err != nil {
			log.Fatal("Error: %v\n", err)
			return []string{}
		}
		//check response statuscode
		if response.StatusCode != 200 {
			log.Fatal("Error: %v\n", response.Status)
			return []string{}
		}
		//check if there are no more releases
		if len(releases) == 0 {
			fmt.Println("Reached the last release")
			break
		}
		//check rate limit
		if response.Rate.Remaining <= 0 {
			log.Fatal("Error: Rate limit exceeded")
			return []string{}
		}
		for i := 0; i < len(releases); i++ {
			releasesList = append(releasesList, *releases[i].TagName)
		}
	}
	for i := 0; i < len(releasesList); i++ {
		fmt.Println(releasesList[i])
	}
	return releasesList
}
