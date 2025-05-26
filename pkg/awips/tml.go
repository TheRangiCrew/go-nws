package awips

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TML struct {
	Original    string       `json:"original"`
	Time        time.Time    `json:"time"`
	Direction   int          `json:"direction"`
	Speed       int          `json:"speed"`
	SpeedString string       `json:"speedString"`
	Locations   [][2]float64 `json:"location"`
}

func ParseTML(text string, issued time.Time) (*TML, error) {
	tmlRegexp := regexp.MustCompile(`(?m:^(TIME\.\.\.MOT\.\.\.LOC)([A-Za-z0-9 ]*))`)
	original := tmlRegexp.FindString(text)

	if original == "" {
		return nil, nil
	}

	segments := strings.Split(original, " ")[1:]

	parsedTime, err := time.Parse(("1504Z"), segments[0])

	if err != nil {
		return nil, errors.New("could not parse TML time")
	}

	time := time.Date(issued.Year(), issued.Month(), issued.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Now().UTC().Location())

	direction, err := strconv.Atoi(segments[1][:3])

	if err != nil {
		return nil, errors.New("could not parse direction in TML")
	}

	numberRegexp := regexp.MustCompile("[0-9]+")
	speedString := segments[2]
	speed, err := strconv.Atoi(numberRegexp.FindString(speedString))

	if err != nil {
		return nil, errors.New("could not parse speed in TML")
	}

	tml := TML{
		Original:    original,
		Time:        time,
		Direction:   direction,
		Speed:       speed,
		SpeedString: speedString,
		Locations:   [][2]float64{},
	}

	for i := 3; i < len(segments); i += 2 {
		latString, err := strconv.Atoi(segments[i])
		if err != nil {
			return nil, errors.New("failed to parse LAT...LON lat")
		}
		lonString, err := strconv.Atoi(segments[i+1])
		if err != nil {
			return nil, errors.New("failed to parse LAT...LON lon")
		}

		lat := float64(latString) / 100
		lon := (float64(lonString) / 100) * -1
		if lon <= -180.0 {
			lon = lon + 360.0
		}
		location := [2]float64{lon, lat}

		tml.Locations = append(tml.Locations, location)
	}

	return &tml, nil
}
