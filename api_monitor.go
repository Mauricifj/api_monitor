package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	welcome()

	for {
		showMenu()

		option := readOption()

		switch option {
		case 1:
			monitor(SANDBOX)

		case 2:
			monitor(PRODUCTION)

		case 0:
			fmt.Println("See you later, my friend")
			os.Exit(0)

		default:
			fmt.Println("Invalid option")
			os.Exit(-1)
		}
	}
}

// Display and menu
func welcome() {
	version := "1.0"
	fmt.Println("*** API MONITOR - VERSION", version, "***")
	newLine()
}

func showMenu() {
	fmt.Println("--- MENU ---")
	fmt.Println("1 - Start monitoring SANDBOX APIs")
	fmt.Println("2 - Start monitoring PRODUCTION APIs")
	fmt.Println("0 - Quit")
	newLine()
}

func readOption() int {
	fmt.Print("Choose one option, please: ")
	var option int
	fmt.Scan(&option)
	return option
}

func newLine() {
	fmt.Println()
}

// Monitoring
func monitor(environment string) {
	newLine()
	fmt.Println("-- API MONITOR STARTED --")
	newLine()

	urls, err := extractUrlsFromTxt(environment)

	if err != nil {
		fmt.Println("Error on", err)
		os.Exit(-1)
	}

	for _, url := range urls {
		healthiness, err := checkHealthiness(url)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("API:", url, " - IsHealthy:", healthiness)
		}
	}

	newLine()
	fmt.Println("-- API MONITOR ENDED --")
	newLine()
}

func checkHealthiness(url string) (bool, error) {
	response, httpError := http.Get(url)

	if httpError != nil {
		return false, httpError
	}

	apiResponse, convertError := convertToObject(response.Body)

	if convertError != nil {
		return false, convertError
	}

	return apiResponse.IsHealthy, nil
}

// Files
func extractUrlsFromTxt(environment string) ([]string, error) {
	var (
		fileName string
		urls []string
	)

	if environment == PRODUCTION {
		fileName = PRODUCTION
	} else {
		fileName = SANDBOX
	}

	file, err := os.Open(fileName + ".txt")

	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	for {
		url, err := reader.ReadString('\n')
		url = strings.TrimSpace(url)
		urls = append(urls, url)

		if err == io.EOF {
			break
		}
	}

	file.Close()

	return urls, nil
}

// Response converter
func convertToObject(response io.ReadCloser) (apiResponse, error) {
	var convertedObject apiResponse

	responseData, err := ioutil.ReadAll(response)

	if err != nil {
		return apiResponse{}, err
	}

	json.Unmarshal(responseData, &convertedObject)

	return convertedObject, nil
}

type apiResponse struct {
	IsHealthy bool
}

// Environments
const(
	PRODUCTION string = "production"
	SANDBOX    string = "sandbox"
)
