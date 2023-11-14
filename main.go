package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"time"
)

func main() {
	args := os.Args
	inputPath := args[1]

	input, err := os.ReadFile(inputPath)
	if err != nil {
		log.Println("Err opening input file", err)
		return
	}

	type PinboardLink struct {
		Href        string    `json:"href"`
		Description string    `json:"description"`
		Extended    string    `json:"extended"`
		Meta        string    `json:"meta"`
		Hash        string    `json:"hash"`
		Time        time.Time `json:"time"`
		Shared      string    `json:"shared"`
		Toread      string    `json:"toread"`
		Tags        string    `json:"tags"`
	}

	var pbls []PinboardLink

	err = json.Unmarshal(input, &pbls)
	if err != nil {
		log.Println("Err unmarshaling JSON", err)
	}

	var records [][]string
	for _, pbl := range pbls {
		record := []string{pbl.Description, pbl.Href}
		records = append(records, record)
    }

	output, err := os.Create("output.csv")
	if err != nil {
		log.Println("Err creating output file", err)
	}
	defer output.Close()

	writer := csv.NewWriter(output)

	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			log.Println("Err writing record to file", err)
		}
	}
}
