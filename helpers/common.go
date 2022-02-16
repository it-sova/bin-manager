package helpers

import (
	"fmt"
	"sort"
)

// GetLastMapElement sorts map by key returning last element as 2 values
func GetLastMapElement(input map[string]string) (string, string, error) {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	lastKey := keys[len(keys)-1]

	lastValue, ok := input[lastKey]
	if !ok {
		return "", "", fmt.Errorf("Failed to get last map elements of %#v", input)
	}

	return lastKey, lastValue, nil
}

// RemoveEmptyElementsFromStringSlice removes empty elements from slice
func RemoveEmptyElementsFromStringSlice(input []string) []string {
	var res []string
	for _, element := range input {
		if element != "" {
			res = append(res, element)
		}
	}
	return res
}

// StringSliceHasElement returns true if slice has element
func StringSliceHasElement(slice []string, element string) bool {
	for _, value := range slice {
		if element == value {
			return true
		}
	}

	return false
}
