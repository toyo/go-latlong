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
	if len(cds) == 0 {
		panic(cds)
	}
	return cds[0]
}

// S2Point is Center LatLng
func (cds MultiPoint) S2Point() s2.Point {
	return cds.Point().S2Point()
}

// S2Region is nil
func (cds MultiPoint) S2Region() s2.Region {
	return nil
}

// NewMultiPointISO6709 is from ISO6709 latlongs.
func NewMultiPointISO6709(str []byte) (ll MultiPoint) {
	for _, s := range bytes.Split(str, []byte(`/`)) {
		if len(s) != 0 {
			ll = append(ll, NewPointISO6709(s))
		}
	}
	return
}

// Reverse returns reverse order.
func (cds MultiPoint) Reverse() MultiPoint {
	for i, j := 0, len(cds)-1; i < j; i, j = i+1, j-1 {
		cds[i], cds[j] = cds[j], cds[i]
	}
	return cds
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
			if len(t) == 0 {
				continue
			}
			b := NewMultiPointISO6709(t)
			if b == nil {
				return errors.New("Unexpected CharData on Coordinates UnmarshalXML")
			}
			for i := range b {
				*cds = append(*cds, b[i])
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
