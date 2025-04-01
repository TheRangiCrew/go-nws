package products

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/TheRangiCrew/go-nws/pkg/awips"
)

type MCD struct {
	Original         string               `json:"original"`
	Number           int                  `json:"number"`
	Issued           time.Time            `json:"issued"`
	Expires          time.Time            `json:"expires"`
	Concerning       string               `json:"concerning"`
	Polygon          awips.PolygonFeature `json:"polygon"`
	WatchProbability int                  `json:"watch_probability"`
	MostProbTornado  string               `json:"most_prob_tornado"`
	MostProbGust     string               `json:"most_prob_gust"`
	MostProbHail     string               `json:"mopst_prob_hail"`
}

func ParseMCD(text string) (*MCD, error) {

	valueRegexp := regexp.MustCompile("([0-9]+)")

	mcdRegex := regexp.MustCompile("(Mesoscale Discussion )([0-9]{4})")
	mcdString := mcdRegex.FindString(text)
	numberString := valueRegexp.FindString(mcdString)
	if numberString == "" {
		return nil, errors.New("error parsing mcd: No MCD number found")
	}
	number, err := strconv.Atoi(numberString)
	if err != nil {
		return nil, fmt.Errorf("error parsing mcd number: %s", err.Error())
	}

	validRegex := regexp.MustCompile("(Valid|VALID) ([0-9]{6}Z) - ([0-9]{6}Z)\n")
	validString := strings.TrimSpace(validRegex.FindString(text))
	timeRegex := regexp.MustCompile("([0-9]{6}Z)")
	times := timeRegex.FindAllString(validString, 2)

	if len(times) != 2 {
		return nil, fmt.Errorf("error parsing mcd: Invalid number of valid times. Found %d, expected 2", len(times))
	}

	issued, err := time.Parse("021504Z", times[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing mcd issued time: %s", err.Error())
	}
	expires, err := time.Parse("021504Z", times[1])
	if err != nil {
		return nil, fmt.Errorf("error parsing mcd expire time: %s", err.Error())
	}

	concerningRegex := regexp.MustCompile(`(Concerning\.\.\.)(.+)`)
	concerningString := concerningRegex.FindString(text)

	if concerningString == "" {
		return nil, fmt.Errorf("error parsing mcd: No concerning text found")
	}

	concerning := strings.ReplaceAll(concerningString, "Concerning...", "")

	latlon, err := awips.ParseLatLon(text)
	if err != nil {
		return nil, fmt.Errorf("error parsing mcd latlon: %s", err.Error())
	}

	polygon := latlon.Polygon

	probabilityRegexp := regexp.MustCompile(`(Probability of Watch Issuance\.\.\.)(.+)`)
	probabilityString := probabilityRegexp.FindString(text)
	var probability int
	if probabilityString != "" {
		valueString := valueRegexp.FindString(probabilityString)

		if valueString == "" {
			return nil, fmt.Errorf("error parsing mcd: Found probability string but no numbers")
		}

		probability, err = strconv.Atoi(valueString)
		if err != nil {
			return nil, fmt.Errorf("error parsing mcd probability: %s", err.Error())
		}
	}

	probTornadoRegexp := regexp.MustCompile(`(MOST PROBABLE PEAK TORNADO INTENSITY\.\.\.)([\w ]+)`)
	probTornadoString := probTornadoRegexp.FindString(text)
	var probTornado string
	if probTornadoString != "" {
		values := strings.Split(probTornadoString, "...")
		if len(values) < 2 {
			return nil, fmt.Errorf("tornado probability string was found but split returned %d elements", len(values))
		}
		probTornado = values[1]
	}

	probGustRegexp := regexp.MustCompile(`(MOST PROBABLE PEAK TORNADO INTENSITY\.\.\.)([\w ]+)`)
	probGustString := probGustRegexp.FindString(text)
	var probGust string
	if probGustString != "" {
		values := strings.Split(probGustString, "...")
		if len(values) < 2 {
			return nil, fmt.Errorf("gust probability string was found but split returned %d elements", len(values))
		}
		probGust = values[1]
	}

	probHialRegexp := regexp.MustCompile(`(MOST PROBABLE PEAK TORNADO INTENSITY\.\.\.)([\w ]+)`)
	probHailString := probHialRegexp.FindString(text)
	var probHail string
	if probHailString != "" {
		values := strings.Split(probHailString, "...")
		if len(values) < 2 {
			return nil, fmt.Errorf("hail probability string was found but split returned %d elements", len(values))
		}
		probHail = values[1]
	}

	mcd := MCD{
		Original:         text,
		Number:           number,
		Issued:           issued,
		Expires:          expires,
		Concerning:       concerning,
		Polygon:          *polygon,
		WatchProbability: probability,
		MostProbTornado:  probTornado,
		MostProbGust:     probGust,
		MostProbHail:     probHail,
	}

	return &mcd, nil
}
