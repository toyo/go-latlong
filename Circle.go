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
	latprecchordangle := s1.ChordAngleFromAngle(latlng.latprec)
	lngprecchordangle := s1.ChordAngleFromAngle(latlng.lngprec) * s1.ChordAngle(math.Abs(math.Cos(float64(latlng.latprec.Radians()))))
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
		Point:      *NewPoint(0, 0, 0, 0),
		ChordAngle: s1.NegativeChordAngle,
	}
	return &circle
}

// S2Region is getter for s2.Region.
func (c *Circle) S2Region() s2.Cap {
	return s2.CapFromCenterChordAngle(s2.PointFromLatLng(c.Point.LatLng), c.ChordAngle)
}

// Radius returns radius of circle.
func (c *Circle) Radius() Km {
	return EarthArcFromChordAngle(c.ChordAngle)
}

func (c *Circle) String() string {
	return c.Point.String() + "/" + c.Radius().String()
}

// CircumferenceToPole returns circle circumference point nearest to other pole.
func (c *Circle) CircumferenceToPole() Point {
	var l s2.LatLng
	if c.Point.Lat > 0 { // if north side, nearest to south pole.
		l.Lat = c.Lat - c.ChordAngle.Angle()
		l.Lng = c.Lng
	} else { // if south side, nearest to north pole.
		l.Lat = c.Lat + c.ChordAngle.Angle()
		l.Lng = c.Lng
	}
	return Point{LatLng: l}
}

// S2Loop is circumference loop.
// div is nomber of vertices.
func (c *Circle) S2Loop(div int) (loop *s2.Loop) {
	return s2.RegularLoop(s2.PointFromLatLng(c.Point.LatLng), c.Angle(), div)
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
func (c Circle) NewGeoJSONGeometry() *GeoJSONGeometry {
	var g GeoJSONGeometry
	g.Type = "Circle"
	g.Coordinates = []interface{}{c.Point.Lng.Degrees(), c.Point.Lat.Degrees()}
	radius := float64(c.Radius())
	g.Radius = &radius

	return &g
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
		ll := Point{LatLng: s2.LatLngFromPoint(vs[i])}
		lls = append(lls, ll)
	}
	return
}
