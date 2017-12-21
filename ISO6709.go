package latlong

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

func getErrorDeg() (deg float64, degprec float64) {
	deg = 0
	degprec = 360
	return
}

func getDeg(part string, pos int) (deg float64, degprec float64) {
	var err error
	deg, err = strconv.ParseFloat(part, 64)
	if err != nil {
		deg, degprec = getErrorDeg()
		return
	}

	if l := len(part); l == pos {
		degprec = 1
	} else {
		degprec = math.Pow10(pos - l + 1)
	}
	return
}

func getDegMin(part string, pos int) (deg float64, degprec float64) {
	var err error
	if deg, err = strconv.ParseFloat(part[1:pos-2], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}

	var min float64
	if min, err = strconv.ParseFloat(part[pos-2:], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}
	deg += min / 60

	switch part[0] {
	case '-':
		deg = -deg
	case '+':
		break
	default:
		deg, degprec = getErrorDeg()
		return
	}

	if l := len(part); l == pos {
		degprec = 1 / 60
	} else {
		degprec = math.Pow10(pos-l+1) / 60
	}

	return
}

func getDegMinSec(part string, pos int) (deg float64, degprec float64) {
	var err error
	if deg, err = strconv.ParseFloat(part[1:pos-4], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}

	var min float64
	if min, err = strconv.ParseFloat(part[pos-4:pos-2], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}
	deg += min / 60

	var sec float64
	if sec, err = strconv.ParseFloat(part[pos-2:], 64); err != nil {
		deg, degprec = getErrorDeg()
		return
	}
	deg += sec / 3600

	switch part[0] {
	case '-':
		deg = -deg
	case '+':
		break
	default:
		deg, degprec = getErrorDeg()
		return
	}

	if l := len(part); l == pos {
		degprec = float64(1) / 3600
	} else {
		degprec = math.Pow10(pos-l+1) / 3600
	}

	return
}

// NewLatLongISO6709 is from ISO6709 string
func NewLatLongISO6709(iso6709 string) (ll LatLong) {
	re := regexp.MustCompile(`(?P<Latitude>[\+-][\d.]+)(?P<Longitude>[\+-][\d.]+)(?P<Altitude>[\+-][\d.]+)?`)
	var latitude, longitude, latprec, longprec float64
	var altitude *float64

	if iso6709 != "" {
		if re.MatchString(iso6709) {
			match := re.FindStringSubmatch(iso6709)
			for i, name := range re.SubexpNames() {
				part := match[i]
				if i == 0 || name == "" || part == "" {
					continue
				}
				pos := strings.Index(part, ".")
				if pos == -1 {
					pos = len(part)
				}
				switch name {
				case "Latitude":
					if pos < 2 {
						latitude, latprec = getErrorDeg()
					} else if pos < 4 {
						latitude, latprec = getDeg(part, pos)
					} else if pos < 6 {
						latitude, latprec = getDegMin(part, pos)
					} else if pos < 8 {
						latitude, latprec = getDegMinSec(part, pos)
					} else {
						latitude, latprec = getErrorDeg()
					}
				case "Longitude":
					if pos < 3 {
						latitude, latprec = getErrorDeg()
					} else if pos < 5 {
						longitude, longprec = getDeg(part, pos)
					} else if pos < 7 {
						longitude, longprec = getDegMin(part, pos)
					} else if pos < 9 {
						longitude, longprec = getDegMinSec(part, pos)
					} else {
						latitude, latprec = getErrorDeg()
					}
				case "Altitude":
					if a, er := strconv.ParseFloat(part, 64); er == nil {
						altitude = &a
					}
				}
			}
		}
	}

	if altitude != nil {
		ll = NewLatLongAlt(latitude, longitude, latprec, longprec, *altitude)
	} else {
		ll = NewLatLong(latitude, longitude, latprec, longprec)
	}
	return
}
