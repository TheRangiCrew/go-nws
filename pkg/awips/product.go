package awips

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/TheRangiCrew/go-nws/internal/utils"
)

/*
Definitions and components of an AWIPS text product are described in
NWS Directive 10-1701 as of September 30, 2024.

https://www.weather.gov/media/directives/010_pdfs/pd01017001curr.pdf
*/

// An AWIPS text product
type TextProduct struct {
	Text     string               `json:"text"`
	WMO      WMO                  `json:"wmo"`
	AWIPS    AWIPS                `json:"awips"`
	Issued   time.Time            `json:"issued"`
	Office   string               `json:"office"`
	Product  string               `json:"product"`
	Segments []TextProductSegment `json:"segments"`
}

// A text product segment
type TextProductSegment struct {
	Text    string            `json:"text"`
	VTEC    []VTEC            `json:"vtec"`
	UGC     *UGC              `json:"ugc"`
	Expires time.Time         `json:"expires"` // The product expiry time as defined in NWS Directive 10-1701
	Ends    time.Time         `json:"ends"`    // The event end time as defined in NWS Directive 10-1701
	LatLon  *LatLon           `json:"latlon"`
	Tags    map[string]string `json:"tags"`
	TML     *TML              `json:"tml"`
}

// Attempts to parse the given text into a text product including segments & VTEC
func New(text string) (*TextProduct, error) {

	// Get the AWIPS header
	awips, err := ParseAWIPS(text)
	if err != nil {
		return nil, err
	}

	issued, err := GetIssuedTime(text)
	if err != nil {
		return nil, err
	}

	// Get the WMO header
	wmo, err := ParseWMO(text)
	if err != nil {
		return nil, err
	}

	// TODO: Decide if we actually need this
	// bilRegexp := regexp.MustCompile("(?m:^(BULLETIN - |URGENT - |EAS ACTIVATION REQUESTED|IMMEDIATE BROADCAST REQUESTED|FLASH - |REGULAR - |HOLD - |TEST...)(.*))")
	// bil := bilRegexp.FindString(text)

	segments, e := GetSegments(text, issued, wmo)
	if len(e) != 0 {
		return nil, e[0]
	}

	product := TextProduct{
		Text:     text,
		WMO:      wmo,
		AWIPS:    awips,
		Issued:   issued,
		Office:   awips.WFO,
		Product:  awips.Product,
		Segments: segments,
	}

	return &product, nil
}

/*
Attempts to find a product issuing datetime string in the provided text. If a match is found, it may be disseminated. Otherwise, if all else fails, returns time zero.
*/
func GetIssuedTime(text string) (time.Time, error) {
	var issued time.Time
	var err error

	// Find when the product was issued
	issuedRegexp := regexp.MustCompile("[0-9]{3,4} ((AM|PM) [A-Za-z]{3,4}|UTC) ([A-Za-z]{3} ){2}[0-9]{1,2} [0-9]{4}")
	issuedString := issuedRegexp.FindString(text)

	if issuedString != "" {
		// Find if the timezone is UTC
		utcRegexp := regexp.MustCompile("UTC")
		utc := utcRegexp.MatchString(issuedString)
		if utc {
			// Set the UTC timezone
			issued, err = time.ParseInLocation("1504 UTC Mon Jan 2 2006", issuedString, utils.Timezones["UTC"])
		} else {
			/*
				Since the time package cannot handle the time format that is provided in the NWS text products,
				we have to modify the string to include a better seperator between the hour and the minute values
			*/
			tzString := strings.ToUpper(strings.Split(issuedString, " ")[2])
			tz := utils.Timezones[tzString]
			if tz == nil {
				return issued, errors.New("missing timezone " + tzString + " in issued string")
			}
			split := strings.Split(issuedString, " ")
			t := split[0]
			hours := t[:len(t)-2]
			minutes := t[len(t)-2:]
			split[0] = hours + ":" + minutes
			new := strings.Join(split, " ")
			new = strings.Replace(new, tzString+" ", "", -1)
			issued, err = time.ParseInLocation("3:04 PM Mon Jan 2 2006", new, tz)
		}

		if err != nil {
			return issued, errors.New("could not parse issued date line")
		}
	}

	return issued, nil
}

func GetSegments(text string, issued time.Time, wmo WMO) ([]TextProductSegment, []error) {
	// Segment the product
	splits := strings.Split(text, "$$")

	segments := []TextProductSegment{}
	errors := []error{}

	for _, segment := range splits {
		segment = strings.TrimSpace(segment)

		// Assume the segment is the end of the product if it is shorter than 10 characters
		if len(segment) < 20 {
			continue
		}

		ugc, err := ParseUGC(segment)
		if err != nil {
			errors = append(errors, err)
			return nil, errors
		}
		expires := time.Now().UTC()
		if ugc != nil {
			expires = time.Date(issued.Year(), issued.Month(), ugc.Expires.Day(), ugc.Expires.Hour(), ugc.Expires.Minute(), 0, 0, time.UTC)
			if ugc.Expires.Day() > wmo.Issued.Day() && ugc.Expires.Day() == 1 {
				expires = expires.AddDate(0, 1, 0)
			}
			ugc.Merge(issued)
		}

		// Find any VTECs that the segment may have
		vtec, e := ParseVTEC(segment)
		if len(e) != 0 {
			errors = append(errors, e...)
		}

		latlon, err := ParseLatLon(text)
		if err != nil {
			errors = append(errors, err)
			return nil, errors
		}

		tags, e := ParseTags(text)
		if len(e) != 0 {
			errors = append(errors, e...)
		}

		segments = append(segments, TextProductSegment{
			Text:    segment,
			VTEC:    vtec,
			UGC:     ugc,
			Expires: expires,
			LatLon:  latlon,
			Tags:    tags,
		})

	}

	return segments, nil
}

func (product *TextProduct) HasVTEC() bool {
	for _, segment := range product.Segments {
		if segment.HasVTEC() {
			return true
		}
	}
	return false
}

func (segment *TextProductSegment) HasVTEC() bool {
	return len(segment.VTEC) != 0
}

func (segment *TextProductSegment) HasUGC() bool {
	return segment.UGC != nil
}

func (segment *TextProductSegment) IsEmergency() bool {
	emergencyRegexp := regexp.MustCompile(`(TORNADO|FLASH\s+FLOOD)\s+EMERGENCY`)
	return emergencyRegexp.MatchString(segment.Text)
}

func (segment *TextProductSegment) IsPDS() bool {
	pdsRegexp := regexp.MustCompile(`(THIS\s+IS\s+A|This\s+is\s+a)\s+PARTICULARLY\s+DANGEROUS\s+SITUATION`)
	return pdsRegexp.MatchString(segment.Text)
}
