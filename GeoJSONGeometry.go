package latlong

import (
	"encoding/json"

	"github.com/golang/geo/s2"
)

// GeoJSONGeometry is Geometry of GeoJSON
type GeoJSONGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
	Radius      *float64        `json:"radius,omitempty"` // only for Circle, which is GeoJSON specification 1.1 and leter.
}

// Equal return equal or not.
func (geom GeoJSONGeometry) Equal(geom1 GeoJSONGeometry) bool {
	if geom.Type != geom1.Type {
		return false
	}
	if geom.Radius != geom1.Radius {
		return false
	}
	for i := range geom.Coordinates {
		if geom.Coordinates[i] != geom1.Coordinates[i] {
			return false
		}
	}
	return true
}

// S2Region is getter for s2.Region.
func (geom GeoJSONGeometry) S2Region() s2.Region {
	switch geom.Type {
	case "Circle":
		return geom.Circle().S2Region()
	case "LineString":
		return geom.LineString().S2Region()
	case "Polygon":
		return geom.Polygon().S2Region()
	}
	return s2.EmptyRect()
}

// S2Point is getter for center of s2.Point.
func (geom GeoJSONGeometry) S2Point() s2.Point {
	switch geom.Type {
	case "Point":
		c := geom.Point().S2Point()
		return c
	case "Circle":
		c := geom.Circle().S2Point()
		return c
	case "LineString":
		c := geom.LineString().S2Point()
		return c
	case "Polygon":
		c := geom.Polygon().S2Point()
		return c
	}
	panic(geom.Type)
}

// S2LatLng returns s2.LatLng
func (geom GeoJSONGeometry) S2LatLng() s2.LatLng {
	return s2.LatLngFromPoint(geom.S2Point())
}

// Polygon extract Polygon
func (geom GeoJSONGeometry) Polygon() *Polygon {
	if geom.Type != "Polygon" {
		return nil
	}
	var co []MultiPoint
	err := json.Unmarshal(geom.Coordinates, &co)
	if err != nil {
		panic("Error")
	}

	switch len(co) {
	case 0:
		panic("No Polygon!")
	case 1:
		return &Polygon{MultiPoint: co[0]}
	default:
		panic("Polygon has hole! Not implemented")
	}
}

// LineString extract LineString
func (geom GeoJSONGeometry) LineString() *LineString {
	if geom.Type != "LineString" {
		return nil
	}
	var coor MultiPoint
	err := json.Unmarshal(geom.Coordinates, &coor)
	if err != nil {
		panic("Error")
	}

	return &LineString{MultiPoint: coor}
}

// Circle extract Circle
func (geom GeoJSONGeometry) Circle() (ls *Circle) {
	if geom.Type != "Circle" {
		return nil
	}
	var coor Point
	err := json.Unmarshal(geom.Coordinates, &coor)
	if err != nil {
		panic("Error")
	}

	radius := geom.Radius
	if radius == nil {
		ls = NewEmptyCircle()
	} else {
		ls = NewCircle(coor, Km(*radius))
	}

	return
}

// Point extract Point
func (geom GeoJSONGeometry) Point() *Point {
	if geom.Type != "Point" {
		return nil
	}
	var coor Point
	err := json.Unmarshal(geom.Coordinates, &coor)
	if err != nil {
		panic("Error")
	}

	return &coor
}
