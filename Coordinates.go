package latlong

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/golang/geo/s2"
)

// Coordinates is slice of LatLong
type Coordinates []*Coordinate

// NewLatLongsISO6709 is from ISO6709 latlongs.
func NewLatLongsISO6709(str string) (ll *Coordinates, err error) {
	ll = new(Coordinates)
	for _, s := range strings.Split(str, "/") {
		if s != "" {
			l := NewLatLongISO6709(s)
			if err == nil {
				*ll = append(*ll, l)
			} else {
				return
			}
		}
	}
	err = nil
	return
}

// S2Polyline is getter for s2.Polyline ([]s2.Point).
func (latlongs *Coordinates) S2Polyline() (ps s2.Polyline) {
	for _, v := range *latlongs {
		ps = append(ps, v.S2Point())
	}
	return
}

// S2Loop is getter for s2.Loop.
func (latlongs *Coordinates) S2Loop() *s2.Loop {
	lo := s2.LoopFromPoints(latlongs.S2Polyline())
	if lo.TurningAngle() < 0 { // if loop is not CCW but CW,
		lo.Invert() // Change to CCW.
	}
	return lo
}

func (a *Coordinates) unset(i int) {
	l := *a
	if i >= len(l) {
		return
	}
	l = append(l[:i], l[i+1:]...)
	*a = l
}

// Uniq merge same element
func (a *Coordinates) Uniq() {
	if a != nil && len(*a) >= 2 {
		ls := *a

		l := len(ls)
		if ls[l-2].Intersects(ls[l-1].Rect) { //Same point, different precision
			if ls[l-2].Area() < ls[l-1].Area() {
				ls.unset(l - 1)
				ls.Uniq()
			} else {
				ls.unset(l - 2)
				ls.Uniq()
			}
		} else {
			ll := ls[l-1]
			lm := ls[:l-1]
			lm.Uniq()
			ls = append(lm, ll)
		}
		*a = ls
	}
}

func (a Coordinates) String() string {
	a.Uniq()

	var ss []string
	for _, l := range a {
		ss = append(ss, l.String())
	}
	return strings.Join(ss, ",")
}

// UnmarshalXML is Unmarshal function but NOT WORK.
func (a *Coordinates) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var token xml.Token

	token, err = d.Token()
	if err != nil {
		return
	}
	if token == io.EOF {
		err = errors.New("Unexpected EOF on LatLongs")
		return
	}

	switch t := token.(type) {
	case xml.CharData:
		var b *Coordinates
		b, err = NewLatLongsISO6709(string(t))
		*a = *b
		return
	default:
		err = errors.New("Unexpected Token on LatLongs")
		return
	}
}
