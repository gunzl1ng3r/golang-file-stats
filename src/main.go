package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	queryDuration := time.Now().UnixMilli()

	debug := flag.Bool("debug", false, "enable debug mode")
	baseDirectory := flag.String("baseDirectory", "./", "path to directory you want to evaluate")
	outputFilePath := flag.String("outputFilePath", "./", "path to the director where outputFile should be written to")

	// mandatory call so flags are really usable
	flag.Parse()

	// create name of outputFile based on baseDirectory (replace / with _)
	outputFile := strings.Join(removeEmptyElementsFromSlice(strings.Split(*baseDirectory, "/")), "_")
	lockFilePath := *outputFilePath + outputFile + ".lock"

	// check if lockfile exists
	_, err := os.Stat(outputFile + ".lock")

	// if stat on lockfile returns "file does not exist", we are good to continue
	if errors.Is(err, os.ErrNotExist) {
		if *debug {
			log.Println("DEBUG - creating lockfile (" + lockFilePath + ") to prevent concurrent runs.")
			lockFile, err := os.Create(lockFilePath)
			if err != nil {
				log.Println("ERROR - Somthing went wrong trying to create lockfile" + lockFilePath)
				log.Panic(err)
			}
			defer lockFile.Close()
		}
	} else {
		// if we find a lockfile, we abort
		log.Println("INFO - Job already running. Exiting.")
		os.Exit(0)
	}

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

	}

	fileContent = append(fileContent, "# HELP file_size_query_duration The time in milliseconds it took to get the size of a given path.")
	fileContent = append(fileContent, "# TYPE file_size_query_duration gauge")
	queryDuration = time.Now().UnixMilli() - queryDuration
	fileContent = append(fileContent, "file_size_query_duration{path=\""+*baseDirectory+"\"} "+strconv.FormatInt(queryDuration, 10))

	// finally call function to write results to file
	writeMetricFile(*outputFilePath+outputFile+".prom", fileContent, *debug)

	err = os.Remove(lockFilePath)
	if err != nil {
		log.Println("ERROR - failed to delete lockfile (" + lockFilePath + ")")
	} else {
		if *debug {
			log.Println("DEBUG - successfully removed lockfile (" + lockFilePath + ")")
		}
	}
}
