package latlong

import (
	"github.com/golang/geo/s2"
)

// Polygon inherited MultiPoint
type Polygon struct {
	MultiPoint
}

// S2Region is getter for s2.Loop.
func (cds *Polygon) S2Region() *s2.Loop {
	ps := make(s2.Polyline, len(cds.MultiPoint))
	for i := range cds.MultiPoint {
		ps[i] = cds.MultiPoint[i].S2Point()
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

// NewGeoJSONFeature returns GeoJSONFeature.
func (cds Polygon) NewGeoJSONFeature(property interface{}) *GeoJSONFeature {
	var g GeoJSONFeature
	g.Type = "Feature"
	g.Geometry = cds.NewGeoJSONGeometry()
	g.Property = property
	return &g
}
