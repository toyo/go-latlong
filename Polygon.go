package latlong

import "github.com/golang/geo/s2"

type Polygon struct {
	LineString
}

// S2Loop is getter for s2.Loop.
func (cds Polygon) S2Region() *s2.Loop {
	lo := s2.LoopFromPoints(*cds.LineString.S2Region())
	lo.Normalize() // if loop is not CCW but CW, change to CCW.
	return lo
}
