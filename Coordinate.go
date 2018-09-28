package latlong

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"googlemaps.github.io/maps"
)

// LatLng is Latitude & Longitude with precision.
type LatLng struct {
	s2.LatLng
	latprec float64
	lngprec float64
}

// Coordinate is LatLng and Altitude.
type Coordinate struct {
	LatLng
	alt *float64 // altitude
}

// MarshalJSON is a marshaler for JSON.
func (latlong *LatLng) MarshalJSON() ([]byte, error) {
	s := latlong.lngString() + "," + latlong.latString()
	return []byte("[" + s + "]"), nil
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
	latlongalt := new(Coordinate)
	latlongalt.LatLng.LatLng = s2.LatLngFromDegrees(latitude, longitude)
	latlongalt.latprec = latprec
	latlongalt.lngprec = longprec
	latlongalt.alt = altitude
	return latlongalt
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

// S2LatLng is getter for s2.LatLng
func (latlong LatLng) S2LatLng() s2.LatLng {
	return latlong.LatLng
}

// S2Point is getter for s2.Point
func (latlong LatLng) S2Point() s2.Point {
	return s2.PointFromLatLng(latlong.S2LatLng())
}

// S2Cap is getter for s2.Cap
func (latlong *LatLng) S2Cap(radius s1.Angle) s2.Cap {
	return s2.CapFromCenterAngle(latlong.S2Point(), radius)
}

// DistanceAngle in radian.
func (latlong *LatLng) DistanceAngle(latlong1 *LatLng) s1.Angle {
	return latlong.Distance(latlong1.LatLng)
}

// DistanceEarthKm in km at surface.
func (latlong *LatLng) DistanceEarthKm(latlong1 *LatLng) Km {
	return EarthArcFromAngle(latlong.DistanceAngle(latlong1))
}

var msgCatalog = map[string]struct {
	latN   string
	latS   string
	lngE   string
	lngW   string
	elv    string
	ground string
	dep    string
}{
	"ja": {
		latN:   "北緯%s度、",
		latS:   "南緯%s度、",
		lngE:   "東経%s度",
		lngW:   "西経%s度",
		elv:    "、標高%.0fm",
		ground: "、ごく浅く",
		dep:    "、深さ%.0fkm",
	},
	"en": {
		latN:   "lat.%sN, ",
		latS:   "lat.%sS, ",
		lngE:   "long.%sE",
		lngW:   "long.%sW",
		elv:    ", elv.%.0fm",
		ground: ", shallow ground",
		dep:    ", dep.%.0fkm",
	},
}

// LatString is string getter for latitude
func (latlong LatLng) LatString() (s string) {
	latprec := int(-math.Log10(latlong.latprec))
	if latprec < 0 {
		latprec = 0
	}

	lat := latlong.Lat.Degrees()
	if lat >= 0 {
		s += fmt.Sprintf(msgCatalog[Config.Lang].latN, strconv.FormatFloat(lat, 'f', latprec, 64))
	} else {
		s += fmt.Sprintf(msgCatalog[Config.Lang].latS, strconv.FormatFloat(-lat, 'f', latprec, 64))
	}
	return
}

// latString is string getter for latitude
func (latlong LatLng) latString() string {
	latprec := int(-math.Log10(latlong.latprec))
	if latprec < 0 {
		latprec = 0
	}
	return strconv.FormatFloat(latlong.Lat.Degrees(), 'f', latprec, 64)
}

// LngString is string getter for longitude
func (latlong LatLng) LngString() (s string) {
	lngprec := int(-math.Log10(latlong.lngprec))
	if lngprec < 0 {
		lngprec = 0
	}

	lng := latlong.Lng.Degrees()
	if lng >= 0 {
		s += fmt.Sprintf(msgCatalog[Config.Lang].lngE, strconv.FormatFloat(lng, 'f', lngprec, 64))
	} else {
		s += fmt.Sprintf(msgCatalog[Config.Lang].lngW, strconv.FormatFloat(-lng, 'f', lngprec, 64))
	}
	return
}

// lngString is string getter for longitude
func (latlong LatLng) lngString() string {
	lngprec := int(-math.Log10(latlong.lngprec))
	if lngprec < 0 {
		lngprec = 0
	}

	return strconv.FormatFloat(latlong.Lng.Degrees(), 'f', lngprec, 64)
}

// AltString is string getter for altitude
func (latlong Coordinate) AltString() (s string) {
	if latlong.alt != nil {
		if *latlong.alt > 0 {
			s += fmt.Sprintf(msgCatalog[Config.Lang].elv, *latlong.alt)
		} else if *latlong.alt > -10000 {
			s += fmt.Sprintf(msgCatalog[Config.Lang].ground)
		} else {
			s += fmt.Sprintf(msgCatalog[Config.Lang].dep, *latlong.alt/(-1000))
		}
	}
	return
}

// AltString is string getter for altitude
func (latlong Coordinate) altString() string {
	if latlong.alt != nil {
		return strconv.FormatFloat(*latlong.alt, 'f', 0, 64)
	}
	return ""
}

func (latlong LatLng) String() string {
	return latlong.LatString() + latlong.LngString()
}

func (latlong Coordinate) String() string {
	return latlong.LatString() + latlong.LngString() + latlong.AltString()
}

// PrecisionArea returns area size of precicion.
func (latlong LatLng) PrecisionArea() float64 {
	return latlong.latprec * latlong.lngprec
}

// PrecString is Precision String()
func (latlong LatLng) PrecString() (s string) {
	if Config.Lang == "ja" {
		s = fmt.Sprintf("緯度誤差%f度、経度誤差%f度", latlong.latprec, latlong.lngprec)
	} else {
		s = fmt.Sprintf("lat. error %fdeg., long. error %fdeg.", latlong.latprec, latlong.lngprec)
	}
	return
}

// MapsLatLng return maps.LatLng ( "googlemaps.github.io/maps" )
func (latlong LatLng) MapsLatLng() maps.LatLng {
	return maps.LatLng{Lat: latlong.Lat.Degrees(), Lng: latlong.Lng.Degrees()}
}

func isErrorDeg(deg float64, degprec float64) bool {
	degerr, degprecerr := getErrorDeg()
	if deg == degerr && degprec == degprecerr {
		return true
	}
	return false
}

func getErrorDeg() (deg float64, degprec float64) {
	deg = 0
	degprec = 360
	return
}

func getDeg(part string, pos int) (deg float64, degprec float64) {
	var err error
	deg, err = strconv.ParseFloat(part, 64)
	if err != nil {
		deg, degprec = getErrorDeg()
		return
	}

	if l := len(part); l == pos {
		degprec = 1
	} else {
		degprec = math.Pow10(pos - l + 1)
	}
	return
}

func getDegMin(part string, pos int) (deg float64, degprec float64) {
	var err error
	if deg, err = strconv.ParseFloat(part[1:pos-2], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}

	var min float64
	if min, err = strconv.ParseFloat(part[pos-2:], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}
	deg += min / 60

	switch part[0] {
	case '-':
		deg = -deg
	case '+':
		break
	default:
		deg, degprec = getErrorDeg()
		return
	}

	if l := len(part); l == pos {
		degprec = 1 / 60
	} else {
		degprec = math.Pow10(pos-l+1) / 60
	}

	return
}

func getDegMinSec(part string, pos int) (deg float64, degprec float64) {
	var err error
	if deg, err = strconv.ParseFloat(part[1:pos-4], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}

	var min float64
	if min, err = strconv.ParseFloat(part[pos-4:pos-2], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}
	deg += min / 60

	var sec float64
	if sec, err = strconv.ParseFloat(part[pos-2:], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}
	deg += sec / 3600

	switch part[0] {
	case '-':
		deg = -deg
	case '+':
		break
	default:
		deg, degprec = getErrorDeg()
		return
	}

	if l := len(part); l == pos {
		degprec = float64(1) / 3600
	} else {
		degprec = math.Pow10(pos-l+1) / 3600
	}

	return
}

func getLat(part string) (latitude float64, latprec float64) {
	part = strings.TrimSpace(part)
	pos := strings.Index(part, ".")
	if pos == -1 {
		pos = len(part)
	}

	if pos < 2 && false {
		latitude, latprec = getErrorDeg()
	} else if pos < 4 {
		latitude, latprec = getDeg(part, pos)
	} else if pos < 6 {
		latitude, latprec = getDegMin(part, pos)
	} else if pos < 8 {
		latitude, latprec = getDegMinSec(part, pos)
	} else {
		latitude, latprec = getErrorDeg()
	}
	return
}

func getLng(part string) (longitude float64, longprec float64) {
	part = strings.TrimSpace(part)
	pos := strings.Index(part, ".")
	if pos == -1 {
		pos = len(part)
	}

	if pos < 3 && false {
		longitude, longprec = getErrorDeg()
	} else if pos < 5 {
		longitude, longprec = getDeg(part, pos)
	} else if pos < 7 {
		longitude, longprec = getDegMin(part, pos)
	} else if pos < 9 {
		longitude, longprec = getDegMinSec(part, pos)
	} else {
		longitude, longprec = getErrorDeg()
	}
	return
}

func getAlt(part string) (altitude *float64) {
	part = strings.TrimSpace(part)
	if a, er := strconv.ParseFloat(part, 64); er == nil {
		altitude = &a
	}
	return
}

type config struct {
	HTTPClient       *http.Client
	GoogleAPIKey     string
	GoogleMapsAPIURL string
	YahooJPClientID  string
	YahooJPAPIURL    string
	Lang             string
}

// Config is an configuration of environment.
var Config = config{
	HTTPClient:    &http.Client{},
	Lang:          "en", //= "ja" 	// Lang is an string language
	YahooJPAPIURL: "https://map.yahooapis.jp/geoapi/V1/reverseGeoCoder",
}
