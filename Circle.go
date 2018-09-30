package latlong

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Circle is s2.Cap
type Circle struct {
	//LatLng
	//s1.Angle
	s2.Cap
	latprec s1.Angle
	lngprec s1.Angle
}

// NewCircle is constuctor for Cap
func NewCircle(latlng LatLng, km Km) *Circle {
	circle := Circle{
		Cap:     s2.CapFromCenterAngle(s2.PointFromLatLng(latlng.LatLng), km.EarthAngle()),
		latprec: latlng.latprec,
		lngprec: latlng.lngprec,
	}
	return &circle
}

// Center returns LatLng of center.
func (c *Circle) Center() LatLng {
	ll := LatLng{
		LatLng:  s2.LatLngFromPoint(c.Cap.Center()),
		latprec: c.latprec,
		lngprec: c.lngprec,
	}
	return ll
}

// Radius returns radius of circle.
func (c *Circle) Radius() Km {
	return EarthArcFromAngle(c.Cap.Radius())
}

func (c *Circle) String() string {
	return c.Center().String() + "/" + c.Radius().String()
}
