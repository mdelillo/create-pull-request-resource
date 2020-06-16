package main

import (
	"encoding/json"
	"fmt"
	. "github.com/pivotal/create-pull-request-resource"
	"log"
	"os"
)

func main() {

	var request InRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Fatalf("failed to read request: %s", err.Error())
	}

	checkPutResponse := fmt.Sprintf(`[{ "ref": "%s" }]`, request.Version.Ref)

	fmt.Println(string(checkPutResponse))
}
