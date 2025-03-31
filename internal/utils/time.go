package utils

import "time"

var Timezones = map[string]*time.Location{
	"GMT":  time.FixedZone("GMT", 0*60*60),
	"UTC":  time.FixedZone("UTC", 0*60*60),
	"AST":  time.FixedZone("AST", -4*60*60),
	"EST":  time.FixedZone("EST", -5*60*60),
	"EDT":  time.FixedZone("EDT", -4*60*60),
	"CST":  time.FixedZone("CST", -6*60*60),
	"CDT":  time.FixedZone("CDT", -5*60*60),
	"MST":  time.FixedZone("MST", -7*60*60),
	"MDT":  time.FixedZone("MDT", -6*60*60),
	"PST":  time.FixedZone("PST", -8*60*60),
	"PDT":  time.FixedZone("PDT", -7*60*60),
	"AKST": time.FixedZone("AKST", -9*60*60),
	"AKDT": time.FixedZone("AKDT", -8*60*60),
	"HST":  time.FixedZone("HST", -10*60*60),
	"SST":  time.FixedZone("SST", -11*60*60),
	"CHST": time.FixedZone("CHST", 10*60*60),
}
