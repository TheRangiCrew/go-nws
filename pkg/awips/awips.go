package awips

import (
	"errors"
	"regexp"
	"strings"
)

// AWIPS Header
type AWIPS struct {
	Original string `json:"original"`
	Product  string `json:"product"`
	WFO      string `json:"wfo"`
}

var AWIPSRegexp = `(?m:^[A-Z0-9]{4,6}[ ]*\n)`

func HasAWIPS(text string) bool {
	// Find the AWIPS header
	awipsRegex := regexp.MustCompile(AWIPSRegexp)
	original := awipsRegex.FindString(text)
	return original != ""
}

func ParseAWIPS(text string) (AWIPS, error) {
	// Find the AWIPS header
	awipsRegex := regexp.MustCompile(AWIPSRegexp)
	original := awipsRegex.FindString(text)
	if original == "" {
		return AWIPS{}, errors.New("could not find AWIPS header")
	}
	// Trim the end
	original = strings.TrimSpace(original[:len(original)-1])

	// Product is the first three characters
	product := strings.TrimSpace(original[0:3])
	// Issuing office is the final three products
	wfo := strings.TrimSpace(original[3:])

	return AWIPS{
		Original: original,
		Product:  product,
		WFO:      wfo,
	}, nil
}
