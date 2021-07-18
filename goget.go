package main

import (
	"fmt"
	"io"
	// "flag"
	"net/http"
	"os"
)

func getData(url string) {
	res, err := http.Get(url)

	if err != nil {
		fmt.Println("An error ocurred fetching request.")
		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("An error ocurred reading body of response.")
	}

	strBody := string(body)
	fmt.Println(strBody)

	defer res.Body.Close()
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Improper input.")
		os.Exit(1)
	}

	if args[0] == "help" {
		fmt.Println("goget <METHOD> <URL> <OPTIONS>")
		fmt.Println("METHOD: get, post")
	} else {
		switch args[0] {
		case "get":
			url := args[1]
			getData(url)
		}
	}
}