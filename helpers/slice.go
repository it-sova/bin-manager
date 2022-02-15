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
