package latlong

import (
	"errors"
	"math"

	geohash "github.com/TomiHiltunen/geohash-golang"
)

// NewLatLongGeoHash is from GeoHash
// http://geohash.org/
func NewLatLongGeoHash(geoHash string) (latlong LatLong, err error) {
	if bb := geohash.Decode(geoHash); bb != nil {
		latlong = NewLatLong(bb.Center().Lat(), bb.Center().Lng(), bb.NorthEast().Lat()-bb.SouthWest().Lat(), bb.NorthEast().Lng()-bb.SouthWest().Lng())
		//fmt.Println(bb.NorthEast(), bb.SouthWest())
	} else {
		err = errors.New("Geohash decode error")
	}
	return
}

func (latlong *LatLong) geoHash(precision int) string {
	return geohash.EncodeWithPrecision(latlong.Lat(), latlong.Lng(), precision)
}

// GeoHash5 returns GeoHash string.
func (latlong *LatLong) GeoHash5() string {
	return latlong.geoHash(5)
}

// GeoHash6 returns GeoHash string.
func (latlong *LatLong) GeoHash6() string {
	return latlong.geoHash(6)
}

// GeoHash returns GeoHash string.
func (latlong *LatLong) GeoHash() string {
	const floaterr = 1 + 5E-10

	geohashlatbits := -math.Log2(latlong.Size().Lat.Degrees()/45) + 2 // div by 180 = 45 * 2^2
	geohashlngbits := -math.Log2(latlong.Size().Lng.Degrees()/45) + 3 // div by 360 = 45 * 2^3
	//fmt.Printf("lat %.99f, lng %.99f\n", geohashlatbits, geohashlngbits)
	//fmt.Printf("lat %.9f, lng %.9f\n", latlong.Size().Lat.Degrees(), latlong.Size().Lng.Degrees())

	geohashlat2len, geohashlatlen2mod := math.Modf(geohashlatbits / 5 * floaterr)
	//fmt.Printf("lat %f mod %f\n", geohashlat2len, geohashlatlen2mod)

	var geohashlatlen int
	if geohashlatlen2mod >= 0.4 {
		geohashlatlen = int(geohashlat2len)*2 + 1
	} else {
		geohashlatlen = int(geohashlat2len) * 2
	}

	geohashlng2len, geohashlnglen2mod := math.Modf(geohashlngbits / 5 * floaterr)
	//fmt.Printf("lng %f mod %f\n", geohashlng2len, geohashlnglen2mod)

	var geohashlnglen int
	if geohashlnglen2mod >= 0.6 {
		geohashlnglen = int(geohashlng2len)*2 + 1
	} else {
		geohashlnglen = int(geohashlng2len) * 2
	}
	//fmt.Printf("%d, %d\n", geohashlatlen, geohashlnglen)

	if geohashlatlen < geohashlnglen {
		return latlong.geoHash(geohashlatlen)
	}
	return latlong.geoHash(geohashlnglen)
}
