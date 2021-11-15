package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "country.csv", "file containing list of countries for covis stats")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("failed to open the csv file %s\n", *csvFilename))
	}
	//r := csv.NewReader(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var list []string

	for scanner.Scan() {
		list = append(list, strings.ToUpper(scanner.Text()))
	}
	for i, country := range list {
		fmt.Printf("%d. %s\n", i+1, country)
	}
	fmt.Println("Enter country name:")
	var val string
	fmt.Scanf("%s", &val)
	val = strings.ToUpper(val)
	countryFound := find(list, val)
	if countryFound {
		fmt.Printf("Country %s found. Data fetching please wait...\n", val)
		APIcall(val)
	} else {
		fmt.Printf("Country %s not found. Please enter from list above", val)

	}

	// countries, err := r.ReadAll()
	// if err != nil {
	// 	exit("Failed to open parse the provided csv file")
	// }
	// for i, name := range countries {
	// 	fmt.Printf("%d. %s\n", i+1, name)
	// }
	// fmt.Println("Enter country name you wish to get covid statistics of: ")
	// //var val string
	// val := ""
	// fmt.Scanf("%s\n", &val)
	// find(countries, val)
	// if val != "" {
	// 	fmt.Println("found the country!! Fetching for data please wait....")
	// }

}

func find(countries []string, val string) bool {
	for i, _ := range countries {
		if countries[i] == val {
			return true
		}
	}
	return false
}

func APIcall(val string) {
	url := "https://covid-19-data.p.rapidapi.com/country"

	req, _ := http.NewRequest("GET", url, nil)

	// req.Header.Add("content-type", "application/json")
	q := req.URL.Query()
	q.Add("name", val)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("x-rapidapi-host", "covid-19-data.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "89611c7851msh195755ef7dc1161p1ddc0ejsnef8e72d70806")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	//fmt.Println(string(body))
	var responseObject Response
	json.Unmarshal(body, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)
	// a := fs.FileMode(os.ModeAppend)
	// info, err := os.Stat("covidstats.json")
	// _ = info
	// if os.IsNotExist(err) {
	// 	//fmt.Println("file doesnot exist")
	// 	os.Create("covidstats.txt")
	// } else {
	// 	file, err := os.OpenFile("covidstats.json", os.O_APPEND|os.O_WRONLY, 0644)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer file.Close()
	// 	file.WriteString(string(body))
	// }
	ioutil.WriteFile("covidstats.json", body, 0644)

	// file, err := os.OpenFile("covidstats.txt", os.O_APPEND|os.O_WRONLY, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// file.WriteString(string(body))

}

type Response []struct {
	Country        string    `json:"country"`
	Code           string    `json:"code"`
	ConfirmedCases int       `json:"confirmed"`
	RecoveredCases int       `json:"recovered"`
	CriticalCases  int       `json:"critical"`
	Deaths         int       `json:"deaths"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	LastChange     time.Time `json:"lastChange"`
	LastUpdated    time.Time `json:"lastUpdate"`
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)

}
