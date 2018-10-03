package latlong

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Circle is s2.Cap
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

// S2Region is getter for s2.S2Region.
func (c *Circle) S2Region() *s2.Cap {
	cap := s2.CapFromCenterChordAngle(s2.PointFromLatLng(c.Point.LatLng), c.ChordAngle)
	return &cap
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
// radian is one vertex degree.
func (c *Circle) S2Loop(radian s1.Angle) (loop *s2.Loop) {
	p := s2.PointFromLatLng(c.CircumferenceToPole().LatLng)
	axis := s2.PointFromLatLng(c.Point.LatLng)

	var pss s2.Polyline
	for angle := s1.Angle(0); angle < 2*math.Pi; angle += radian {
		pss = append(pss, s2.Rotate(p, axis, angle))
	}
	pss = append(pss, s2.Rotate(p, axis, 0))

	loop = s2.LoopFromPoints(pss)

	return
}

// S2LatLngs is circumference loop by []s2.LatLng.
// radian is one vertex degree.
func (c *Circle) S2LatLngs(radian s1.Angle) (lls []s2.LatLng) {
	vs := c.S2Loop(radian).Vertices()
	for i := range vs {
		lls = append(lls, s2.LatLngFromPoint(vs[i]))
	}
	return
}

// LatLngs is circumference loop by []LatLng.LatLngs
// radian is one vertex degree.
func (c *Circle) LatLngs(radian s1.Angle) (lls []Point) {
	vs := c.S2Loop(radian).Vertices()
	for i := range vs {
		ll := Point{LatLng: s2.LatLngFromPoint(vs[i])}
		lls = append(lls, ll)
	}
	return
}

// CircleGeoJSON is GeoJSON format for Circle.
type CircleGeoJSON struct {
	Type       string     `json:"type"`
	Coordinate Coordinate `json:"coordinates"`
	Radius     *Km        `json:"radius"`
}

// MarshalJSON is a marshaler for JSON.
func (c Circle) MarshalJSON() (data []byte, err error) {
	var js CircleGeoJSON
	js.Type = "Circle"
	js.Coordinate.Point = c.Point
	km := EarthArcFromChordAngle(c.ChordAngle)
	if !math.IsNaN(float64(km)) {
		js.Radius = &km
	}
	return json.Marshal(&js)
}

// UnmarshalJSON is a unmarshaler for JSON.
func (c *Circle) UnmarshalJSON(data []byte) error {
	var js CircleGeoJSON
	err := json.Unmarshal(data, &js)
	if js.Type != "Circle" {
		err := errors.New("JSON Type is Circle but Type != Circle")
		return err
	}
	c.Point = js.Coordinate.Point
	if js.Radius != nil {
		c.ChordAngle = js.Radius.EarthChordAngle()
	}
	return err
}
