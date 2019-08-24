package latlong

import (
	"bytes"
	"encoding/json"
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

// UnmarshalText is from ISO6709 latlongs.
func (cds *MultiPoint) UnmarshalText(str []byte) error {
	for _, s := range bytes.Split(str, []byte(`/`)) {
		if len(s) != 0 {
			var p Point
			err := p.UnmarshalText(s)
			if err != nil {
				return err
			}
			*cds = append(*cds, p)
		}
	}
	return nil
}

// UnmarshalJSON is from ISO6709 latlongs.
func (cds *MultiPoint) UnmarshalJSON(str []byte) error {
	var v []Point

	if err := json.Unmarshal(str, &v); err != nil {
		return err
	}

	*cds = v
	return nil
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

// Radiusp is un-used
func (cds MultiPoint) Radiusp() *float64 {
	return nil
}

// NewGeoJSONGeometry returns GeoJSONGeometry.
func (cds MultiPoint) NewGeoJSONGeometry() *GeoJSONGeometry {
	var g GeoJSONGeometry
	g.geo = cds
	return &g
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
