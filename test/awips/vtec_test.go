package awips_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/stretchr/testify/assert"
)

// Test cases mostly derived from NWS Directive 10-1703 as of 1 October, 2024
// https://www.weather.gov/media/directives/010_pdfs/pd01017003curr.pdf
func TestVTEC(t *testing.T) {

	arr, err := awips.ParseVTEC("/O.NEW.KMKX.GL.A.0002.111203T0000Z-111203T1200Z/")
	assert.Empty(t, err)
	assert.NotEmpty(t, arr, "vtec array should not be empty ")

	vtec := arr[0]
	assert.Equalf(t, "O", vtec.Class, "vtec class should be O, got %v", vtec.Class)
	assert.Equalf(t, "NEW", vtec.Action, "vtec action should be NEW, got %v", vtec.Action)
	assert.Equalf(t, "KMKX", vtec.WFO, "vtec wfo should be KMKX, got %v", vtec.WFO)
	assert.Equalf(t, "GL", vtec.Phenomena, "vtec phenomena should be GL, got %v", vtec.Phenomena)
	assert.Equalf(t, "A", vtec.Significance, "vtec significance should be A, got %v", vtec.Significance)
	assert.Equalf(t, 2, vtec.EventNumber, "vtec event number should be 2, got %d", vtec.EventNumber)
	assert.Equalf(t, time.Date(2011, time.December, 03, 0, 0, 0, 0, time.UTC), vtec.Start, "vtec start should be UTC 00:00 3 December, 2011, got %v", vtec.Start)
	assert.Equalf(t, time.Date(2011, time.December, 03, 12, 0, 0, 0, time.UTC), vtec.End, "vtec end should be UTC 12:00 3 December, 2011, got %v", vtec.Start)

	arr, err = awips.ParseVTEC("/O.EXT.KICT.FL.W.0007.000000T0000Z-000000T0000Z/")
	assert.Empty(t, err)
	assert.NotEmpty(t, arr, "vtec array should not be empty ")

	vtec = arr[0]
	assert.Equal(t, "O", vtec.Class)
	assert.Equal(t, "EXT", vtec.Action)
	assert.Equal(t, "KICT", vtec.WFO)
	assert.Equal(t, "FL", vtec.Phenomena)
	assert.Equal(t, "W", vtec.Significance)
	assert.Equal(t, 7, vtec.EventNumber)
	assert.Equal(t, time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), vtec.Start)
	assert.Equal(t, time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), vtec.End)

	arr, err = awips.ParseVTEC("/O.CON.KLWX.CF.Y.0004.000000T0000Z-110408T1500Z/")
	assert.Empty(t, err)
	assert.NotEmpty(t, arr, "vtec array should not be empty ")

	vtec = arr[0]
	assert.Equal(t, "O", vtec.Class)
	assert.Equal(t, "CON", vtec.Action)
	assert.Equal(t, "KLWX", vtec.WFO)
	assert.Equal(t, "CF", vtec.Phenomena)
	assert.Equal(t, "Y", vtec.Significance)
	assert.Equal(t, 4, vtec.EventNumber)
	assert.Equal(t, time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), vtec.Start)
	assert.Equal(t, time.Date(2011, time.April, 8, 15, 0, 0, 0, time.UTC), vtec.End)

	arr, err = awips.ParseVTEC("O.CON.KLWX.CF.Y.0004.000000T0000Z-110408T1500Z/")
	assert.Empty(t, err)
	assert.NotEmpty(t, arr, "vtec array should not be empty ")
	arr, err = awips.ParseVTEC("/O.CON.KLWX.CF.Y.0004.000000T0000Z-110408T1500Z")
	assert.Empty(t, err)
	assert.NotEmpty(t, arr, "vtec array should not be empty ")
	arr, err = awips.ParseVTEC("O.CON.KLWX.CF.Y.0004.000000T0000Z-110408T1500Z")
	assert.Empty(t, err)
	assert.NotEmpty(t, arr, "vtec array should not be empty ")
	// Fail for missing timezone Z
	arr, err = awips.ParseVTEC("O.CON.KLWX.CF.Y.0004.000000T0000Z-110408T1500")
	assert.NotEmpty(t, err)
	assert.Empty(t, arr, "vtec array should be empty ")
	// Fail for invalid class
	arr, err = awips.ParseVTEC("/L.CON.KLWX.CF.Y.0004.000000T0000Z-110408T1500Z/")
	assert.NotEmpty(t, err)
	assert.Empty(t, arr, "vtec array should be empty ")
	// Fail for invalid significance
	arr, err = awips.ParseVTEC("/O.CON.KLWX.CF.T.0004.000000T0000Z-110408T1500Z/")
	assert.NotEmpty(t, err)
	assert.Empty(t, arr, "vtec array should be empty ")
	// Fail for invalid segment count
	arr, err = awips.ParseVTEC("/O.CON.KLWX.CF.0004.000000T0000Z-110408T1500Z/")
	assert.Empty(t, err)
	assert.Empty(t, arr, "vtec array should be empty ")
}

// Test cases mostly derived from NWS Directive 10-1703 as of 1 October, 2024
// https://www.weather.gov/media/directives/010_pdfs/pd01017003curr.pdf
func TestMultiVTEC(t *testing.T) {
	text := `
/O.NEW.KPSR.DU.Y.0001.120105T1400Z-120106T0400Z/
/O.EXT.KPSR.WI.Y.0001.120105T1000Z-120106T0400Z/
	`
	arr, err := awips.ParseVTEC(text)
	assert.Empty(t, err)
	assert.Len(t, arr, 2)

	// 1 VTEC, 1 H-VTEC
	text = `
/O.EXT.KICT.FL.W.0010.000000T0000Z-000000T0000Z/
/CFVK1.3.ER.110629T1727Z.110702T0000Z.000000T0000Z.NR/
	`
	arr, err = awips.ParseVTEC(text)
	assert.Empty(t, err)
	assert.Len(t, arr, 1)

	// VTEC 1 malformed, VTEC 2 parsed
	text = `
/O.NEW.KPSR.DU.Y.0001.120105T1400Z-120106T0400Z/
/O.EX.KPSR.WI.Y.0001.120105T1000Z-120106T0400Z/
	`
	arr, err = awips.ParseVTEC(text)
	assert.NotEmpty(t, err)
	assert.Len(t, arr, 1)
}

func TestProduct(t *testing.T) {
	file, err := os.Open("../../assets/awips/test/TCVCAE.txt")
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	assert.Greater(t, len(text), 0)

	arr, er := awips.ParseVTEC(text)
	assert.Empty(t, er)
	assert.Len(t, arr, 18)
}
