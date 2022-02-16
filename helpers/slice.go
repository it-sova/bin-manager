package helpers

func RemoveEmptyElementsFromStringSlice(input []string) []string {
	var res []string
	for _, element := range input {
		if element != "" {
			res = append(res, element)
		}
	}
	return res
}

func StringSliceHasElement(slice []string, element string) bool {
	for _, value := range slice {
		if element == value {
			return true
		}
	}

	return false
}
