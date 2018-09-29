package latlong

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/golang/geo/s2"
	"googlemaps.github.io/maps"
)

// Coordinates is slice of LatLong
type Coordinates []*Coordinate

// NewLatLongsISO6709 is from ISO6709 latlongs.
func NewLatLongsISO6709(str string) *Coordinates {
	var ll Coordinates
	for _, s := range strings.Split(str, "/") {
		if s != "" {
			l := NewLatLongISO6709(s)
			if l != nil {
				ll = append(ll, l)
			} else {
				return nil
			}
		}
	}
	return &ll
}

// S2Polyline is getter for s2.Polyline ([]s2.Point).
func (cds Coordinates) S2Polyline() (ps s2.Polyline) {
	for _, cd := range cds {
		ps = append(ps, cd.S2Point())
	}
	return
}

// S2Loop is getter for s2.Loop.
func (cds Coordinates) S2Loop() *s2.Loop {
	lo := s2.LoopFromPoints(cds.S2Polyline())
	lo.Normalize() // if loop is not CCW but CW, change to CCW.
	return lo
}

// MapsLatLng covert to google maps.
func (cds Coordinates) MapsLatLng() (mlls []maps.LatLng) {
	for _, cd := range cds {
		mlls = append(mlls, cd.MapsLatLng())
	}
	return
}

func (cds *Coordinates) unset(i int) {
	l := *cds
	if i >= len(l) {
		return
	}
	l = append(l[:i], l[i+1:]...)
	*cds = l
}

// Uniq merge same element
func (cds *Coordinates) Uniq() {
	if cds != nil && len(*cds) >= 2 {
		ls := *cds

		l := len(ls)
		if ls[l-2].Lat == ls[l-1].Lat && ls[l-2].Lng == ls[l-1].Lng { //Same point, different precision
			if ls[l-2].PrecisionArea() < ls[l-1].PrecisionArea() {
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
		*cds = ls
	}
}

func (cds Coordinates) String() string {
	cds.Uniq()

	var ss []string
	for _, l := range cds {
		ss = append(ss, l.String())
	}
	return strings.Join(ss, ",")
}

// UnmarshalXML is Unmarshal function but NOT WORK.
func (cds *Coordinates) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	token, err := d.Token()
	if err != nil {
		return err
	}
	if token == io.EOF {
		err = errors.New("Unexpected EOF on LatLongs")
		return err
	}

	switch t := token.(type) {
	case xml.CharData:
		if b := NewLatLongsISO6709(string(t)); b != nil {
			*cds = *b
			return nil
		}
		return errors.New("Unexpected CharData on Coordinates UnmarshalXML")
	default:
		return errors.New("Unexpected Token on LatLongs")
	}
}
