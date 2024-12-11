package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const dir = "states/"

type Data struct {
	Items []string `json:items`
}

type Space struct {
	Name    string
	Link    string
	Snippet string
}

type SpaceList []Space

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

func createJsonTest(filename string) {

	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("json not found")
	}

	var data Data
	// var state State

	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("invalid json")
	}

	data.Items = append(data.Items, "new item")
	fmt.Println(data.Items)

	updatedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("updatedData error")
	}

	err = os.WriteFile(filename, updatedData, 0644)
	if err != nil {
		fmt.Println("json write error")
	}
}

func createStateJson(state string, spaceList SpaceList, filename string) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("json not found")
		return
	}
	defer file.Close()

	var data map[string]interface{}

	decoder := json.NewDecoder(file)
	if err != nil {
		fmt.Println("invalid json")
		return
	}

	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	data[state] = spaceList

	updatedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("updatedData error")
		return
	}

	err = os.WriteFile(filename, updatedData, 0644)
	if err != nil {
		fmt.Println("json write error")
		return
	}
}

func getMakerspaces(state string) {
	s := strings.ToLower(state)
	s = strings.Replace(s, " ", "+", -1)

	rows := [][]string{}
	var spaceList SpaceList

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

		space := Space{
			Name:    title,
			Link:    link,
			Snippet: s[0],
		}

		spaceList = append(spaceList, space)
	})

	createStateJson(state, spaceList, "state.json")

}

// Takes a state as param, returns Google search results
func googleSearch(query string, count int) (*goquery.Document, error) {

	url := "https://www.google.com/search?q=" + query + "&gl=us&hl=en&num=" + strconv.Itoa(count)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(url)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	return doc, err
}

func main() {
	fmt.Println("Go Google Search Scraper")

	for i, state := range stateList {
		fmt.Println(i, state)
		getMakerspaces(state)
		time.Sleep(10 * time.Second)
	}

}
