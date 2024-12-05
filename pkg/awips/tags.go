package awips

import (
	"fmt"
	"regexp"
	"strings"
)

type tagOpt struct {
	Tag       string
	Regexp    string
	Possibles []string
}

var tags = []tagOpt{
	{
		Tag:       "tornado",
		Regexp:    `TORNADO\.\.\.([A-Z ]+)`,
		Possibles: []string{"POSSIBLE", "RADAR INDICATED", "OBSERVED"},
	},
	{
		Tag:       "damage",
		Regexp:    `(TORNADO|THUNDERSTORM|FLASH FLOOD) DAMAGE THREAT\.\.\.([A-Z ]+)`,
		Possibles: []string{"CONSIDERABLE", "DESTRUCTIVE", "CATASTROPHIC"},
	},
	{
		Tag:       "hailThreat",
		Regexp:    `HAIL THREAT\.\.\.([A-Z ]+)`,
		Possibles: []string{"RADAR INDICATED", "OBSERVED"},
	},
	{
		Tag:    "hail",
		Regexp: `.*(HAIL|MAX HAIL SIZE)\.\.\.[><\.0-9]+\s?IN`,
	},
	{
		Tag:       "windThreat",
		Regexp:    `WIND THREAT\.\.\.([A-Z ]+)`,
		Possibles: []string{"RADAR INDICATED", "OBSERVED"},
	},
	{
		Tag:    "wind",
		Regexp: `.*(WIND|MAX WIND GUST)\.\.\.[><\.0-9]+\s?(MPH|KTS)`,
	},
	{
		Tag:       "flashFlood",
		Regexp:    `FLASH FLOOD\.\.\.([A-Z ]+)`,
		Possibles: []string{"RADAR INDICATED", "OBSERVED"},
	},
	{
		Tag:    "expectedRainfall",
		Regexp: `EXPECTED RAINFALL RATE\.\.\.(.)+`,
	},
	{
		Tag:       "damFailure",
		Regexp:    `(DAM|LEVEE) FAILURE\.\.\.(.)+`,
		Possibles: []string{"IMMINENT", "OCCURRING"},
	},
	{
		Tag:       "spout",
		Regexp:    `.*(LANDSPOUT|WATERSPOUT)\.\.\.(.)+`,
		Possibles: []string{"POSSIBLE", "OBSERVED"},
	},
	{
		Tag:       "snowSquall",
		Regexp:    `SNOW SQUALL\.\.\.([A-Z ]+)`,
		Possibles: []string{"RADAR INDICATED", "OBSERVED"},
	},
	{
		Tag:       "snowSquallImpact",
		Regexp:    `SNOW SQUALL IMPACT\.\.\.([A-Z ]+)`,
		Possibles: []string{"SIGNIFICANT"},
	},
}

func ParseTags(text string) (map[string]string, []error) {
	err := []error{}

	output := make(map[string]string)
	for _, tag := range tags {
		regex := regexp.MustCompile(tag.Regexp)
		found := regex.FindString(text)

		if found == "" {
			continue
		}

		value := strings.Split(found, "...")[1]

		if tag.Possibles != nil {
			valid := false
			for _, p := range tag.Possibles {
				if p == value {
					valid = true
					break
				}
			}

			if !valid {
				err = append(err, fmt.Errorf("unusual tag found for %s: %s", tag.Tag, value))
			}
		}

		output[tag.Tag] = value

	}

	return output, err
}
