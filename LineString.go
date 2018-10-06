package latlong

import (
	"encoding/json"

	"github.com/golang/geo/s2"
	"googlemaps.github.io/maps"
)

// LineString inherited MultiPoint
type LineString struct {
	MultiPoint
}

// Type returns this type
func (LineString) Type() string {
	return "LineString"
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

// S2Point is Center
func (cds LineString) S2Point() s2.Point {
	return cds.MultiPoint[0].S2Point()
}

// MapsLatLng convert to google maps.
func (cds LineString) MapsLatLng() (mlls []maps.LatLng) {
	for _, cd := range cds.MultiPoint {
		mlls = append(mlls, cd.MapsLatLng())
	}
	return
}

// S2Region is getter for s2.Polyline ([]s2.Point).
func (cds LineString) S2Region() s2.Region {
	ps := make(s2.Polyline, len(cds.MultiPoint))
	for i := range cds.MultiPoint {
		ps[i] = cds.MultiPoint[i].S2Point()
	}
	return &ps
}

// Radiusp is un-used
func (cds LineString) Radiusp() *float64 {
	return nil
}

// CapBound is for s2.Region interface.
func (cds *LineString) CapBound() s2.Cap {
	return cds.S2Region().CapBound()
}

// RectBound is for s2.Region interface.
func (cds *LineString) RectBound() s2.Rect {
	return cds.S2Region().RectBound()
}

// ContainsCell is for s2.Region interface.
func (cds *LineString) ContainsCell(c s2.Cell) bool {
	return cds.S2Region().ContainsCell(c)
}

// IntersectsCell is for s2.Region interface.
func (cds *LineString) IntersectsCell(c s2.Cell) bool {
	return cds.S2Region().IntersectsCell(c)
}

// ContainsPoint is for s2.Region interface.
func (cds *LineString) ContainsPoint(p s2.Point) bool {
	return cds.S2Region().ContainsPoint(p)
}

// CellUnionBound is for s2.Region interface.
func (cds *LineString) CellUnionBound() []s2.CellID {
	return cds.S2Region().CellUnionBound()
}

// MarshalJSON is a marshaler for JSON.
func (cds LineString) MarshalJSON() ([]byte, error) {
	return json.Marshal(&cds.MultiPoint)
}

// UnmarshalJSON is a unmarshaler for JSON.
func (cds *LineString) UnmarshalJSON(data []byte) (err error) {
	err = json.Unmarshal(data, &cds.MultiPoint)
	if err != nil {
		panic(err)
	}
	return nil
}

// NewGeoJSONGeometry returns GeoJSONGeometry.
func (cds LineString) NewGeoJSONGeometry() GeoJSONGeometry {
	var g GeoJSONGeometry
	g.geo = cds
	return g
}

// NewGeoJSONFeature returns GeoJSONFeature.
func (cds LineString) NewGeoJSONFeature(property interface{}) *GeoJSONFeature {
	var g GeoJSONFeature
	g.Type = "Feature"
	g.Geometry = cds.NewGeoJSONGeometry()
	g.Property = property
	return &g
}
