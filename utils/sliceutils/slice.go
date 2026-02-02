package sliceutils

import (
	"fmt"
	"slices"
)

type SliceOperation string

const (
	OpAppend    SliceOperation = "append"
	OpRemove    SliceOperation = "remove"
	OpOverwrite SliceOperation = "overwrite"
)

func ToStringSlice(slice []any) []string {
	out := []string{}
	for _, v := range slice {
		out = append(out, fmt.Sprintf("%v", v))
	}
	return out
}

func Unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ContainsStringSlice(container, contained []string) bool {

	if len(contained) > len(container) {
		return false
	}

	for _, str := range contained {
		if !slices.Contains(container, str) {
			return false
		}
	}
	return true
}

// Processes operation op to slice ini over values.
// E.g., if op = append, values will be appended to ini;
// if op = overwrite, values will overwrite ini;
// if op = remove, values will be removed from ini.
func ProcessSliceOperation(
	ini []string,
	op SliceOperation,
	values []string,
) []string {

	switch op {
	case OpAppend:
		return append(ini, values...)
	case OpOverwrite:
		return values
	case OpRemove:
		out := make([]string, 0, len(ini))
		for _, s := range ini {
			if !slices.Contains(values, s) {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}

}

// Returns a slice of strings containing the intersection
// between all input string slices.
func IntersectionString(slices ...[]string) []string {

	if len(slices) == 0 {
		return []string{}
	}

	if len(slices) == 1 {
		return slices[0]
	}

	// it's faster to index mapsList than sweep arrays:
	// will keep only 1 array.
	// others will be transformed into mapsList
	mapsList := make(
		[]map[string]struct{}, 0, len(slices)-1,
	)

	for _, slice := range slices[1:] {
		m := map[string]struct{}{}

		for _, str := range slice {
			m[str] = struct{}{}
		}

		mapsList = append(mapsList, m)
	}

	// now, sweep the first slice and try to find
	// in all other maps

	out := make([]string, 0, len(slices[0]))

	// prevent duplicate
	got := map[string]struct{}{}

	for _, v := range slices[0] {

		_, ok := got[v]
		if ok {
			// repeated: go to next
			continue
		}

		found := true
		for _, mp := range mapsList {
			_, ok := mp[v]
			if !ok {
				found = false
				break
			}
		}

		if found {
			out = append(out, v)
			got[v] = struct{}{}
		}
	}

	return out

}
