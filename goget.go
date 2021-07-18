package main

import (
	"bytes"
	"fmt"
	"io"

	// "flag"
	"encoding/json"
	"net/http"
	"os"
)

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

	// If we get to this point, no errors
	// Parse body and pretty-print
	var bodyObj interface{}
	json.Unmarshal(body, &bodyObj)
	formatBody, err := json.MarshalIndent(bodyObj, "", "  ")

	if err != nil {
		fmt.Println("An error occurred parsing body of response")
		return
	}

	strBody := string(formatBody)
	fmt.Println("Response:")
	fmt.Println(strBody)

	defer res.Body.Close()
}

func postData(url, filename string) {
	// Convert file to JSON object for post
	file, err := os.Open(filename + ".json")

	if err != nil {
		fmt.Println("An error occurred reading JSON file")
		return
	}

	defer file.Close()

	// Read file for sending
	data, _ := io.ReadAll(file)
	body := bytes.NewBuffer(data)

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

	respStr := string(respBody)
	fmt.Println(respStr)
}

func main() {
	args := os.Args[1:]

	// No args provided
	if len(args) < 1 {
		fmt.Println("Improper input; use \"goget help\" if confused or stuck")
		os.Exit(1)
	}

	if args[0] == "help" {
		fmt.Println("goget <METHOD> <URL> <OPTIONS>")
		fmt.Println("METHOD: get, post")
	// If not help, send command thru switch
	} else {
		switch args[0] {
		case "get":
			getData(args[1])
		case "post":
			postData(args[1], args[2])
			
		// Incorrect command provided
		default:
			fmt.Println("Invalid command; use \"goget help\" if confused or stuck")
		}
	}
}