package awips

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UGC struct {
	Original string    `json:"original"`
	States   []State   `json:"states"`
	Expires  time.Time `json:"expires"`
}

type State struct {
	ID    string   `json:"id"`   // State ID
	Type  string   `json:"type"` // "C" for counties or "Z" for zones
	Areas []string `json:"areas"`
}

func ParseUGC(text string) (*UGC, error) {
	// Find the start of the UGC
	ugcStartRegex := regexp.MustCompile("(?m:^[A-Z]{2}(C|Z)[A-Z0-9]{3}(-|>))")
	startIndex := ugcStartRegex.FindStringIndex(text)
	if startIndex == nil {
		return nil, nil
	}
	start := text[startIndex[0]:]

	// Find the end of the UGC
	ugcEndRegex := regexp.MustCompile("([0-9]{6}-)")
	endIndex := ugcEndRegex.FindStringIndex(start)
	if endIndex == nil {
		return nil, nil
	}

	// Subtract 1 to remove the - at the end of the UGC
	original := start[:endIndex[1]-1]

	// Make it one long string and segment it
	segments := strings.Split(strings.ReplaceAll(original, "\n", ""), "-")

	// Get the expiry time
	expires, err := time.Parse("021504", segments[len(segments)-1])
	if err != nil {
		return nil, errors.New("could not parse UGC datetime")
	}
	segments = segments[:len(segments)-1]

	// Group everything into states since that is the order of the UGC
	states := []State{}
	currentState := -1
	alphabetRegexp := regexp.MustCompile("[A-Z]")
	// UGC uses > to specify a range of zones/counties
	bracketRegexp := regexp.MustCompile(">")

	for _, s := range segments {
		// If the
		if alphabetRegexp.MatchString(s) {
			currentState++
			states = append(states, State{
				ID:    s[0:2],
				Type:  s[2:3],
				Areas: []string{},
			})
			s = s[3:]
		}

		// Get the range of the zones/counties
		if bracketRegexp.MatchString(s) {
			start, err := strconv.Atoi(s[:3])
			if err != nil {
				return nil, errors.New("could not parse UGC int")
			}

			end, err := strconv.Atoi(s[4:])
			if err != nil {
				return nil, errors.New("could not parse UGC int")
			}

			for i := start; i <= end; i++ {
				// Format the ugc to be at least three digits padded with zeros
				states[currentState].Areas = append(states[currentState].Areas, fmt.Sprintf("%03d", i))
			}
		} else {
			states[currentState].Areas = append(states[currentState].Areas, s)
		}
	}

	return &UGC{
		Original: original,
		States:   states,
		Expires:  expires,
	}, nil
}

func (ugc *UGC) Merge(t time.Time) {
	ugc.Expires = time.Date(t.Year(), t.Month(), ugc.Expires.Day(), ugc.Expires.Hour(), ugc.Expires.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}
