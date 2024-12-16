package awips

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var VTECClass = map[string]string{
	"O": "Operational",
	"T": "Test",
	"E": "Experimental",
	"X": "Experiemtnal VTEC",
}

var VTECAction = map[string]string{
	"NEW": "issues",
	"CON": "continues",
	"EXA": "expands area to include",
	"EXT": "extends time of",
	"EXB": "extends time and expands area to include",
	"UPG": "issues upgrade to",
	"CAN": "cancels",
	"EXP": "expires",
	"ROU": "routine",
	"COR": "corrects",
}

var VTECSignificance = map[string]string{
	"W": "Warning",
	"Y": "Advisory",
	"A": "Watch",
	"S": "Statement",
	"O": "Outlook",
	"N": "Synopsis",
	"F": "Forecast",
}

var VTECPhenomena = map[string]string{
	"AF": "Ashfall",
	"AS": "Air Stagnation",
	"BH": "Beach Hazard",
	"BS": "Blowing Snow",
	"BW": "Brisk Wind",
	"BZ": "Blizzard",
	"CF": "Coastal Flood",
	"CW": "Cold Weather",
	"DF": "Debris Flow",
	"DS": "Dust Storm",
	"DU": "Blowing Dust",
	"EC": "Extreme Cold",
	"EH": "Excessive Heat",
	"EW": "Extreme Wind",
	"FA": "Flood",
	"FF": "Flash Flood",
	"FG": "Dense Fog",
	"FL": "Flood",
	"FR": "Frost",
	"FW": "Fire Weather",
	"FZ": "Freeze",
	"UP": "Freezing Spray",
	"GL": "Gale",
	"HF": "Hurricane Force Wind",
	"HI": "Inland Hurricane",
	"HS": "Heavy Snow",
	"HT": "Heat",
	"HU": "Hurricane",
	"HW": "High Wind",
	"HY": "Hydrologic",
	"HZ": "Hard Freeze",
	"IP": "Sleet",
	"IS": "Ice Storm",
	"LB": "Lake Effect Snow and Blowing Snow",
	"LE": "Lake Effect Snow",
	"LO": "Low Water",
	"LS": "Lakeshore Flood",
	"LW": "Lake Wind",
	"MA": "Marine",
	"MF": "Marine Dense Fog",
	"MH": "Marine Ashfall",
	"MS": "Marine Dense Smoke",
	"RB": "Small Craft for Rough",
	"RP": "Rip Currents",
	"SB": "Snow and Blowing",
	"SC": "Small Craft",
	"SE": "Hazardous Seas",
	"SI": "Small Craft for Winds",
	"SM": "Dense Smoke",
	"SN": "Snow",
	"SQ": "Snow Squall",
	"SR": "Storm",
	"SS": "Storm Surge",
	"SU": "High Surf",
	"SV": "Severe Thunderstorm",
	"SW": "Small Craft for Hazardous Seas",
	"TI": "Inland Tropical Storm",
	"TO": "Tornado",
	"TR": "Tropical Storm",
	"TS": "Tsunami",
	"TY": "Typhoon",
	"WC": "Wind Chill",
	"WI": "Wind",
	"WS": "Winter Storm",
	"WW": "Winter Weather",
	"XH": "Extreme Heat",
	"ZF": "Freezing Fog",
	"ZR": "Freezing Rain",
}

type VTEC struct {
	Original     string `json:"original"`
	Class        string `json:"class"`
	Action       string `json:"action"`
	WFO          string `json:"wfo"`
	Phenomena    string `json:"phenomena"`
	Significance string `json:"significance"`
	EventNumber  int    `json:"event_number"`
	StartString  string
	Start        *time.Time `json:"start"`
	EndString    string
	End          *time.Time `json:"end"`
}

const VTECRegexp = `([A-Z])\.([A-Z]+)\.([A-Z]+)\.([A-Z]+)\.([A-Z])\.([0-9]+)\.([0-9TZ]+)-([0-9TZ]+)`

func ParseVTEC(text string) ([]VTEC, []error) {
	// Find the VTECs
	vtecRegex := regexp.MustCompile(VTECRegexp)
	instances := vtecRegex.FindAllString(text, -1)

	// There could be more than one
	var vtecs []VTEC
	// We will return an array of errors for debugging individual VTECs instead of failing a whole product parse
	var err []error

	for _, original := range instances {

		segments := strings.Split(original, ".")

		if len(segments) < 6 {
			err = append(err, fmt.Errorf("length of segments is %d, expected 6 for %s", len(segments), original))
			continue
		}

		// Get VTEC class
		class := segments[0]
		if _, ok := VTECClass[class]; !ok {
			err = append(err, fmt.Errorf("invalid class %s for %s", class, original))
			continue
		}

		// Get VTEC action
		action := segments[1]
		if _, ok := VTECAction[action]; !ok {
			err = append(err, fmt.Errorf("invalid action %s for %s", action, original))
			continue
		}

		// Get WFO
		wfo := segments[2]

		// Get phenomena
		phenomena := segments[3]
		if _, ok := VTECPhenomena[phenomena]; !ok {
			err = append(err, fmt.Errorf("invalid phenomena %s for %s", phenomena, original))
			continue
		}

		// Get significance
		significance := segments[4]
		if _, ok := VTECSignificance[significance]; !ok {
			err = append(err, fmt.Errorf("invalid significance %s for %s", significance, original))
			continue
		}

		// Get tracking number
		etnString := segments[5]
		etn, e := strconv.Atoi(etnString)
		if e != nil {
			err = append(err, fmt.Errorf("invalid etn %s for %s", etnString, original))
			continue
		}

		// Get time
		datetimeString := segments[6]
		dateSegments := strings.Split(datetimeString, "-")

		layout := "060102T1504Z"

		var start *time.Time
		var end *time.Time

		zeroRegexp := regexp.MustCompile("000000T0000Z")

		// Sort out start datetime
		if !zeroRegexp.MatchString(dateSegments[0]) {
			t, e := time.Parse(layout, dateSegments[0])
			if e != nil {
				err = append(err, fmt.Errorf("failed to parse start time %s for %s", dateSegments[0], original))
				continue
			}

			start = &t
		}

		if !zeroRegexp.MatchString(dateSegments[1]) {
			t, e := time.Parse(layout, dateSegments[1])
			if e != nil {
				err = append(err, fmt.Errorf("failed to parse start time %s for %s", dateSegments[1], original))
				continue
			}

			end = &t
		}

		vtecs = append(vtecs, VTEC{
			Original:     original,
			Class:        class,
			Action:       action,
			WFO:          wfo,
			Phenomena:    phenomena,
			Significance: significance,
			EventNumber:  etn,
			StartString:  dateSegments[0],
			Start:        start,
			EndString:    dateSegments[1],
			End:          end,
		})
	}

	return vtecs, err
}

func (vtec *VTEC) PhenomenaString() string {
	switch vtec.Phenomena {
	case "FW":
		return "Red Flag"
	}
	return VTECPhenomena[vtec.Phenomena]
}

func (vtec *VTEC) SignificanceString() string {
	return VTECSignificance[vtec.Significance]
}

func (vtec *VTEC) Title(isEmergency bool) string {
	if isEmergency {
		return vtec.PhenomenaString() + " Emergency"
	}
	return vtec.PhenomenaString() + " " + vtec.SignificanceString()
}
