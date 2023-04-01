package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No destination file provided")
		return
	}
	file := os.Args[1]

	table := newSinTable(file)
	if err := table.execute(); err != nil {
		fmt.Println("Template Execution Error: ", err)
		return
	}
}
