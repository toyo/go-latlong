package latlong

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Coordinate is LatLng and Altitude.
type Coordinate struct {
	LatLng
	alt *float64 // altitude
}

// MarshalJSON is a marshaler for JSON.
func (latlong *Coordinate) MarshalJSON() ([]byte, error) {
	s := latlong.lngString() + "," + latlong.latString()
	if latlong.alt != nil {
		s += "," + latlong.altString()
	}
	return []byte("[" + s + "]"), nil
}

// UnmarshalJSON is a unmarshaler for JSON.
func (latlong *Coordinate) UnmarshalJSON(data []byte) (err error) {
	s := strings.TrimSpace(string(data))

	if s[0] != '[' {
		return errors.New("Unknown JSON format (not starting '[')")
	}
	s = s[1:]

	if s[len(s)-1] != ']' {
		return errors.New("Unknown JSON format (not ending ']')")
	}
	s = s[:len(s)-1]

	v := strings.Split(s, ",")
	switch len(v) {
	case 2:
		lat, latprec := getLat(v[1])
		if isErrorDeg(lat, latprec) {
			err = fmt.Errorf("Error Degreee on JSON Lat %s", v[1])
		}
		lng, lngprec := getLng(v[0])
		if isErrorDeg(lng, lngprec) {
			err = fmt.Errorf("Error Degreee on JSON Lng %s", v[0])
		}

		*latlong = *NewLatLongAlt(lat, lng, latprec, lngprec, nil)
	case 3:
		lat, latprec := getLat(v[1])
		if isErrorDeg(lat, latprec) {
			err = fmt.Errorf("Error Degreee on JSON Lat %s", v[1])
		}
		lng, lngprec := getLng(v[0])
		if isErrorDeg(lng, lngprec) {
			err = fmt.Errorf("Error Degreee on JSON Lng %s", v[0])
		}

		*latlong = *NewLatLongAlt(lat, lng, latprec, lngprec, getAlt(v[2]))
	default:
		return errors.New("unknown JSON Coordinate format")
	}

	return
}

// NewLatLongAlt is from latitude, longitude and altitude.
func NewLatLongAlt(latitude, longitude, latprec, longprec float64, altitude *float64) *Coordinate {
	var latlongalt Coordinate
	latlongalt.LatLng.LatLng = s2.LatLngFromDegrees(latitude, longitude)
	latlongalt.latprec = s1.Angle(latprec) * s1.Degree
	latlongalt.lngprec = s1.Angle(longprec) * s1.Degree
	latlongalt.alt = altitude
	return &latlongalt
}

// NewLatLongISO6709 is from ISO6709 string
func NewLatLongISO6709(iso6709 string) *Coordinate {
	re := regexp.MustCompile(`(?P<Latitude>[\+-][\d.]+)(?P<Longitude>[\+-][\d.]+)(?P<Altitude>[\+-][\d.]+)?`)

	if re.MatchString(iso6709) {
		match := re.FindStringSubmatch(iso6709)

		var latitude, longitude, latprec, longprec float64
		var altitude *float64

		for i, name := range re.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}

			switch name {
			case "Latitude":
				latitude, latprec = getLat(match[i])
			case "Longitude":
				longitude, longprec = getLng(match[i])
			case "Altitude":
				altitude = getAlt(match[i])
			}
		}
		return NewLatLongAlt(latitude, longitude, latprec, longprec, altitude)
	}
	return nil
}

// Scan is for fmt.Scanner
func (latlong *Coordinate) Scan(state fmt.ScanState, verb rune) (err error) {
	var token []byte
	token, err = state.Token(false, nil)
	if err == nil {
		*latlong = *NewLatLongISO6709(string(token))
	}
	return
}

var msgCatalog = map[string]struct {
	latN   string
	latS   string
	lngE   string
	lngW   string
	elv    string
	ground string
	dep    string
	comma  string
}{
	"ja": {
		latN:   "北緯%s度",
		latS:   "南緯%s度",
		lngE:   "東経%s度",
		lngW:   "西経%s度",
		elv:    "標高%.0fm",
		ground: "ごく浅く",
		dep:    "深さ%.0fkm",
		comma:  "、",
	},
	"en": {
		latN:   "lat.%sN",
		latS:   "lat.%sS",
		lngE:   "long.%sE",
		lngW:   "long.%sW",
		elv:    "elv.%.0fm",
		ground: "shallow ground",
		dep:    "dep.%.0fkm",
		comma:  ", ",
	},
}

// AltString is string getter for altitude
func (latlong Coordinate) altString() string {
	if latlong.alt != nil {
		return strconv.FormatFloat(*latlong.alt, 'f', 0, 64)
	}
	return ""
}

func (latlong Coordinate) String() string {
	var ss []string
	ss = append(ss, latlong.LatString())
	ss = append(ss, latlong.LngString())
	if latlong.alt != nil {
		if *latlong.alt > 0 {
			ss = append(ss, fmt.Sprintf(msgCatalog[Config.Lang].elv, *latlong.alt))
		} else if *latlong.alt > -10000 {
			ss = append(ss, fmt.Sprintf(msgCatalog[Config.Lang].ground))
		} else {
			ss = append(ss, fmt.Sprintf(msgCatalog[Config.Lang].dep, *latlong.alt/(-1000)))
		}
	}
	return strings.Join(ss, msgCatalog[Config.Lang].comma)
}

func getAlt(part string) (altitude *float64) {
	part = strings.TrimSpace(part)
	if a, er := strconv.ParseFloat(part, 64); er == nil {
		altitude = &a
	}
	return
}
