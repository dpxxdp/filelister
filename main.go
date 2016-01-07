package main

import (
	"fmt"
	"flag"
	"path/filepath"
	"log"
	"errors"
    "encoding/json"
    
    "gopkg.in/yaml.v2"
)

func main() {
	pathFlag := flag.String("path", ".", "path to folder")
	recursiveFlag := flag.Bool("r", false, "when set, list files recursively")
	outputFlag := flag.String("output", "text", "<json|yaml|text>")
	flag.Parse()

	path, err := filepath.Abs(*pathFlag)
	if(err != nil) {
		log.Fatal(err)
	}

	output := resolveOutputFlag(*outputFlag)
	if(output == 0) {
		log.Fatal(errors.New("invalid output type; choose <json|yaml|text>"))
	}

	var dirInfo DirInfo
	filepath.Walk(path, dirInfo.WalkAndBuildFileInformation(*recursiveFlag))
    
    switch output {
    	case textOut:
        	dirInfo.String()
    	case jsonOut:
        	jsonEncoded, err := json.Marshal(dirInfo.Files)
        	if err!=nil {
        	    log.Fatal(err)
        	}
        	fmt.Println(string(jsonEncoded))
    	case yamlOut:
        	yamlEncoded, _ := yaml.Marshal(dirInfo.Files)
        	if err!=nil {
          		log.Fatal(err)
       		}
			fmt.Println(string(yamlEncoded))
    }
}

//Output records the return format of the filelister
type Output int

const (
	textOut Output = iota + 1
	jsonOut
	yamlOut
)

func resolveOutputFlag(outputFlag string) Output {
	switch outputFlag {
    case "text":
        return textOut
    case "json":
    	return jsonOut
    case "yaml":
    	return yamlOut
    }
    return 0
}