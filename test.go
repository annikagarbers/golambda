package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	// If there is an arg on the command line, use it; that's where
	// the JSON is passed by the JavaScript launcher. Otherwise,
	// read from stdin (for local testing).

	var snsJson []byte

	if len(os.Args) > 1 {
		snsJson = []byte(os.Args[1])
	} else {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("error reading json... %s`", err)
			os.Exit(1)
		}
		snsJson = bytes
	}

	// Okay, we have the raw bytes from the Lambda post; parse
	// that *outer* JSON data so we can get to the inner JSON
	// from GitHub.
	var lambdaMsg map[string]interface{}
	jsonErr := json.Unmarshal(snsJson, &lambdaMsg)
	if jsonErr != nil {
		fmt.Printf("error parsing lambda message: %s", jsonErr)
		os.Exit(1)
	}

	// Got the JSON in an interface. Pull out the first "record"
	// and let's get the GitHub JSON out of it... unfortunately
	// it's pretty far down in there.
	records := lambdaMsg["Records"].([]interface{})
	r := records[0].(map[string]interface{})

	sns := r["Sns"].(map[string]interface{})
	message := sns["Message"].(string)

	// For testing, re-enable these lines to just dump the raw JSON
	// or whatever's in the "Message" element.
	// fmt.Println(message)
	// os.Exit(0)

	var githubMsg map[string]interface{}
	jsonErr = json.Unmarshal([]byte(message), &githubMsg)
	if jsonErr != nil {
		fmt.Printf("Message does not appear to be JSON: %s", message)
		os.Exit(1)
	}
	fmt.Printf("parsed github message is  %s\n", githubMsg)

}
