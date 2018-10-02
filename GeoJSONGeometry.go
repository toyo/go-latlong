package latlong

import (
	"errors"
)

// GeoJSONGeometry is Geometry of GeoJSON
type GeoJSONGeometry struct {
	Type        string        `json:"type"`
	Coordinates []interface{} `json:"coordinates"`
	Radius      *float64      `json:"radius"` // only for Circle, which is GeoJSON specification 1.1 and leter.
}

// Polygon extract Polygon
func (geom GeoJSONGeometry) Polygon() (ls Polygon, err error) {
	if geom.Type != "Polygon" {
		err = errors.New("No Polygon")
		return
	}
	ls = make(Polygon, len(geom.Coordinates[0].([]interface{})))
	for i := range geom.Coordinates[0].([]interface{}) {
		g := geom.Coordinates[0].([]interface{})[i].([]interface{})
		ls[i] = NewLatLng(g[1].(float64), g[0].(float64), 0, 0)
	}
	return
}

// LineString extract LineString
func (geom GeoJSONGeometry) LineString() (ls LineString, err error) {
	if geom.Type != "LineString" {
		err = errors.New("No LineString")
		return
	}
	ls = make(LineString, len(geom.Coordinates))
	for i := range geom.Coordinates {
		g := geom.Coordinates[i].([]interface{})
		ls[i] = NewLatLng(g[1].(float64), g[0].(float64), 0, 0)
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
		ls = *NewCircle(*NewLatLng(geom.Coordinates[1].(float64), geom.Coordinates[0].(float64), 0, 0), Km(*radius))
	}

	return
}
