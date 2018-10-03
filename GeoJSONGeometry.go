package latlong

import (
	"github.com/golang/geo/s2"
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

// Polygon extract Polygon
func (geom GeoJSONGeometry) Polygon() *Polygon {
	if geom.Type != "Polygon" {
		return nil
	}
	switch len(geom.Coordinates) {
	case 0:
		panic("No Polygon!")
	case 1:
		coor := geom.Coordinates[0].([]interface{})
		mp := make(MultiPoint, len(coor))
		for i := range coor {
			g := coor[i].([]interface{})
			mp[i] = NewPoint(g[1].(float64), g[0].(float64), 0, 0)
		}
		return &Polygon{MultiPoint: mp}
	default:
		panic("Polygon has hole! Not implemented")
	}

}

// LineString extract LineString
func (geom GeoJSONGeometry) LineString() *LineString {
	if geom.Type != "LineString" {
		return nil
	}
	coor := geom.Coordinates
	mp := make(MultiPoint, len(coor))
	for i := range coor {
		g := coor[i].([]interface{})
		mp[i] = NewPoint(g[1].(float64), g[0].(float64), 0, 0)
	}
	return &LineString{MultiPoint: mp}
}

// Circle extract Circle
func (geom GeoJSONGeometry) Circle() (ls *Circle) {
	if geom.Type != "Circle" {
		return nil
	}

	radius := geom.Radius
	if radius == nil {
		ls = NewEmptyCircle()
	} else {
		ls = NewCircle(*NewPoint(geom.Coordinates[1].(float64), geom.Coordinates[0].(float64), 0, 0), Km(*radius))
	}

	return
}
