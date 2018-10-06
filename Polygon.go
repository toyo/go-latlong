package latlong

import (
	"encoding/json"

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

// CapBound is for s2.Region interface.
func (cds *Polygon) CapBound() s2.Cap {
	return cds.S2Region().CapBound()
}

// RectBound is for s2.Region interface.
func (cds *Polygon) RectBound() s2.Rect {
	return cds.S2Region().RectBound()
}

// ContainsCell is for s2.Region interface.
func (cds *Polygon) ContainsCell(c s2.Cell) bool {
	return cds.S2Region().ContainsCell(c)
}

// IntersectsCell is for s2.Region interface.
func (cds *Polygon) IntersectsCell(c s2.Cell) bool {
	return cds.S2Region().IntersectsCell(c)
}

// ContainsPoint is for s2.Region interface.
func (cds *Polygon) ContainsPoint(p s2.Point) bool {
	return cds.S2Region().ContainsPoint(p)
}

// CellUnionBound is for s2.Region interface.
func (cds *Polygon) CellUnionBound() []s2.CellID {
	return cds.S2Region().CellUnionBound()
}

// Center is Center LatLng
func (cds *Polygon) S2Point() s2.Point {
	return cds.S2Region().Centroid()
}

// NewGeoJSONGeometry returns GeoJSONGeometry.
func (cds Polygon) NewGeoJSONGeometry() *GeoJSONGeometry {
	var g GeoJSONGeometry
	g.Type = "Polygon"
	var err error
	mmp := []MultiPoint{cds.MultiPoint}
	g.Coordinates, err = json.Marshal(&mmp)
	if err != nil {
		panic("Error")
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
