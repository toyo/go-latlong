package latlong

import (
	"math"
	"unicode"
)

// NewLatLongGridLocator is from Grid Locator.
// https://en.wikipedia.org/wiki/Maidenhead_Locator_System
func NewLatLongGridLocator(gl string) *Coordinate {
	latitude := float64(-90)
	longitude := float64(-180)

	latprec := float64(10) * 24
	lonprec := float64(20) * 24

loop:
	for i, c := range gl {
		c = unicode.ToUpper(c)
		switch i % 4 {
		case 0:
			if unicode.IsUpper(c) {
				lonprec /= 24
				longitude += lonprec * float64(c-'A')
			} else {
				break loop
			}
		case 1:
			if unicode.IsUpper(c) {
				latprec /= 24
				latitude += latprec * float64(c-'A')
			} else {
				break loop
			}
		case 2:
			if unicode.IsDigit(c) {
				lonprec /= 10
				longitude += lonprec * float64(c-'0')
			} else {
				break loop
			}
		case 3:
			if unicode.IsDigit(c) {
				latprec /= 10
				latitude += latprec * float64(c-'0')
			} else {
				break loop
			}
		}
	}
	return NewLatLongAlt(latitude+latprec/2, longitude+lonprec/2, latprec, lonprec, nil)
}

// GridLocator is from Grid Locator.
// https://en.wikipedia.org/wiki/Maidenhead_Locator_System
func (latlong *Coordinate) GridLocator() string {
	const floaterr = 1 + 1E-11

	var gl []rune

	latitude := latlong.Lat() + 90
	longitude := latlong.Lng() + 180

	latprec := float64(10) * 24
	lonprec := float64(20) * 24

loop:
	for i := 0; ; i++ {
		switch i % 4 {
		case 0:
			lonprec /= 24
			if lonprec*floaterr < latlong.Size().Lng.Degrees() {
				//fmt.Printf("lon %.15f, %.15f", lonprec, latlong.Size().Lng.Degrees())
				break loop
			}
			c := math.Floor(longitude / lonprec)
			gl = append(gl, rune(byte(c)+'A'))
			longitude -= c * lonprec
		case 1:
			latprec /= 24
			if latprec*floaterr < latlong.Size().Lat.Degrees() {
				//fmt.Printf("lat %.15f, %.15f", latprec, latlong.Size().Lat.Degrees())
				break loop
			}
			c := math.Floor(latitude / latprec)
			gl = append(gl, rune(byte(c)+'A'))
			latitude -= c * latprec
		case 2:
			lonprec /= 10
			if lonprec*floaterr < latlong.Size().Lng.Degrees() {
				//fmt.Printf("lon %.15f, %.15f", lonprec, latlong.Size().Lng.Degrees())
				break loop
			}
			c := math.Floor(longitude / lonprec)
			gl = append(gl, rune(byte(c)+'0'))
			longitude -= c * lonprec

		case 3:
			latprec /= 10
			if latprec*floaterr < latlong.Size().Lat.Degrees() {
				//fmt.Printf("lat %.15f, %.15f", latprec, latlong.Size().Lat.Degrees())
				break loop
			}
			c := math.Floor(latitude / latprec)
			gl = append(gl, rune(byte(c)+'0'))
			latitude -= c * latprec
		}
	}

	l := len(gl)
	if l%2 == 1 {
		gl = gl[:l-1]
	}
	return string(gl)
}
