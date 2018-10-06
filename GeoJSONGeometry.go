package latlong

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/golang/geo/s2"
)

// Geometry is interface for each geometry class @ GeoJSON.
type Geometry interface {
	Equal(Geometry) bool
	S2Region() s2.Region
	S2Point() s2.Point
	Radiusp() *float64
	Type() string
	String() string
}

// GeoJSONGeometry is Geometry of GeoJSON
/*
type GeoJSONGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
	Radius      *float64        `json:"radius,omitempty"` // only for Circle, which is GeoJSON specification 1.1 and leter.
}
*/
type GeoJSONGeometry struct {
	geo Geometry
}

// NewGeoJSONGeometry is constructor
func NewGeoJSONGeometry(g Geometry) (geom GeoJSONGeometry) {
	geom.geo = g
	return
}

// Geo returns contents.
func (geom GeoJSONGeometry) Geo() Geometry {
	return geom.geo
}

// Equal return equal or not.
func (geom GeoJSONGeometry) Equal(geom1 GeoJSONGeometry) bool {
	return geom.geo.Equal(geom1.geo)
}

// S2Region is getter for s2.Region.
func (geom GeoJSONGeometry) S2Region() s2.Region {
	return geom.geo.S2Region()
}

// S2Point is getter for center of s2.Point.
func (geom GeoJSONGeometry) S2Point() s2.Point {
	return geom.geo.S2Point()
}

// S2LatLng returns s2.LatLng
func (geom GeoJSONGeometry) S2LatLng() s2.LatLng {
	return s2.LatLngFromPoint(geom.S2Point())
}

func (geom GeoJSONGeometry) String() string {
	return geom.geo.String()
}

// MarshalJSON is a marshaler for JSON.
func (geom GeoJSONGeometry) MarshalJSON() ([]byte, error) {
	var js struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
		Radius      *float64        `json:"radius,omitempty"` // only for Circle, which is GeoJSON specification 1.1 and leter.
	}

	var err error
	if geom.geo != nil {
		js.Type = geom.geo.Type()
	} else {
		js.Type = "Null"
	}
	js.Coordinates, err = json.Marshal(geom.geo)
	if err != nil {
		return nil, err
	}
	if geom.geo != nil {
		js.Radius = geom.geo.Radiusp()
	}
	return json.Marshal(&js)
}

// UnmarshalJSON is a unmarshaler for JSON.
func (geom *GeoJSONGeometry) UnmarshalJSON(data []byte) error {
	var js struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
		Radius      *float64        `json:"radius,omitempty"` // only for Circle, which is GeoJSON specification 1.1 and leter.
	}

	err := json.Unmarshal(bytes.TrimSpace(data), &js)
	if err != nil {
		return err
	}
	switch js.Type {
	case "Polygon":
		var p Polygon
		err := json.Unmarshal(js.Coordinates, &p)
		geom.geo = p
		return err
	case "LineString":
		var p LineString
		err := json.Unmarshal(js.Coordinates, &p)
		geom.geo = p
		return err
	case "Point":
		var p Point
		err := json.Unmarshal(js.Coordinates, &p)
		geom.geo = p
		return err
	case "Circle":
		var p Point
		err := json.Unmarshal(js.Coordinates, &p)
		radius := js.Radius
		if radius == nil {
			geom.geo = *NewEmptyCircle()
		} else {
			geom.geo = *NewCircle(p, Km(*radius))
		}
		return err
	case "Null":
		geom.geo = nil
		return nil
	}
	return fmt.Errorf("Unknown type %s", js.Type)
}
