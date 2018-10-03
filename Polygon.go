package latlong

import (
	"github.com/golang/geo/s2"
)

// Polygon is Polygon.
type Polygon struct {
	MultiPoint
}

// S2Region is getter for s2.Loop.
func (cds Polygon) S2Region() *s2.Loop {
	var ps s2.Polyline
	for _, cd := range cds.MultiPoint {
		ps = append(ps, cd.S2Point())
	}
	l := s2.LoopFromPoints(ps)
	if !l.IsNormalized() {
		l.Invert()
	}
	return l
}

// NewGeoJSONGeometry returns GeoJSONGeometry.
func (cds Polygon) NewGeoJSONGeometry() *GeoJSONGeometry {
	var g GeoJSONGeometry
	g.Type = "Polygon"
	g.Coordinates = make([]interface{}, len(cds.MultiPoint))
	for i := range cds.MultiPoint {
		g.Coordinates[i] = cds.MultiPoint[i]
	}
	return &g
}
