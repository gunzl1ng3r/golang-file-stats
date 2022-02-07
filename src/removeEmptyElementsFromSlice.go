package main

func removeEmptyElementsFromSlice(inputSlice []string) (outputSlice []string) {

	// loop over handed slice and remove anything that is an empty string
	for _, value := range inputSlice {
		if value != "" {
			outputSlice = append(outputSlice, value)
		}
	}

	return
}
