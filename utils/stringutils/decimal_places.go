package stringutils

import (
	"fmt"
	"strconv"
	"strings"
)

// Parses value by placing placesToRight decimal places
// to the right.
// "2.554", 2 -> 255
func ParseValue(
	value string,
	placesToRight int,
) (
	int,
	error,
) {

	if len(value) == 0 {
		return 0, fmt.Errorf("empty input")
	}

	isNegative := value[0] == '-'

	parts := strings.Split(value, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid format")
	} else if len(parts) == 1 {
		decPlaces := ""
		for range placesToRight {
			decPlaces += "0"
		}
		parts = append(parts, decPlaces)
	}

	intPart, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer part: %w", err)
	}

	centsPart := 0
	centsStr := parts[1]

	if len(centsStr) >= placesToRight {
		centsStr = centsStr[0:placesToRight]
	} else {
		//centsStr += "0"
		for range placesToRight - len(centsStr) {
			centsStr += "0"
		}
	}

	if len(centsStr) > 0 {
		centsPart, err = strconv.Atoi(centsStr)
		if err != nil {
			return 0, fmt.Errorf("invalid cents value: %w", err)
		}

		if isNegative {
			centsPart *= -1
		}
	}

	return intPart*intPow(10, placesToRight) +
		centsPart, nil
}

// Formats value by placing placesToLeft decimal places
// to the left.
// 44523, 2 -> 445.23
func FormatValue(
	value int,
	placesToLeft int,
) string {

	if placesToLeft == 0 {
		return fmt.Sprintf("%v", value)
	}

	intPart := value / intPow(10, placesToLeft)
	centsPart := value % intPow(10, placesToLeft)
	sign := ""
	if value < 0 {
		sign = "-"
		intPart = -intPart
		centsPart = -centsPart
	}
	return fmt.Sprintf("%s%d.%0*d",
		sign, intPart, placesToLeft, centsPart)
}

func intPow(
	n int,
	m int,
) int {

	if m == 0 {
		return 1
	}

	if m == 1 {
		return n
	}

	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result

}
