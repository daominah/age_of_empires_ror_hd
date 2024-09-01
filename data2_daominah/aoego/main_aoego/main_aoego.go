package main

import (
	"flag"
	"log"
)

// inputFilePath is the path to the "*.ai" file that defines the strategy
var inputFilePath string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	flag.StringVar(&inputFilePath,
		"i",
		`D:\game\age_of_empires_ror_hd\data2_daominah\Assyria_Archer.ai`,
		`the path to the "*.ai" file that defines the strategy`)
	flag.Parse()
	if inputFilePath == "" {
		log.Printf("empty input file path")
		return
	}

	log.Printf("input file path: %v", inputFilePath)
}
