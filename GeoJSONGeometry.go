package latlong

import (
	"errors"
)

// GeoJSONGeometry is Geometry of GeoJSON
type GeoJSONGeometry struct {
	Type        string        `json:"type"`
	Coordinates []interface{} `json:"coordinates"`
	Radius      *float64      `json:"radius,omitempty"` // only for Circle, which is GeoJSON specification 1.1 and leter.
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

// Polygon extract Polygon
func (geom GeoJSONGeometry) Polygon() (ls Polygon, err error) {
	if geom.Type != "Polygon" {
		err = errors.New("No Polygon")
		return
	}
	mp := make(MultiPoint, len(geom.Coordinates[0].([]interface{})))
	ls = Polygon{MultiPoint: mp}
	for i := range geom.Coordinates[0].([]interface{}) {
		g := geom.Coordinates[0].([]interface{})[i].([]interface{})
		ls.MultiPoint[i] = NewPoint(g[1].(float64), g[0].(float64), 0, 0)
	}
	return
}

// LineString extract LineString
func (geom GeoJSONGeometry) LineString() (ls LineString, err error) {
	if geom.Type != "LineString" {
		err = errors.New("No LineString")
		return
	}
	mp := make(MultiPoint, len(geom.Coordinates))
	ls = LineString{MultiPoint: mp}
	for i := range geom.Coordinates {
		g := geom.Coordinates[i].([]interface{})
		ls.MultiPoint[i] = NewPoint(g[1].(float64), g[0].(float64), 0, 0)
	}
	return
}

// Circle extract Circle
func (geom GeoJSONGeometry) Circle() (ls Circle, err error) {
	if geom.Type != "Circle" {
		err = errors.New("No Circle")
		return
	}

	radius := geom.Radius
	if radius == nil {
		ls = *NewEmptyCircle()
	} else {
		ls = *NewCircle(*NewPoint(geom.Coordinates[1].(float64), geom.Coordinates[0].(float64), 0, 0), Km(*radius))
	}

	return
}
