package latlong

import (
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Circle is s2.Cap
type Circle struct {
	LatLng
	s1.ChordAngle
}

// NewCircle is constuctor for Circle
func NewCircle(latlng LatLng, km Km) *Circle {
	circle := Circle{
		LatLng:     latlng,
		ChordAngle: s1.ChordAngleFromAngle(km.EarthAngle()),
	}
	return &circle
}

// NewPointCircle is constructor for Circle with radius = prec
func NewPointCircle(latlng LatLng) *Circle {
	latprecchordangle := s1.ChordAngleFromAngle(latlng.latprec)
	lngprecchordangle := s1.ChordAngleFromAngle(latlng.lngprec) * s1.ChordAngle(math.Abs(math.Cos(float64(latlng.latprec.Radians()))))
	var precchordangle s1.ChordAngle

	if latprecchordangle > lngprecchordangle {
		precchordangle = latprecchordangle
	} else {
		precchordangle = lngprecchordangle
	}

	circle := Circle{
		LatLng:     latlng,
		ChordAngle: precchordangle,
	}
	return &circle

}

// NewEmptyCircle is constructor for Circle with empty.
func NewEmptyCircle() *Circle {
	circle := Circle{
		LatLng:     *NewLatLng(0, 0, 0, 0),
		ChordAngle: s1.NegativeChordAngle,
	}
	return &circle
}

// S2Region is getter for s2.S2Region.
func (c *Circle) S2Region() *s2.Cap {
	cap := s2.CapFromCenterChordAngle(s2.PointFromLatLng(c.LatLng.LatLng), c.ChordAngle)
	return &cap
}

// Radius returns radius of circle.
func (c *Circle) Radius() Km {
	return EarthArcFromChordAngle(c.ChordAngle)
}

func (c *Circle) String() string {
	return c.LatLng.String() + "/" + c.Radius().String()
}
