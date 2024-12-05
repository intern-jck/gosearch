package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Entry struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

type State struct {
	Name    string
	Entries []Entry
}

var stateList = []string{
	"Alabama",
	"Alaska",
	"Arizona",
	"Arkansas",
	"California",
	"Colorado",
	"Connecticut",
	"Delaware",
	"Florida",
	"Georgia",
	"Hawaii",
	"Idaho",
	"Illinois",
	"Indiana",
	"Iowa",
	"Kansas",
	"Kentucky",
	"Louisiana",
	"Maine",
	"Maryland",
	"Massachusetts",
	"Michigan",
	"Minnesota",
	"Mississippi",
	"Missouri",
	"Montana",
	"Nebraska",
	"Nevada",
	"New Hampshire",
	"New Jersey",
	"New Mexico",
	"New York",
	"North Carolina",
	"North Dakota",
	"Ohio",
	"Oklahoma",
	"Oregon",
	"Pennsylvania",
	"Rhode Island",
	"South Carolina",
	"South Dakota",
	"Tennessee",
	"Texas",
	"Utah",
	"Vermont",
	"Virginia",
	"Washington",
	"West Virginia",
	"Wisconsin",
	"Wyoming",
}

const dir = "states/"

// Takes a state as param, returns Google search results
func searchForMakerspace(query string) (*goquery.Document, error) {

	// q := strings.ToLower(query)
	// q = strings.Replace(q, " ", "+", -1)

	url := "https://www.google.com/search?q=makerspace" + query + "&gl=us&hl=en&num=100"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// doc, err := goquery.NewDocumentFromReader(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	return goquery.NewDocumentFromReader(res.Body)
}

func createJson(entries []Entry) {

	// Create json file to save results
	// file, err := os.Create("data.json")
	// file, err := os.OpenFile("data.json", os.O_CREATE|os.O_WRONLY, 0644)
	file, err := os.Create("data.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := map[string]interface{}{
		"stateName": entries,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		panic(err)
	}
}

func createCsv(filename string, rows [][]string) {

	// Create the CSV file
	file, err := os.Create(dir + filename + ".csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the data to the CSV file
	for _, entry := range rows {
		// fmt.Println(entry[0])
		// fmt.Println(entry[1])
		// fmt.Println(entry[2])

		err := writer.Write(entry)
		if err != nil {
			panic(err)
		}
	}
}

func getMakerspaces(state string) {
	s := strings.ToLower(state)
	s = strings.Replace(s, " ", "+", -1)

	rows := [][]string{}
	doc, err := searchForMakerspace(s)
	if err != nil {
		fmt.Println(err)
	}

	c := 0
	doc.Find("div.g").Each(func(i int, result *goquery.Selection) {
		title := result.Find("h3").First().Text()
		link, _ := result.Find("a").First().Attr("href")
		snippet := result.Find(".VwiC3b").First().Text()
		// s := strings.SplitAfter(snippet, ".")
		s := regexp.MustCompile(`[.!;]`).Split(snippet, -1)

		// fmt.Printf("Title: %s\n", title)
		// fmt.Printf("Link: %s\n", link)
		// // fmt.Printf("Snippet: %s\n", snippet)
		// fmt.Println(s[0])
		// fmt.Println()

		row := []string{title, link, s[0]}
		rows = append(rows, row)

		c++
	})

	createCsv(s, rows)
}

func main() {
	fmt.Println("Go Google Search Scraper")

	// Search state for makerspace
	// stateName := "New Jersey"

	// doc, err := searchForMakerspace(stateName)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Parse through results
	// entries := []Entry{}

	// createJson(entries)
	// getMakerspaces("New Jersey")

	for i, state := range stateList {
		fmt.Println(i, state)
		getMakerspaces(state)
	}
}
