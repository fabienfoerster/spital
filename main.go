package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Box represent a chirurgical box
type Box struct {
	RegistrationNumber string `db:"registration_number" json:"registration_number"`
	Name               string `db:"name" json:"name"`
	Information        string `db:"information" json:"information"`
	Specialty          string `db:"specialty" json:"specialty"`
}

// Instrument present a chirurgical instrument
type Instrument struct {
	Ref  string `db:"ref" json:"ref"`
	Name string `db:"name" json:"name"`
}

// BoxComposition represent the composition of a chirurgical box
type BoxComposition struct {
	InstrumentID string `db:"instrumentid" json:"instrumentid"`
	Quantity     int    `db:"quantity" json:"quantity"`
}

func initDBBoxes(url, filename string) {
	url = url + "/api/v1/boxes"
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		log.Fatalf("Unable to open file: %s\n", filename)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		log.Fatalf("Unable to read file: %s\n", filename)
	}

	var aBox Box
	var counter int
	for _, each := range rawCSVData {
		aBox.RegistrationNumber = each[0]
		aBox.Name = each[1]
		aBox.Information = each[2]
		aBox.Specialty = each[4]
		fmt.Println(aBox)
		if sendData(url, aBox) == true {
			counter++
		}
	}
	log.Printf("%d were succesfully added", counter)

}

func initDBTools(url, filename string) {
	url = url + "/api/v1/instruments"

	csvFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		log.Fatalf("Unable to open file: %s\n", filename)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		log.Fatalf("Unable to read file: %s\n", filename)
	}

	var aTool Instrument
	var counter int
	for _, each := range rawCSVData {
		aTool.Ref = each[0]
		aTool.Name = each[1]
		fmt.Println(aTool)
		if sendData(url, aTool) == true {
			counter++
		}
	}
	log.Printf("%d were succesfully added", counter)
}

func initDBComposition(url, filename string) {

	csvFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		log.Fatalf("Unable to open file: %s\n", filename)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		log.Fatalf("Unable to read file: %s\n", filename)
	}

	var aBoxComposition BoxComposition
	var counter int
	for _, each := range rawCSVData {
		newURL := fmt.Sprintf("%s/api/v1/boxes/%s/content", url, each[0])
		aBoxComposition.InstrumentID = each[1]
		quantity, _ := strconv.ParseInt(each[2], 10, 0)
		aBoxComposition.Quantity = int(quantity)
		fmt.Println(aBoxComposition)
		if sendData(newURL, aBoxComposition) == true {
			counter++
		}
	}
	log.Printf("%d were succesfully added", counter)
}

func sendData(url string, value interface{}) bool {
	json, err := json.Marshal(value)
	if err != nil {
		log.Println("Marshalling failed")
		log.Fatalln(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))
	defer resp.Body.Close()
	if err != nil {
		log.Println("Something went wrong")
		log.Fatalln(err)
	}
	if resp.StatusCode == 201 {
		log.Println("Succesfully inserted into server")
		return true
	}
	log.Println("Failed to insert into server")
	return false

}

func main() {
	args := os.Args[1:]
	var url string
	flag.StringVar(&url, "url", "http://localhost:5000", "url ressource")
	flag.Parse()
	initDBBoxes(url, args[0])
	initDBTools(url, args[1])
	initDBComposition(url, args[2])

}
