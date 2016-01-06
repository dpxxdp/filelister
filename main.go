package main

import (
	"os"
	"fmt"
	"flag"
	"path/filepath"
	"log"
	"errors"
)

func main() {
	pathFlag := flag.String("path", ".", "path to folder")
	fmt.Printf("pathflag set to: %v\n", *pathFlag)
	recursiveFlag := flag.Bool("r", false, "when set, list files recursively")
	fmt.Printf("recursiveflag set to: %v\n", *recursiveFlag)
	outputFlag := flag.String("output", "text", "<json|yaml|text>")
	fmt.Printf("outputFlag set to: %v\n", *outputFlag)
	flag.Parse()

	path, err := filepath.Abs(*pathFlag)

	if(err != nil) {
		log.Fatal(err)
	}

	output := resolveOutputFlag(*outputFlag)
	if(output == 0) {
		log.Fatal(errors.New("invalid output type; choose <json|yaml|text>"))
	}

	fmt.Printf("output: %v\n", output)

	fileSlice := make([]FileInformation, 0, 100)

	fmt.Printf("walking filepath...\n")
	filepath.Walk(path, WalkAndBuildFileInformation(&fileSlice, *recursiveFlag))

	for _,fileInfo := range fileSlice {
		fmt.Printf("%s\n", fileInfo)
	}
}

func WalkAndBuildFileInformation(filePtr *[]FileInformation, recursive bool) filepath.WalkFunc {
	return filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		fileInfo := FileInformation {
			Path: 			path,
			ModTime: 		info.ModTime(),
			IsLink: 		info.Mode() == os.ModeSymlink,
			IsDir: 			info.IsDir(),
			LinksTo: 		"TODO",
			Size:			info.Size(),
			Name:			info.Name(),
		}

		fileSlice := *filePtr
		*filePtr = append(fileSlice, fileInfo)

		if !recursive {
			return filepath.SkipDir
		} else {
			return nil
		}
    })
}

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