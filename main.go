package main

import (
	"encoding/json"
	"fmt"
	"github.com/simonfalke-01/cbr-cli/kooky"
	_ "github.com/simonfalke-01/cbr-cli/kooky/browser/all"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	checkArgs(3)

	cookie := kooky.ReadCookies(kooky.Valid, kooky.DomainHasSuffix(`codebreaker.xyz`), kooky.Name(`google-login-session`))[0].Value
	problemName := os.Args[1]
	urlStr := fmt.Sprintf("https://codebreaker.xyz/problem/%s", problemName)

	resp, err := http.Get(fmt.Sprintf("https://cbr-api.simonfalke.studio/api/getSubmissions?problemId=%v", problemName))
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var attemptsList []Attempt
	err = json.NewDecoder(resp.Body).Decode(&attemptsList)
	if err != nil {
		panic(err)
	}

	latestId := func(id int, err error) int {
		if err != nil {
			panic(err)
		}
		return id
	}(strconv.Atoi(attemptsList[0].ID)) + 1

	solutionPath := os.Args[2]
	// read solutionPath
	code, err := os.ReadFile(solutionPath)
	if err != nil {
		panic(err)
	}

	resp = submit(urlStr, string(code), cookie)

	fmt.Printf("Submitted solution for %v. ", problemName)
	if resp.StatusCode == 200 {
		fmt.Println("Success!")
		askOpenBrowser(fmt.Sprintf("https://codebreaker.xyz/submission/%v", latestId))
	} else {
		fmt.Println("Failure!")
	}
}
