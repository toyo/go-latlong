package latlong

import (
	"encoding/json"

	"github.com/golang/geo/s2"
)

// Polygon inherited MultiPoint
type Polygon struct {
	LineString
}

/*
// UnmarshalText is from ISO6709 latlongs.
func (cds *Polygon) UnmarshalText(b []byte) error {
	return cds.MultiPoint.UnmarshalText(b)
}
*/

// Type returns this type
func (Polygon) Type() string {
	return "Polygon"
}

// S2Loop is getter for s2.Loop.
func (cds Polygon) S2Loop() *s2.Loop {
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

// S2Region is getter for s2.Region.
func (cds Polygon) S2Region() s2.Region {
	return cds.S2Loop()
}

// CapBound is for s2.Region interface.
func (cds *Polygon) CapBound() s2.Cap {
	return cds.S2Loop().CapBound()
}

// RectBound is for s2.Region interface.
func (cds *Polygon) RectBound() s2.Rect {
	return cds.S2Loop().RectBound()
}

// ContainsCell is for s2.Region interface.
func (cds *Polygon) ContainsCell(c s2.Cell) bool {
	return cds.S2Loop().ContainsCell(c)
}

// IntersectsCell is for s2.Region interface.
func (cds *Polygon) IntersectsCell(c s2.Cell) bool {
	return cds.S2Loop().IntersectsCell(c)
}

// ContainsPoint is for s2.Region interface.
func (cds *Polygon) ContainsPoint(p s2.Point) bool {
	return cds.S2Loop().ContainsPoint(p)
}

// CellUnionBound is for s2.Region interface.
func (cds *Polygon) CellUnionBound() []s2.CellID {
	return cds.S2Loop().CellUnionBound()
}

// S2Point is Center LatLng
func (cds Polygon) S2Point() s2.Point {
	return cds.S2Loop().Centroid()
}

// Radiusp is un-used
func (cds Polygon) Radiusp() *float64 {
	return nil
}

// MarshalJSON is a marshaler for JSON.
func (cds Polygon) MarshalJSON() ([]byte, error) {
	return json.Marshal(&[]MultiPoint{cds.MultiPoint})
}

// UnmarshalJSON is a unmarshaler for JSON.
func (cds *Polygon) UnmarshalJSON(data []byte) (err error) {
	var co []MultiPoint
	err = json.Unmarshal(data, &co)
	if err != nil {
		panic(err)
	}

	switch len(co) {
	case 0:
		panic("No Polygon!")
	case 1:
		cds.MultiPoint = co[0]
	default:
		panic("Polygon has hole! Not implemented")
	}
	return nil
}

// NewGeoJSONGeometry returns GeoJSONGeometry.
func (cds Polygon) NewGeoJSONGeometry() *GeoJSONGeometry {
	var g GeoJSONGeometry
	g.geo = cds
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
