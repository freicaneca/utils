package stringutils

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"math/rand"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letters_numbers = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// IsEmptyString checks if the string "s" entered is empty
func IsEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// GetStringInBetween Returns empty string if no start string found
func GetStringInBetween(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]

}

func ImplodeInt(sep string, intSl []int) string {
	var out string

	if len(intSl) > 0 {
		out += fmt.Sprintf("%d", intSl[0])
		for i := 1; i < len(intSl); i++ {
			out += fmt.Sprintf("%s%d", sep, intSl[i])
		}
	}

	return out
}

func ImplodeString(sep string, strSl []string) string {
	var out string

	if len(strSl) > 0 {
		out += fmt.Sprint(strSl[0])
		for i := 1; i < len(strSl); i++ {
			out += fmt.Sprintf("%s%s", sep, strSl[i])
		}
	}

	return out
}

func NormalizeString(text string) string {

	t := transform.Chain(
		//norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
		norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, err := transform.String(t, text)

	if err != nil {
		return "" //, fmt.Errorf("failed to transform string %#v", err)
	}

	re := regexp.MustCompile(`[^0-9A-za-z \\/ \\s]`)

	return re.ReplaceAllString(result, "") //, nil

}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

// RandomNumbersString returns a string of "n" random numbers
func RandomNumbersString(size int) string {
	randomNumbers := rand.Perm(100)
	out := ImplodeInt("", randomNumbers)[:size]
	return out
}

// RandomCharsString returns a string of "n" random letters,
// possibly mixed-case.
func RandomLettersString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandomString returns a string of "n" random letters or numbers.
func RandomString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = letters_numbers[rand.Intn(len(letters_numbers))]
	}
	return string(b)
}

func Hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
