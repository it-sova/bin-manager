package helpers

import (
	"fmt"
	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"sort"
)

// GetLastVersionFromMap sorts map by key returning last element as 2 values
func GetLastVersionFromMap(input map[string]string) (string, string, error) {
	if len(input) == 0 {
		return "", "", fmt.Errorf("unable to get last element of empty map")
	}

	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}

	log.Debugf("Original keys: %#v", keys)

	keys = SortVersions(keys)

	log.Debugf("Sorted keys: %#v", keys)

	lastKey := keys[len(keys)-1]

	lastValue, ok := input[lastKey]
	if !ok {
		return "", "", fmt.Errorf("Failed to get last map elements of %#v", input)
	}

	return lastKey, lastValue, nil
}

func SortVersions(input []string) []string {
	var result []string
	versions := make([]*version.Version, len(input))
	for i, raw := range input {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}

	// After this, the versions are properly sorted
	sort.Sort(version.Collection(versions))

	// Return slice of versions in original format
	for _, sortedVersion := range versions {
		result = append(result, sortedVersion.Original())
	}

	return result
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
