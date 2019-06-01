package main

import (
	"fmt"
	"log"
	"os"

	"github.com/das-frama/website/pkg/markdown"
)

func main() {
	file, err := os.Open("new_day.md")
	if err != nil {
		log.Fatal(err)
	}

	output, err := markdown.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(output)
}
