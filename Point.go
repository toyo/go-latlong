package latlong

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	lat Angle
	lng Angle
	alt *float64 // altitude
}

// Type returns this type
func (Point) Type() string {
	return "Point"
}

// NewLatLongAlt is from latitude, longitude and altitude.
func NewLatLongAlt(lat, lng Angle, altitude *float64) *Point {
	var latlongalt Point
	latlongalt.lat = lat
	latlongalt.lng = lng
	latlongalt.alt = altitude
	return &latlongalt
}

// NewPointISO6709 is from ISO6709 string
func NewPointISO6709(iso6709 []byte) *Point {
	re := regexp.MustCompile(`(?P<Latitude>[\+-][\d.]+)(?P<Longitude>[\+-][\d.]+)(?P<Altitude>[\+-][\d.]+)?`)

	if re.Match(iso6709) {
		match := re.FindSubmatch(iso6709)

		var lat, lng Angle
		var altitude *float64

		for i, name := range re.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}

			switch name {
			case "Latitude":
				lat = AngleFromBytes(match[i])
			case "Longitude":
				lng = AngleFromBytes(match[i])
			case "Altitude":
				altitude = getAlt(match[i])
			}
		}
		return NewLatLongAlt(lat, lng, altitude)
	}
	return nil
}

// NewPointFromS2Point is from s2.Point
func NewPointFromS2Point(p s2.Point) Point {
	s2ll := s2.LatLngFromPoint(p)
	return Point{lat: NewAngleFromS1Angle(s2ll.Lat, 0), lng: NewAngleFromS1Angle(s2ll.Lng, 0)}
}

// Lat is getter for latitude.
func (latlong *Point) Lat() Angle {
	return latlong.lat
}

// Lng is getter for longitude.
func (latlong *Point) Lng() Angle {
	return latlong.lng
}

// Equal is true if coordinate is same.
func (latlong Point) Equal(latlong1 Geometry) bool {
	return latlong == latlong1.(Point)
}

// Scan is for fmt.Scanner
func (latlong *Point) Scan(state fmt.ScanState, verb rune) (err error) {
	var token []byte
	token, err = state.Token(false, nil)
	if err == nil {
		*latlong = *NewPointISO6709(token)
	}
	return
}

// S2LatLng is getter for s2.LatLng
func (latlong Point) S2LatLng() s2.LatLng {
	return s2.LatLng{Lat: latlong.Lat().S1Angle(), Lng: latlong.Lng().S1Angle()}
}

// S2Point is getter for s2.Point
func (latlong Point) S2Point() s2.Point {
	return s2.PointFromLatLng(latlong.S2LatLng())
}

// S2Region is getter for s2.Loop.
func (latlong Point) S2Region() s2.Region {
	return s2.PointFromLatLng(latlong.S2LatLng())
}

// Radiusp is un-used
func (latlong Point) Radiusp() *float64 {
	return nil
}

// DistanceAngle in radian.
func (latlong *Point) DistanceAngle(latlong1 *Point) s1.Angle {
	return latlong.S2LatLng().Distance(latlong1.S2LatLng())
}

// DistanceEarthKm in km at surface.
func (latlong *Point) DistanceEarthKm(latlong1 *Point) Km {
	return EarthArcFromAngle(latlong.DistanceAngle(latlong1))
}

// LatString is string getter for latitude
func (latlong Point) LatString() (s string) {
	lat := latlong.Lat().Degrees()
	if lat >= 0 {
		s += fmt.Sprintf(msgCatalog[Config.Lang].latN, strconv.FormatFloat(lat, 'f', latlong.lat.preclog(), 64))
	} else {
		s += fmt.Sprintf(msgCatalog[Config.Lang].latS, strconv.FormatFloat(-lat, 'f', latlong.lat.preclog(), 64))
	}
	//s += "精度" + strconv.FormatFloat(latlong.latprec.Degrees(), 'f', 5, 64)
	return
}

// latString is string getter for latitude
func (latlong Point) latString() string {
	return strconv.FormatFloat(latlong.Lat().Degrees(), 'f', latlong.lat.preclog(), 64)
}

// LngString is string getter for longitude
func (latlong Point) LngString() (s string) {
	lng := latlong.Lng().Degrees()
	if lng >= 0 {
		s += fmt.Sprintf(msgCatalog[Config.Lang].lngE, strconv.FormatFloat(lng, 'f', latlong.lng.preclog(), 64))
	} else {
		s += fmt.Sprintf(msgCatalog[Config.Lang].lngW, strconv.FormatFloat(-lng, 'f', latlong.lng.preclog(), 64))
	}
	//s += "精度" + strconv.FormatFloat(latlong.lngprec.Degrees(), 'f', 5, 64)
	return
}

// lngString is string getter for longitude
func (latlong Point) lngString() string {
	return strconv.FormatFloat(latlong.Lng().Degrees(), 'f', latlong.lng.preclog(), 64)
}

func getAlt(part []byte) (altitude *float64) {
	part = bytes.TrimSpace(part)
	if a, er := strconv.ParseFloat(string(part), 64); er == nil {
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
	return latlong.lat.PrecDegrees() * latlong.lng.PrecDegrees()
}

// PrecString is Precision String()
func (latlong Point) PrecString() (s string) {
	if Config.Lang == "ja" {
		s = fmt.Sprintf("緯度誤差%f度、経度誤差%f度", latlong.lat.PrecDegrees(), latlong.lng.PrecDegrees())
	} else {
		s = fmt.Sprintf("lat. error %fdeg., long. error %fdeg.", latlong.lat.PrecDegrees(), latlong.lng.PrecDegrees())
	}
	return
}

// MapsLatLng return maps.LatLng ( "googlemaps.github.io/maps" )
func (latlong Point) MapsLatLng() maps.LatLng {
	return maps.LatLng{Lat: latlong.Lat().Degrees(), Lng: latlong.Lng().Degrees()}
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
	Lang:          "en", // or "ja" 	// Lang is an string language
	YahooJPAPIURL: "https://map.yahooapis.jp/geoapi/V1/reverseGeoCoder",
}

// MarshalJSON is a marshaler for JSON.
func (latlong Point) MarshalJSON() ([]byte, error) {
	var ll []Angle

	if latlong.alt != nil {
		ll = make([]Angle, 3)
		ll[2].radian = s1.Angle(*latlong.alt) * s1.Degree
		ll[2].radianprec = 1
	} else {
		ll = make([]Angle, 2)
	}

	ll[0] = latlong.lng
	ll[1] = latlong.lat

	return json.Marshal(&ll)
}

// UnmarshalJSON is a unmarshaler for JSON.
func (latlong *Point) UnmarshalJSON(data []byte) (err error) {
	var ll []Angle

	err = json.Unmarshal(bytes.TrimSpace(data), &ll)

	if len(ll) < 2 {
		return errors.New("unknown JSON Coordinate format")
	}

	latlong.lng = ll[0]
	latlong.lat = ll[1]

	if len(ll) > 2 {
		altitude := ll[2].radian.Degrees()
		latlong.alt = &altitude
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

// NewGeoJSONGeometry returns GeoJSONGeometry.
func (latlong Point) NewGeoJSONGeometry() GeoJSONGeometry {
	var g GeoJSONGeometry
	g.geo = latlong

	return g
}

// NewGeoJSONFeature returns GeoJSONFeature.
func (latlong Point) NewGeoJSONFeature(property interface{}) *GeoJSONFeature {
	var g GeoJSONFeature
	g.Type = "Feature"
	g.Geometry = latlong.NewGeoJSONGeometry()
	g.Property = property
	return &g
}
