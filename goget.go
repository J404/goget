package main

import (
	"bytes"
	"fmt"
	"io"
	"flag"
	"encoding/json"
	"net/http"
	"os"
)

// Returns a pretty-print string of the supplied JSON byte arr
func jsonToPretty(data []byte) string {
	var formatted interface{}
	json.Unmarshal(data, &formatted)
	jsonStr, err := json.MarshalIndent(formatted, "", "  ")

	if err != nil {
		fmt.Println("An error occurred parsing body of response")
		os.Exit(1)
	}

	return string(jsonStr)
}

func getData(url string) {
	res, err := http.Get(url)

	// Error fetching response
	if err != nil {
		fmt.Println("An error occurred fetching request.")
		return
	}

	body, err := io.ReadAll(res.Body)

	// Error reading body
	if err != nil {
		fmt.Println("An error occurred reading body of response.")
		return
	}

	fmt.Println("Response:")
	fmt.Println(jsonToPretty(body))

	defer res.Body.Close()
}

func postData(url, filename string, manualData bool) {
	var body *bytes.Buffer


	// If no manual data, get data from json file
	if !manualData {
		// Convert file to JSON object for post
		file, err := os.Open(filename + ".json")

		if err != nil {
			fmt.Println("An error occurred reading JSON file")
			return
		}

		defer file.Close()

		// Read file for sending
		data, _ := io.ReadAll(file)
		body = bytes.NewBuffer(data)
	} else {
		
		// Otherwise, "filename" should be treated as raw JSON data
		byteArr := []byte(filename)
		body = bytes.NewBuffer(byteArr)
	}
	

	res, err := http.Post(url, "application/json", body)

	if err != nil {
		fmt.Println("Error POSTing to URL")
		return
	}

	fmt.Printf("Response (status %d):\n", res.StatusCode)

	respBody, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		fmt.Println("Error reading response")
		return
	}

	fmt.Println("Response:")
	fmt.Println(jsonToPretty(respBody))
}

func printHelp() {
	fmt.Println("goget <FLAGS> <METHOD> <URL> <OPTIONS>")

	fmt.Println("FLAGS")
	fmt.Println("-manual: enter JSON object manually ({\"KEY\":\"VALUE\"})")
	fmt.Println("METHOD: get, post")
	fmt.Println("URL: url to query")

	fmt.Println("OPTIONS:")
	fmt.Println("If POSTing, give filename of JSON data file (no extension)")
	fmt.Println("If POSTING, flag -data and enter raw JSON data (NEED TO ESCAPE QUOTES W/IN JSON)g")
}

func main() {
	manualData := flag.Bool("manual", false, "Manually JSON data for POST request")

	flag.Parse()
	args := os.Args[1:]
	
	if args[0] == "-manual" {
		args = os.Args[2:]
	}

	// No args provided
	if len(args) < 1 {
		fmt.Println("Improper input; use \"goget help\" if confused or stuck")
		os.Exit(1)
	}

	if args[0] == "help" {
		printHelp()

	// If not help, send command thru switch
	} else {
		switch args[0] {
		case "get":
			getData(args[1])
		case "post":
			if args[2] != "" {
				postData(args[1], args[2], *manualData)
			}

		// Incorrect command provided
		default:
			fmt.Println("Invalid command; use \"goget help\" if confused or stuck")
		}
	}
}