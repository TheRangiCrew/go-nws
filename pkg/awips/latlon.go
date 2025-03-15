package awips

import (
	"errors"
	"regexp"
	"strconv"
)

type LatLon struct {
	Original string          `json:"original"`
	Points   [][]float64     `json:"points"`
	Polygon  *PolygonFeature `json:"-"`
}

type PolygonFeature struct {
	Type        string        `json:"type"` // Polygon
	Coordinates [][][]float64 `json:"coordinates"`
}

type MultiPolygonFeature struct {
	Type        string          `json:"type"` // MultiPolygon
	Coordinates [][][][]float64 `json:"coordinates"`
}

func ParsePoint(segments []string) (*[]float64, error) {
	if len(segments[0]) > 5 {
		s := segments[0]
		latInit, err := strconv.Atoi(s[0:4])
		if err != nil {
			return nil, errors.New("failed to parse point latitude")
		}

		lonInit, err := strconv.Atoi(s[4:8])
		if err != nil {
			return nil, errors.New("failed to parse point longitude")
		}

		lat := (float64(latInit) / 100)
		lon := (float64(lonInit) / 100) * -1
		if lon > -20.0 {
			lon = lon + -100
		}

		return &[]float64{lon, lat}, nil
	} else {
		if len(segments) > 2 {
			return nil, errors.New("point string was not the correct size")
		}
		latInit, err := strconv.Atoi(segments[0])
		if err != nil {
			return nil, errors.New("failed to parse point latitude")
		}
		lonInit, err := strconv.Atoi(segments[1])
		if err != nil {
			return nil, errors.New("failed to parse point longitude")
		}

		lat := (float64(latInit) / 100)
		lon := (float64(lonInit) / 100) * -1
		if lon <= -180.0 {
			lon = lon + 360.0
		}
		return &[]float64{lon, lat}, nil
	}
}

func ParseLatLon(text string) (*LatLon, error) {
	latlonRegexp := regexp.MustCompile(`(?m)(^LAT\.\.\.LON\s+(\d+\s*)+)`)
	original := latlonRegexp.FindString(text)

	if original == "" {
		return nil, nil
	}

	segmentRegexp := regexp.MustCompile("[0-9]{4,}")
	segments := segmentRegexp.FindAllString(original, -1)

	points := [][]float64{}
	if len(segments[0]) > 5 {
		for _, s := range segments {
			point, err := ParsePoint([]string{s})
			if err != nil {
				return nil, err
			}
			points = append(points, *point)
		}
	} else {
		for i := 0; i < len(segments); i += 2 {
			point, err := ParsePoint([]string{segments[i], segments[i+1]})
			if err != nil {
				return nil, err
			}
			points = append(points, *point)
		}
	}
	if points[0][0] != points[len(points)-1][0] || points[0][1] != points[len(points)-1][1] {
		points = append(points, points[0])
	}

	polygon := PolygonFeature{
		Type:        "Polygon",
		Coordinates: [][][]float64{points},
	}

	return &LatLon{
		Original: original,
		Points:   points,
		Polygon:  &polygon,
	}, nil

}
