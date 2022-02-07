package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
)

func main() {

	debug := flag.Bool("debug", false, "enable debug mode")
	baseDirectory := flag.String("baseDirectory", "./", "path to directory you want to evaluate")
	outputFilePath := flag.String("outputFilePath", "./", "path to the director where outputFile should be written to")

	// mandatory call so flags are really usable
	flag.Parse()

	// create name of outputFile based on baseDirectory (replace / with _)
	outputFile := strings.Join(removeEmptyElementsFromSlice(strings.Split(*baseDirectory, "/")), "_")

	var objects = make(map[string]int64)
	var baseDirectoryDepth = len(strings.Split(*baseDirectory, "/"))
	var fileContent []string
	fileContent = append(fileContent, "# HELP file_size_total The size of a given path in bytes.")
	fileContent = append(fileContent, "# TYPE file_size_total gauge")

	if *debug {
		log.Println("DEBUG - baseDirectory is set to: " + *baseDirectory)
		log.Println("DEBUG - baseDirectoryDepth is set to: " + strconv.Itoa(baseDirectoryDepth))
	}

	// add the start directory with a size of zero to avoid surprises
	objects[*baseDirectory] = 0

	// call function to get dir size recursively
	objects = parseDirContents(*baseDirectory, objects, *debug)

	// iterate over the returned values
	for key, value := range objects {
		// if the number of "/" is identical to those of the baseDirectory we have a file that can be added
		if len(strings.Split(key, "/")) == baseDirectoryDepth {
			if *debug {
				log.Println(key + " is " + strconv.FormatInt(value, 10) + " large")
			}
			// fileContent = append(fileContent, "file_size_total{path=\""+key+"\"} "+strconv.FormatInt(value, 10)+" "+strconv.Itoa(int(time.Now().Unix())))
			fileContent = append(fileContent, "file_size_total{path=\""+key+"\"} "+strconv.FormatInt(value, 10))
		}

		// directories in baseDirectory have an additional "/"
		if len(strings.Split(key, "/")) == baseDirectoryDepth+1 {
			// only if the last character is a "/" we have a subdirectory (everything else could be a file in the subdirectory)
			if key[len(key)-1:] == "/" {
				if *debug {
					log.Println(key + " is " + strconv.FormatInt(value, 10) + " large")
				}
				// fileContent = append(fileContent, "file_size_total{path=\""+key+"\"} "+strconv.FormatInt(value, 10)+" "+strconv.Itoa(int(time.Now().Unix())))
				fileContent = append(fileContent, "file_size_total{path=\""+key+"\"} "+strconv.FormatInt(value, 10))
			}
		}

		// finally call function to write results to file
		writeMetricFile(*outputFilePath+outputFile+".prom", fileContent, *debug)
	}
}
