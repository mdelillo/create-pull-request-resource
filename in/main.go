package in

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <sources directory>", os.Args[0])
		os.Exit(1)
	}

	json.NewDecoder(os.Stdin)

	inPutResponse := fmt.Sprintf(`{"inImplemented":"no"}`)
	fmt.Println(string(inPutResponse))
}