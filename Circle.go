package latlong

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Circle is s2.Cap
type Circle struct {
	s2.LatLng
	s1.Angle
}

// NewCircle is constuctor for Cap
func NewCircle(latlng LatLng, km Km) *Circle {
	cap := Circle{LatLng: latlng.LatLng, Angle: km.EarthAngle()}
	return &cap
}

// S2Cap is getter for s2.Cap
func (c *Circle) S2Cap() *s2.Cap {
	cap := s2.CapFromCenterAngle(s2.PointFromLatLng(c.LatLng), c.Angle)
	return &cap
}
