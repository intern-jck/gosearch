package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
	doc, err := googleSearch("makerspace+"+s, 100)
	if err != nil {
		fmt.Println(err)
	}

	doc.Find("div.g").Each(func(i int, result *goquery.Selection) {
		title := result.Find("h3").First().Text()
		link, _ := result.Find("a").First().Attr("href")
		snippet := result.Find(".VwiC3b").First().Text()
		s := regexp.MustCompile(`[.!;]`).Split(snippet, -1)

		row := []string{title, link, s[0]}
		rows = append(rows, row)
	})

	createCsv(s, rows)
}

// Takes a state as param, returns Google search results
func googleSearch(query string, count int) (*goquery.Document, error) {

	url := "https://www.google.com/search?q=" + query + "&gl=us&hl=en&num=" + strconv.Itoa(count)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	return goquery.NewDocumentFromReader(res.Body)
}

func main() {
	fmt.Println("Go Google Search Scraper")

	for i, state := range stateList {
		fmt.Println(i, state)
		getMakerspaces(state)
	}
}
