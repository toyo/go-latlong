package latlong

import (
	"github.com/golang/geo/s2"
	"googlemaps.github.io/maps"
)

// LineString is slice of LatLong
type LineString struct {
	MultiPoint
}

// S2Polyline is getter for s2.Polyline ([]s2.Point).
func (cds LineString) S2Polyline() s2.Polyline {
	var ps s2.Polyline
	for _, cd := range cds.MultiPoint {
		ps = append(ps, cd.S2Point())
	}
	return ps
}

// S2Loop is getter for s2.Loop.
func (cds LineString) S2Loop() *s2.Loop {
	lo := s2.LoopFromPoints(cds.S2Polyline())
	lo.Normalize() // if loop is not CCW but CW, change to CCW.
	return lo
}

// MapsLatLng convert to google maps.
func (cds LineString) MapsLatLng() (mlls []maps.LatLng) {
	for _, cd := range cds.MultiPoint {
		mlls = append(mlls, cd.MapsLatLng())
	}
	return
}

// S2Region is getter for s2.Polyline ([]s2.Point).
func (cds *LineString) S2Region() *s2.Polyline {
	ps := make(s2.Polyline, len(cds.MultiPoint))
	for i := range cds.MultiPoint {
		ps[i] = cds.MultiPoint[i].S2Point()
	}
	return &ps
}

// NewGeoJSONGeometry returns GeoJSONGeometry.
func (cds LineString) NewGeoJSONGeometry() *GeoJSONGeometry {
	var g GeoJSONGeometry
	g.Type = "LineString"
	g.Coordinates = make([]interface{}, len(cds.MultiPoint))
	for i := range cds.MultiPoint {
		g.Coordinates[i] = cds.MultiPoint[i]
	}
	return &g
}
