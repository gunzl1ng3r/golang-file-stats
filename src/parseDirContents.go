package main

import (
	"log"
	"os"
	"strconv"
)

func parseDirContents(directory string, objects map[string]int64, debug bool) (returnObjects map[string]int64) {

	// set working variables
	returnObjects = objects

	// Open the directory.
	if debug {
		log.Println("DEBUG - Opening " + directory)
	}
	outputDirRead, _ := os.Open(directory)

	// Call Readdir to get all files.
	outputDirFiles, _ := outputDirRead.Readdir(0)

	outputDirRead.Close()

	for _, entry := range outputDirFiles {
		if entry.IsDir() {
			// upon finding a directory, we need to recuse down that directory
			fileName := directory + entry.Name() + "/"
			if debug {
				log.Println("DEBUG - Directory: " + fileName)
			}
			tempObjects := parseDirContents(fileName, returnObjects, debug)
			if debug {
				log.Println("DEBUG: Directory: Increase " + directory + " by " + strconv.FormatInt(tempObjects[directory], 10) + " from " + directory)
			}
			// increase the directory we have been called for by the size of the directory we just parsed
			returnObjects[directory] = returnObjects[directory] + tempObjects[fileName]
		} else {
			// files are easy, simply add their size to the direcory total
			fileName := directory + entry.Name()
			returnObjects[fileName] = entry.Size()
			if debug {
				log.Println("DEBUG - File: " + fileName)
				log.Println("DEBUG - File: Increase " + directory + " by " + strconv.FormatInt(returnObjects[fileName], 10) + " from " + fileName)
			}
			// increase the directory we have been called for by the size of the file we just found
			returnObjects[directory] = returnObjects[directory] + returnObjects[fileName]
		} // endElse
	} // end range outputDirFiles

	if debug {
		log.Println("DEBUG - Ending run of parseDirContents for " + directory)
	}
	return
}
