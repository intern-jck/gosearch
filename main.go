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

	"github.com/PuerkitoBio/goquery"
)

const dir = "states/"

type Data struct {
	Items []string `json:items`
}

type Space struct {
	Name string
	Link string
}

type StateList []Space

type SpaceList map[string]StateList

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

func createStateJson(filename string) {

	var stateList StateList

	// spaceList := make(map[string]StateList)

	space1 := Space{
		Name: "space 1",
		Link: "www.space_1.com",
	}

	space2 := Space{
		Name: "space 2",
		Link: "www.space_2.com",
	}

	space3 := Space{
		Name: "space 3",
		Link: "www.space_3.com",
	}

	stateList = append(stateList, space1)
	stateList = append(stateList, space2)
	stateList = append(stateList, space3)

	// fmt.Println(stateList)

	// spaceList["state_1"] = stateList
	// spaceList["state_2"] = stateList
	// spaceList["state_3"] = stateList
	// fmt.Println(spaceList)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("json not found")
	}
	defer file.Close()

	var data map[string]interface{}

	decoder := json.NewDecoder(file)
	if err != nil {
		fmt.Println("invalid json")
	}

	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Print the map
	fmt.Println(data)

	// data.Items = append(data.Items, "new item")
	// fmt.Println(data.Items)

	data["state_1"] = stateList

	fmt.Println(data)

	updatedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("updatedData error")
	}

	err = os.WriteFile(filename, updatedData, 0644)

	if err != nil {
		fmt.Println("json write error")
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

	// for i, state := range stateList {
	// 	fmt.Println(i, state)
	// 	getMakerspaces(state)
	// }

	// createJsonTest("state.json")
	createStateJson("state.json")
}
