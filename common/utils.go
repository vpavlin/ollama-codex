package common

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrettyPrint(input interface{}) {
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(data))
}
