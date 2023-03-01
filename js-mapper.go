package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a URL or file as an argument")
		return
	}

	input := os.Args[1]

	// Check if the argument is a file
	if _, err := os.Stat(input); err == nil {
		// Open the file
		f, err := os.Open(input)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		// Read the file line by line
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			// Set the URL variable as the scanned line
			url := scanner.Text()

			// Send the URL to the grabber function
			grabber(url)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	} else {
		// Set the URL variable as the input argument
		url := input

		// Send the URL to the grabber function
		grabber(url)
	}
}

func grabber(url2 string) {
	var rnd strings.Builder
	for i := 0; i < 8; i++ {
		rnd.WriteString(strconv.Itoa(rand.Intn(10)))
	}
	normalString := rnd.String()
	queryString := "?cache=" + normalString + ""
	fullUrl := strings.Join([]string{url2, queryString}, "")
	u, err := url.Parse(fullUrl)
	if err != nil {
		fmt.Println(err)

	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	bodyStr := string(bodyBytes)

	if strings.Contains(bodyStr, "sourceMappingURL") {
		u, err := url.Parse(url2)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Found Mapping URL - " + url2 + "")
		u.RawQuery = ""
		u.Fragment = ""
		newURL := u.String() + ".map"
		fmt.Println("Mapping URL - " + newURL + "")
		req2, err := http.NewRequest("GET", newURL, nil)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		resp2, err := client.Do(req2)
		if err != nil {
			fmt.Println(err)
			return
		}
		bodyBytes2, err := ioutil.ReadAll(resp2.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		bodyStr2 := string(bodyBytes2)
		fmt.Println("body2:", bodyStr2)
		if strings.Contains(bodyStr2, "mappings") {
			f, err := os.OpenFile("mapping.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer f.Close()

			_, err = f.WriteString("" + newURL + "\n")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println("Mapping URL written to mapping.txt")
		}
	} else {
		fmt.Println("No SourceMap Found - " + url2 + "")

	}

}
