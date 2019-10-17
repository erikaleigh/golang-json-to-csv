package main

import (
	"fmt"
	"os"
	"sync"

	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Response struct {
	Items []Item
}
type Item struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Contact   Contact `json:"contact"`
	Pet       string  `json:"pet"`
	Car       string  `json:"car"`
}

type Contact struct {
	Email    string  `json:"email"`
	Timezone string  `json:"timezone"`
	Address  Address `json:"address"`
}

type Address struct {
	StreetName  string `json:"street_name"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}


// processRow takes an entry from a .json file, creates a record with the info we
// need from it, and pushes this record to a channel
// It calls Done() on the waitGroup to unblock when finished
func processRow(item Item, ch chan []string, wg *sync.WaitGroup) {
	record := []string{
		item.ID,
		item.FirstName,
		item.LastName,
		item.Contact.Email,
		item.Contact.Timezone,
		item.Contact.Address.StreetName,
		item.Contact.Address.City,
		item.Contact.Address.Country,
		item.Contact.Address.CountryCode,
		item.Pet,
		item.Car,
	}

	ch <- record

	wg.Done()
}

// writeRow writes a record to the output file, and calls Done() on the waitGroup to unblock when finished
func writeRow(outputFile *os.File, row []string, wg *sync.WaitGroup) {
	w := csv.NewWriter(outputFile)

	if err := w.Write(row); err != nil {
		fmt.Println(err)
	}

	w.Flush()
	wg.Done()
}

// getFileNames searches the ./files/json/ folder and returns a slice containing
// the names of each of the files found that end in .json
func getFileNames() []string {
	var fileNames []string
	files, err := ioutil.ReadDir("./files/json/")
	if err != nil {
		return []string{}
	}

	for _, file := range files {
		// Check that the file is not a directory, and that it has the correct extension of .json
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			fileNames = append(fileNames, fmt.Sprintf("./files/json/%s", file.Name()))
		}
	}

	return fileNames
}

func main() {
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	var response Response
	var outputFile *os.File

	// Find the names of all of the files in the ./files/json folder with an extension of .json
	files := getFileNames()

	// Create the output .csv file that all of the json entries will be added to
	// This file will be returned in the folder ./files/csv/
	outputFile, err := os.Create("./files/csv/output.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer outputFile.Close()

	// For each of the .json files found in the ./files/json folder,
	// read the file and process each entry, adding a new entry
	// with the info to the output .csv file
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		}

		// unmarshal the data to be converted from the .json file
		err = json.Unmarshal(data, &response)
		if err != nil {
			fmt.Println(err)
		}

		// Create a channel to push each of the .json entries into
		ch := make(chan []string)
		defer close(ch)

		items := response.Items
		// For each entry in the .json data, process the data we need,
		// and push it into the channel
		// Use a waitGroup to block while the goroutine processes
		for _, item := range items {
			wg1.Add(1)
			go processRow(item, ch, &wg1)
		}

		go func() {
			// For each entry in the channel, write the row into the .csv file
			// Use a waitGroup to block while the goroutine processes
			for row := range ch {
				wg2.Add(1)
				go writeRow(outputFile, row, &wg2)
			}
		}()

		wg1.Wait()
		wg2.Wait()
	}
}