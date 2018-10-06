package latlong

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/golang/geo/s2"
)

// MultiPoint is slice of *Point
type MultiPoint []Point

// Type returns this type
func (MultiPoint) Type() string {
	return "MultiPoint"
}

// Point returns the first point.
func (cds MultiPoint) Point() Point {
	return cds[0]
}

// S2Point is Center LatLng
func (cds MultiPoint) S2Point() s2.Point {
	return cds[0].S2Point()
}

// S2Region is nil
func (cds MultiPoint) S2Region() s2.Region {
	return nil
}

// NewMultiPointISO6709 is from ISO6709 latlongs.
func NewMultiPointISO6709(str []byte) *MultiPoint {
	var ll MultiPoint
	for _, s := range bytes.Split(str, []byte(`/`)) {
		if len(s) != 0 {
			l := NewPointISO6709(s)
			if l != nil {
				ll = append(ll, *l)
			} else {
				return nil
			}
		}
	}
	return &ll
}

func (cds *MultiPoint) reverse() {
	for i, j := 0, len(*cds)-1; i < j; i, j = i+1, j-1 {
		(*cds)[i], (*cds)[j] = (*cds)[j], (*cds)[i]
	}
}

func (cds *MultiPoint) unset(i int) {
	l := *cds
	if i >= len(l) {
		return
	}
	l = append(l[:i], l[i+1:]...)
	*cds = l
}

// Uniq merge same element
func (cds *MultiPoint) Uniq() {
	if cds != nil && len(*cds) >= 2 {
		ls := *cds

		l := len(ls)
		if ls[l-2].Lat() == ls[l-1].Lat() && ls[l-2].Lng() == ls[l-1].Lng() { //Same point, different precision
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

func (cds MultiPoint) String() string {
	var ss []string
	for _, l := range cds {
		ss = append(ss, l.String())
	}
	return strings.Join(ss, ",")
}

// UnmarshalXML is Unmarshal function.
func (cds *MultiPoint) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
			b := NewMultiPointISO6709(t)
			if b == nil {
				return errors.New("Unexpected CharData on Coordinates UnmarshalXML")
			}
			for i := range *b {
				*cds = append(*cds, (*b)[i])
			}
		case xml.EndElement:
			return nil
		default:
			return fmt.Errorf("Unexpected Token on LatLongs %v", reflect.TypeOf(token))
		}
	}
}

// Radiusp is un-used
func (cds MultiPoint) Radiusp() *float64 {
	return nil
}

// NewGeoJSONGeometry returns GeoJSONGeometry.
func (cds MultiPoint) NewGeoJSONGeometry() GeoJSONGeometry {
	var g GeoJSONGeometry
	g.geo = cds
	return g
}

// NewGeoJSONFeature returns GeoJSONFeature.
func (cds MultiPoint) NewGeoJSONFeature(property interface{}) *GeoJSONFeature {
	var g GeoJSONFeature
	g.Type = "Feature"
	g.Geometry = cds.NewGeoJSONGeometry()
	g.Property = property
	return &g
}

// Equal return bool
func (cds MultiPoint) Equal(c1 Geometry) bool {
	c := c1.(MultiPoint)
	if len(cds) != len(c) {
		return false
	}
	for i := range cds {
		if cds[i] != c[i] {
			return false
		}
	}

	return true
}
