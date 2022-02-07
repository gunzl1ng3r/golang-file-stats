package main

import (
	"fmt"
	"os"
)

func writeMetricFile(metricFilePath string, content []string, debug bool) {

	// create the file
	metricFile, err := os.Create(metricFilePath)
	if err != nil {
		fmt.Println(err)
	}
	// close the file with defer
	defer metricFile.Close()

	// iterate over the passed slice and write its content to file
	for key := range content {
		metricFile.WriteString(content[key] + "\n")
	}
}
