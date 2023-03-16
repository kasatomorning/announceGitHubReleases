package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const FILENAME = "memo.txt"

func main() {
	// use -i option to initialize file
	// if not, existing file is used
	var (
		i = flag.Bool("i", false, "initialize file")
	)
	flag.Parse()
	isInit := *i
	if isInit {
		fmt.Println("initializing file...")
		//initializing file
		releases, getErr := getReleasesData()
		if getErr != nil {
			log.Fatal("Error: %v\n", getErr)
		}
		initializingErr := initializeFile(releases)
		if initializingErr != nil {
			log.Fatal("Error: %v\n", initializingErr)
		}
		return
	}
	//use existing file
	fmt.Println("loading file...")
	postreleases, loadingerr := loadFile(FILENAME)
	if loadingerr != nil {
		log.Fatal("Error: %v\n", loadingerr)
	}
	nowreleases, getErr := getReleasesData()
	if getErr != nil {
		log.Fatal("Error: %v\n", getErr)
	}
	newData, compareErr := compareAndWriteToFile(postreleases, nowreleases, FILENAME)
	if compareErr != nil {
		log.Fatal("Error: %v\n", compareErr)
	}
	if len(newData) != 0 {
		announceDataToXXX(newData)
	}
}

func initializeFile(array []string) error {
	file, err := os.Create(FILENAME)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range array {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// load file and return the array of string per line and error
func loadFile(fn string) ([]string, error) {
	file, err := os.Open(fn)
	if err != nil {
		return []string{}, err
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
	fmt.Println("Loaded " + fn + " successfully")
	return returnData, nil
}

func compareAndWriteToFile(postData []string, nowData []string, fn string) ([]string, error) {
	if len(nowData) == 0 {
		return []string{}, fmt.Errorf("get empty release array")
	}
	if len(postData) == 0 {
		return []string{}, fmt.Errorf("get empty text, please initialize file")
	}
	addedData := make([]string, 0)
	for _, now := range nowData {
		for j, post := range postData {
			if now == post {
				break
			}
			if j == len(postData)-1 {
				addedData = append(addedData, now)
				fmt.Println("Added release: " + now)
			}
		}
	}
	if (len(addedData)) == 0 {
		fmt.Println("No new release")
		return []string{}, nil
	}
	//OpenFile method can overwrite file
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()
	for _, line := range addedData {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return []string{}, err
		}
	}
	return addedData, nil
}

func getReleasesData() ([]string, error) {
	//load .env files
	err := godotenv.Load(".env")
	if err != nil {
		return []string{}, err
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
		//default is 30 -> 100
		ListOptions := github.ListOptions{PerPage: 100, Page: pagecount}
		releases, response, err := client.Repositories.ListReleases(ctx, "yt-dlp", "yt-dlp", &ListOptions)
		//error handling
		if err != nil {
			return []string{}, err
		}
		//check response statuscode
		if response.StatusCode != 200 {
			return []string{}, fmt.Errorf("response status_code is not correct")
		}
		//check if there are no more releases
		if len(releases) == 0 {
			break
		}
		//check rate limit
		if response.Rate.Remaining <= 0 {
			return []string{}, fmt.Errorf("rate limit exceeded")
		}
		for i := 0; i < len(releases); i++ {
			releasesList = append(releasesList, *releases[i].TagName)
		}
	}
	fmt.Printf("Get %d releases\n", len(releasesList))
	return releasesList, nil
}

func announceDataToXXX(data []string) {
	//to be implemented
	fmt.Println("なんかに通知する関数")
}
