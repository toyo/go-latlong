package latlong

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Coordinate is Rectangle and Altitude.
type Coordinate struct {
	s2.Rect
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

// NewLatLong is from latitude and longitude.
func NewLatLong(latitude, longitude, latprec, longprec float64) *Coordinate {
	latlong := new(Coordinate)
	latlong.Rect = s2.RectFromCenterSize(
		s2.LatLngFromDegrees(latitude, longitude),
		s2.LatLngFromDegrees(latprec, longprec))

	return latlong
}

// NewLatLongAlt is from latitude, longitude and altitude.
func NewLatLongAlt(latitude, longitude, latprec, longprec float64, altitude float64) (latlongalt *Coordinate) {
	latlongalt = NewLatLong(latitude, longitude, latprec, longprec)
	latlongalt.alt = &altitude
	return
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

// Lat is getter for latitude
func (latlong Coordinate) Lat() float64 {
	return latlong.Center().Lat.Degrees()
}

// Lng is getter for longitude
func (latlong Coordinate) Lng() float64 {
	return latlong.Center().Lng.Degrees()
}

// DistanceAngle in radian.
func (latlong Coordinate) DistanceAngle(latlong1 *Coordinate) s1.Angle {
	return latlong.Center().Distance(latlong1.Center())
}

// DistanceEarthKm in km at surface.
func (latlong Coordinate) DistanceEarthKm(latlong1 *Coordinate) float64 {
	return float64(latlong.DistanceAngle(latlong1) / 3.14 * 20037.5)
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
func (latlong Coordinate) LatString() (s string) {
	latprec := int(-math.Log10(latlong.Rect.Size().Lat.Degrees()))
	if latprec < 0 {
		latprec = 0
	}

	lat := latlong.Lat()
	if lat >= 0 {
		s += fmt.Sprintf(msgCatalog[Config.Lang].latN, strconv.FormatFloat(lat, 'f', latprec, 64))
	} else {
		s += fmt.Sprintf(msgCatalog[Config.Lang].latS, strconv.FormatFloat(-lat, 'f', latprec, 64))
	}
	return
}

// latString is string getter for latitude
func (latlong Coordinate) latString() string {
	latprec := int(-math.Log10(latlong.Rect.Size().Lat.Degrees()))
	if latprec < 0 {
		latprec = 0
	}
	return strconv.FormatFloat(latlong.Lat(), 'f', latprec, 64)
}

// LngString is string getter for longitude
func (latlong Coordinate) LngString() (s string) {
	lngprec := int(-math.Log10(latlong.Rect.Size().Lng.Degrees()))
	if lngprec < 0 {
		lngprec = 0
	}

	lng := latlong.Lng()
	if lng >= 0 {
		s += fmt.Sprintf(msgCatalog[Config.Lang].lngE, strconv.FormatFloat(lng, 'f', lngprec, 64))
	} else {
		s += fmt.Sprintf(msgCatalog[Config.Lang].lngW, strconv.FormatFloat(-lng, 'f', lngprec, 64))
	}
	return
}

// lngString is string getter for longitude
func (latlong Coordinate) lngString() string {
	lngprec := int(-math.Log10(latlong.Rect.Size().Lng.Degrees()))
	if lngprec < 0 {
		lngprec = 0
	}

	return strconv.FormatFloat(latlong.Lng(), 'f', lngprec, 64)
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

func (latlong Coordinate) String() string {
	return latlong.LatString() + latlong.LngString() + latlong.AltString()
}

// PrecString is Precision String()
func (latlong Coordinate) PrecString() (s string) {
	if Config.Lang == "ja" {
		s = fmt.Sprintf("緯度誤差%f度、経度誤差%f度", latlong.Size().Lat.Degrees(), latlong.Size().Lng.Degrees())
	} else {
		s = fmt.Sprintf("lat. error %fdeg., long. error %fdeg.", latlong.Size().Lat.Degrees(), latlong.Size().Lng.Degrees())
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
