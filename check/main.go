package check

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <sources directory>", os.Args[0])
		os.Exit(1)
	}

	err := json.NewDecoder(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read request: %s", err)
	}
	checkPutResponse := fmt.Sprintf(`{"checkImplemented":"no"}`)
	fmt.Println(string(checkPutResponse))
}