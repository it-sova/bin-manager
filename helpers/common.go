package helpers

import (
	"fmt"
	"sort"
)

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
