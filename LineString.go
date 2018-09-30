package latlong

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/golang/geo/s2"
	"googlemaps.github.io/maps"
)

// LineString is slice of LatLong
type LineString []*LatLng

// S2Region is getter for s2.Polyline ([]s2.Point).
func (cds LineString) S2Region() *s2.Polyline {
	var ps s2.Polyline
	for _, cd := range cds {
		ps = append(ps, cd.S2Point())
	}
	return &ps
}

// MapsLatLng convert to google maps.
func (cds LineString) MapsLatLng() (mlls []maps.LatLng) {
	for _, cd := range cds {
		mlls = append(mlls, cd.MapsLatLng())
	}
	return
}

func (cds *LineString) unset(i int) {
	l := *cds
	if i >= len(l) {
		return
	}
	l = append(l[:i], l[i+1:]...)
	*cds = l
}

// Uniq merge same element
func (cds *LineString) Uniq() {
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

func (cds LineString) String() string {
	cds.Uniq()

	var ss []string
	for _, l := range cds {
		ss = append(ss, l.String())
	}
	return strings.Join(ss, ",")
}

// UnmarshalXML is Unmarshal function.
func (cds *LineString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
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
			b := NewLatLongsISO6709(string(t))
			if b == nil {
				return errors.New("Unexpected CharData on Coordinates UnmarshalXML")
			}
			for _, l := range *b {
				*cds = append(*cds, &l.LatLng)
			}
		case xml.EndElement:
			return nil
		default:
			return fmt.Errorf("Unexpected Token on LatLongs %v", reflect.TypeOf(token))
		}
	}
}
