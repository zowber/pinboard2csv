package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

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

type Label struct {
	Id   *int   `json:"id"`
	Name string `json:"name"`
}

type Link struct {
	Id        *int    `json:"id"`
	Name      string  `json:"name"`
	Url       string  `json:"url"`
	Labels    []Label `json:"labels"`
	CreatedAt int     `json:"createdat"`
}

func getArgs(args []string) (string, string) {
	var inputPath, outputFormat string
	if len(args) == 3 {
		inputPath := args[1]
		outputFormat := args[2]
		return inputPath, outputFormat
	} else if len(args) == 2 {
		inputPath := args[1]
		outputFormat := "csv"
		fmt.Println("No output format specified so defaulting to CSV")
		return inputPath, outputFormat
	} else {
		log.Fatalln("Missing path to input file")
	}
	return inputPath, outputFormat
}

func main() {

	inputPath, outputFormat := getArgs(os.Args)

	input, err := os.ReadFile(inputPath)
	if err != nil {
		log.Println("Err opening input file", err)
		return
	}

	var pbls []PinboardLink

	err = json.Unmarshal(input, &pbls)
	if err != nil {
		log.Println("Err unmarshaling JSON", err)
	}

	if outputFormat == "csv" {

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

	if outputFormat == "json" {

		output, err := os.Create("output.json")
		if err != nil {
			log.Println("Err creating output file", err)
		}
		defer output.Close()

		var links []Link
		for _, pbl := range pbls {
			var link Link
			link.Name = pbl.Description
			link.Url = pbl.Href

			var labels []Label
			var tags = strings.Split(pbl.Tags, " ")
			for _, tag := range tags {
				label := Label{Name: tag}
				labels = append(labels, label)
			}
			link.Labels = labels

			link.CreatedAt = int(pbl.Time.Unix())
			links = append(links, link)
		}

		jsonData, err := json.MarshalIndent(links, "", "  ")
		if err != nil {
			log.Println("Err marshaling json", err)
		}

		output.Write(jsonData)
		if err != nil {
			log.Println("Err writing json to file", err)
		}
	}
}
