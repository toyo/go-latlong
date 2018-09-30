package latlong

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Circle is s2.Cap
type Circle struct {
	LatLng
	s1.ChordAngle
}

// NewCircle is constuctor for Cap
func NewCircle(latlng LatLng, km Km) *Circle {
	circle := Circle{
		LatLng:     latlng,
		ChordAngle: s1.ChordAngleFromAngle(km.EarthAngle()),
	}
	return &circle
}
func NewEmptyCircle() *Circle {
	circle := Circle{
		LatLng:     *NewLatLng(0, 0, 0, 0),
		ChordAngle: s1.NegativeChordAngle,
	}
	return &circle
}

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
