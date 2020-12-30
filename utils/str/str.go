package str

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//ParseLocalKey return pattern word like "key.word1_word2" into "Word1 Word2"
func ParseLocalKey(localKey string) string {
	dotSeparatedWord := strings.Split(localKey, ".")                  // "key.word1_word2" -> ["key", "word1_word2"]
	splittedUnderscoreWord := strings.Split(dotSeparatedWord[1], "_") // "word1_word2" -> ["word1", "word2"]
	joinedWord := strings.Join(splittedUnderscoreWord[:], " ")        // ["word1", "word2"] -> "word1 word2"
	titelizeWord := strings.Title(joinedWord)                         // "word1 word2" -> "Word1 Word2"
	return titelizeWord
}

func ZeroToCountryCode(phoneNumber string, countryCode string) string {
	if strings.HasPrefix(phoneNumber, strings.TrimPrefix(countryCode, "+")) {
		return phoneNumber
	}
	return strings.TrimPrefix(countryCode, "+") + strings.TrimPrefix(phoneNumber, "0")
}

//GenerateKeyFromBaseURL return key for generating presign url
func GenerateKeyFromBaseURL(url string, delimiter string) string {
	if len(url) < 1 {
		return ""
	}
	startDelimiter := len(delimiter)
	endDelimiter := len(url)

	return url[startDelimiter:endDelimiter]
}

// GenerateLocaleKey return pattern locale key. e.g: payment_method.bank_transfer
func GenerateLocaleKey(domain string, localKey string) string {
	localKey = strings.ToLower(localKey)

	return fmt.Sprintf("%s.%s", domain, localKey)
}

// IsStringExistInArray check is input string exist in array of string or not
func IsStringExistInArray(searchStr string, list []string) bool {
	for _, str := range list {
		if searchStr == str {
			return true
		}
	}
	return false
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func StringToArray(str string, delimiter string) []string {
	return strings.Split(str, delimiter)
}

func RemoveWhiteSpace(str string) string {
	return strings.Join(strings.Fields(str), "")
}

// SplitSorter ...
func SplitSorter(sorter string, formatForSQL bool) (sortField string, sortType string, ok bool) {
	separatorIndex := strings.LastIndex(sorter, "_")

	if separatorIndex == -1 {
		return "", "", false
	}

	sortField = sorter[:separatorIndex]
	sortType = sorter[separatorIndex+1:]

	if sortType == "ascend" {
		sortType = "asc"
		if formatForSQL {
			sortType = "ASC"
		}
	} else if sortType == "descend" {
		sortType = "desc"
		if formatForSQL {
			sortType = "DESC NULLS LAST"
		}
	}

	return sortField, sortType, true
}

// StringToFloat64 ...
func StringToFloat64(s string) float64 {
	result, _ := strconv.ParseFloat(s, 64)
	return result
}

// EqualCaseInsensitive check if two string is equal in the lowercase character
func EqualCaseInsensitive(str1 string, str2 string) bool {
	return strings.EqualFold(strings.ToLower(str1), strings.ToLower(str2))
}

// IsStringExistInArrayCaseInsensitive check is input string exist in array compared in lowercase
func IsStringExistInArrayCaseInsensitive(searchStr string, list []string) bool {
	for _, str := range list {
		if EqualCaseInsensitive(searchStr, str) {
			return true
		}
	}
	return false
}

// Capitalize capitalizes the first character of the string.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}
