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

// Point is Latitude & Longitude with precision.
type Point struct {
	s2.LatLng
	latprec s1.Angle
	lngprec s1.Angle
	alt     *float64 // altitude
}

// NewLatLongAlt is from latitude, longitude and altitude.
func NewLatLongAlt(latitude, longitude, latprec, longprec float64, altitude *float64) *Point {
	var latlongalt Point
	latlongalt.LatLng = s2.LatLngFromDegrees(latitude, longitude)
	latlongalt.latprec = s1.Angle(latprec) * s1.Degree
	latlongalt.lngprec = s1.Angle(longprec) * s1.Degree
	latlongalt.alt = altitude
	return &latlongalt
}

// NewPointISO6709 is from ISO6709 string
func NewPointISO6709(iso6709 string) *Point {
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

// NewPoint is from latitude, longitude and altitude.
func NewPoint(latitude, longitude, latprec, longprec float64) *Point {
	var latlongalt Point
	latlongalt.LatLng = s2.LatLngFromDegrees(latitude, longitude)
	latlongalt.latprec = s1.Angle(latprec) * s1.Degree
	latlongalt.lngprec = s1.Angle(longprec) * s1.Degree
	return &latlongalt
}

// Equal is true if coordinate is same.
func (latlong *Point) Equal(latlong1 *Point) bool {
	return latlong.Lat == latlong1.Lat && latlong.Lng == latlong1.Lng
}

// Scan is for fmt.Scanner
func (latlong *Point) Scan(state fmt.ScanState, verb rune) (err error) {
	var token []byte
	token, err = state.Token(false, nil)
	if err == nil {
		*latlong = *NewPointISO6709(string(token))
	}
	return
}

// S2LatLng is getter for s2.LatLng
func (latlong Point) S2LatLng() s2.LatLng {
	return latlong.LatLng
}

// S2Point is getter for s2.Point
func (latlong Point) S2Point() s2.Point {
	return s2.PointFromLatLng(latlong.S2LatLng())
}

// DistanceAngle in radian.
func (latlong *Point) DistanceAngle(latlong1 *Point) s1.Angle {
	return latlong.Distance(latlong1.LatLng)
}

// DistanceEarthKm in km at surface.
func (latlong *Point) DistanceEarthKm(latlong1 *Point) Km {
	return EarthArcFromAngle(latlong.DistanceAngle(latlong1))
}

func (latlong Point) latpreclog() (latprec int) {
	if latlong.lngprec.Degrees() != 0 {
		latprec = int(math.Ceil(-math.Log10(latlong.latprec.Degrees())))
		if latprec < 0 {
			latprec = 0
		}
	} else {
		latprec = 2
	}
	return
}

// LatString is string getter for latitude
func (latlong Point) LatString() (s string) {
	lat := latlong.Lat.Degrees()
	if lat >= 0 {
		s += fmt.Sprintf(msgCatalog[Config.Lang].latN, strconv.FormatFloat(lat, 'f', latlong.latpreclog(), 64))
	} else {
		s += fmt.Sprintf(msgCatalog[Config.Lang].latS, strconv.FormatFloat(-lat, 'f', latlong.latpreclog(), 64))
	}
	//s += "精度" + strconv.FormatFloat(latlong.latprec.Degrees(), 'f', 5, 64)
	return
}

// latString is string getter for latitude
func (latlong Point) latString() string {
	return strconv.FormatFloat(latlong.Lat.Degrees(), 'f', latlong.latpreclog(), 64)
}

func (latlong Point) lngpreclog() (lngprec int) {
	if latlong.lngprec.Degrees() != 0 {
		lngprec = int(math.Ceil(-math.Log10(latlong.lngprec.Degrees())))
		if lngprec < 0 {
			lngprec = 0
		}
	} else {
		lngprec = 2
	}
	return
}

// LngString is string getter for longitude
func (latlong Point) LngString() (s string) {
	lng := latlong.Lng.Degrees()
	if lng >= 0 {
		s += fmt.Sprintf(msgCatalog[Config.Lang].lngE, strconv.FormatFloat(lng, 'f', latlong.lngpreclog(), 64))
	} else {
		s += fmt.Sprintf(msgCatalog[Config.Lang].lngW, strconv.FormatFloat(-lng, 'f', latlong.lngpreclog(), 64))
	}
	//s += "精度" + strconv.FormatFloat(latlong.lngprec.Degrees(), 'f', 5, 64)
	return
}

// lngString is string getter for longitude
func (latlong Point) lngString() string {
	return strconv.FormatFloat(latlong.Lng.Degrees(), 'f', latlong.lngpreclog(), 64)
}

func getAlt(part string) (altitude *float64) {
	part = strings.TrimSpace(part)
	if a, er := strconv.ParseFloat(part, 64); er == nil {
		altitude = &a
	}
	return
}

func (latlong Point) String() string {
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

// PrecisionArea returns area size of precicion.
func (latlong Point) PrecisionArea() float64 {
	return latlong.latprec.Degrees() * latlong.lngprec.Degrees()
}

// PrecString is Precision String()
func (latlong Point) PrecString() (s string) {
	if Config.Lang == "ja" {
		s = fmt.Sprintf("緯度誤差%f度、経度誤差%f度", latlong.latprec, latlong.lngprec)
	} else {
		s = fmt.Sprintf("lat. error %fdeg., long. error %fdeg.", latlong.latprec, latlong.lngprec)
	}
	return
}

// MapsLatLng return maps.LatLng ( "googlemaps.github.io/maps" )
func (latlong Point) MapsLatLng() maps.LatLng {
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
		degprec = float64(1) / 60
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

/*
// MarshalJSON is a marshaler for JSON.
func (latlong Point) MarshalJSON() ([]byte, error) {
	s := latlong.lngString() + "," + latlong.latString()
	return []byte("[" + s + "]"), nil
}

// UnmarshalJSON is a unmarshaler for JSON.
func (latlong *Point) UnmarshalJSON(data []byte) (err error) {
	data = bytes.TrimSpace(data)
	data = bytes.TrimLeft(data, "[")
	data = bytes.TrimRight(data, "]")
	datas := bytes.Split(data, []byte(`,`))

	if len(data) < 2 {
		return fmt.Errorf("Not enough LatLng JSON %d %s", len(data), string(data))
	}

	lat, latprec := getLat(string(datas[1]))
	lng, lngprec := getLat(string(datas[0]))
	*latlong = *NewPoint(lat, lng, latprec, lngprec)
	return
}
*/

// MarshalJSON is a marshaler for JSON.
func (latlong *Point) MarshalJSON() ([]byte, error) {
	s := latlong.lngString() + "," + latlong.latString()
	if latlong.alt != nil {
		s += "," + latlong.altString()
	}
	return []byte("[" + s + "]"), nil
}

// UnmarshalJSON is a unmarshaler for JSON.
func (latlong *Point) UnmarshalJSON(data []byte) (err error) {
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
func (latlong Point) altString() string {
	if latlong.alt != nil {
		return strconv.FormatFloat(*latlong.alt, 'f', 0, 64)
	}
	return ""
}
