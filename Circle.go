package latlong

import (
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Circle is Circle
type Circle struct {
	Point
	s1.ChordAngle
}

// Type returns this type
func (Circle) Type() string {
	return "Circle"
}

// NewCircle is constuctor for Circle
func NewCircle(latlng Point, km Km) *Circle {
	circle := Circle{
		Point:      latlng,
		ChordAngle: s1.ChordAngleFromAngle(km.EarthAngle()),
	}
	return &circle
}

// NewPointCircle is constructor for Circle with radius = prec
func NewPointCircle(latlng Point) *Circle {
	latprecchordangle := s1.ChordAngleFromAngle(latlng.lat.PrecS1Angle())
	lngprecchordangle := s1.ChordAngleFromAngle(latlng.lng.PrecS1Angle()) * s1.ChordAngle(math.Abs(math.Cos(float64(latlng.lat.PrecS1Angle()))))
	var precchordangle s1.ChordAngle

	if latprecchordangle > lngprecchordangle {
		precchordangle = latprecchordangle
	} else {
		precchordangle = lngprecchordangle
	}

	circle := Circle{
		Point:      latlng,
		ChordAngle: precchordangle,
	}
	return &circle

}

// NewEmptyCircle is constructor for Circle with empty.
func NewEmptyCircle() *Circle {
	circle := Circle{
		ChordAngle: s1.NegativeChordAngle,
	}
	return &circle
}

// S2Region is getter for s2.Region.
func (c Circle) S2Region() s2.Region {
	return s2.CapFromCenterChordAngle(s2.PointFromLatLng(c.Point.S2LatLng()), c.ChordAngle)
}

// CapBound is for s2.Region interface.
func (c *Circle) CapBound() s2.Cap {
	return c.S2Region().CapBound()
}

// RectBound is for s2.Region interface.
func (c *Circle) RectBound() s2.Rect {
	return c.S2Region().RectBound()
}

// ContainsCell is for s2.Region interface.
func (c *Circle) ContainsCell(cell s2.Cell) bool {
	return c.S2Region().ContainsCell(cell)
}

// IntersectsCell is for s2.Region interface.
func (c *Circle) IntersectsCell(cell s2.Cell) bool {
	return c.S2Region().IntersectsCell(cell)
}

// ContainsPoint is for s2.Region interface.
func (c *Circle) ContainsPoint(p s2.Point) bool {
	return c.S2Region().ContainsPoint(p)
}

// CellUnionBound is for s2.Region interface.
func (c *Circle) CellUnionBound() []s2.CellID {
	return c.S2Region().CellUnionBound()
}

// Radiusp is un-used
func (c Circle) Radiusp() *float64 {
	r := float64(c.Radius())
	return &r
}

// Radius returns radius of circle.
func (c Circle) Radius() Km {
	return EarthArcFromChordAngle(c.ChordAngle)
}

func (c Circle) String() string {
	return c.Point.String() + "/" + c.Radius().String()
}

// S2Point is Center LatLng
func (c Circle) S2Point() s2.Point {
	return c.Point.S2Point()
}

// S2Loop is circumference loop.
// div is nomber of vertices.
func (c *Circle) S2Loop(div int) (loop *s2.Loop) {
	return s2.RegularLoop(s2.PointFromLatLng(c.Point.S2LatLng()), c.Angle(), div)
}

// S2LatLngs is circumference loop by []s2.LatLng.
// radian is one vertex degree.
func (c *Circle) S2LatLngs(div int) (lls []s2.LatLng) {
	vs := c.S2Loop(div).Vertices()
	for i := range vs {
		lls = append(lls, s2.LatLngFromPoint(vs[i]))
	}
	return
}

// NewGeoJSONGeometry returns GeoJSONGeometry.
func (c Circle) NewGeoJSONGeometry() (g GeoJSONGeometry) {
	g.geo = c
	return
}

// NewGeoJSONFeature returns GeoJSONFeature.
func (c Circle) NewGeoJSONFeature(property interface{}) *GeoJSONFeature {
	var g GeoJSONFeature
	g.Type = "Feature"
	g.Geometry = c.NewGeoJSONGeometry()
	g.Property = property
	return &g
}

// LatLngs is circumference loop by []LatLng.LatLngs
// radian is one vertex degree.
func (c *Circle) LatLngs(div int) (lls []Point) {
	vs := c.S2Loop(div).Vertices()
	for i := range vs {
		s2ll := s2.LatLngFromPoint(vs[i])
		ll := Point{lat: NewAngleFromS1Angle(s2ll.Lat, 0), lng: NewAngleFromS1Angle(s2ll.Lng, 0)}
		lls = append(lls, ll)
	}
	return
}

// Equal return bool
func (c Circle) Equal(c1 Geometry) bool {
	return c == c1
}
