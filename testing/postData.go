package main

import (
	"fmt"

	processor "github.com/m-ahrukh/go-assignment/mainProcessor"
)

func main() {
	jsonResults := make(map[string]bool)
	xmlResults := make(map[string]bool)

	jsonResults["6ba5629bb6a2b9cf6b876aba64e6c90d"] = true //o
	jsonResults["ab896f7700ff79cee22ed7d418802d29"] = true //o
	jsonResults["3fb686990f1ee536e4ad8cdb97963364"] = true //o
	jsonResults["9db54c49a2607ea179054532cdc9ee79"] = true //j
	jsonResults["811c9eea3ff7aeed869538f1af3608af"] = true //j

	xmlResults["9db54c49a2607ea179054532cdc9ee79"] = true
	xmlResults["811c9eea3ff7aeed869538f1af3608ae"] = true //o
	xmlResults["21d63e291b59b475f308d8885f97af31"] = true //o
	xmlResults["811c9eea3ff7aeed869538f1af3608af"] = true
	xmlResults["c9db4c615350ac5b04c09545ad0ed08e"] = true //o

	for id := range jsonResults {
		if xmlResults[id] {
			fmt.Printf("joined %s\n", id)
			err := processor.PostResult("Joined", id)
			if err != nil {
				fmt.Printf("Error posting joined result: %v\n", err)
			}
		} else {
			fmt.Printf("orphaned %s\n", id)
			err := processor.PostResult("Orphaned", id)
			if err != nil {
				fmt.Printf("Error posting oephaned result: %v\n", err)
			}
		}
	}

	for id := range xmlResults {
		if !jsonResults[id] {
			fmt.Printf("orphaned %s\n", id)
			err := processor.PostResult("Orphaned", id)
			if err != nil {
				fmt.Printf("Error posting orphaned result:%v\n", err)
			}
		}
	}
}
