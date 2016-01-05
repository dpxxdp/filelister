package main

import (
	"os"
	"fmt"
	"flag"
	"path/filepath"
	"log"
	"errors"
)

type Output int

const (
	text Output = iota + 1
	json
	yaml
)

func resolveOutputFlag(outputFlag string) Output {
	switch outputFlag {
    case "text":
        return text
    case "json":
    	return json
    case "yaml":
    	return yaml
    }
    return 0
}

func main() {
	pathFlag := flag.String("path", "", "path to folder; required")
	recursiveFlag := flag.Bool("r", false, "when set, list files recursively")
	outputFlag := flag.String("output", "text", "<json|yaml|text>")
	flag.Parse()

	if(*pathFlag == "") {
		log.Fatal(errors.New("path required"))
	}

	output := resolveOutputFlag(*outputFlag)
	if(output == 0) {
		log.Fatal(errors.New("invalid output type; choose <json|yaml|text>"))
	}

	filepath.Walk(*pathFlag, WalkAndBuildFileInformation(*recursiveFlag))
}

func WalkAndBuildFileInformation(recursive bool) filepath.WalkFunc {
	return filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		fmt.Println(path)

		if !recursive {
			return filepath.SkipDir
		} else {
			return nil
		}

    })
}

