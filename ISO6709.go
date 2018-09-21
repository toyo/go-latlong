package latlong

import (
	"regexp"
)

// NewLatLongISO6709 is from ISO6709 string
func NewLatLongISO6709(iso6709 string) (ll *Coordinate) {
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

				switch name {
				case "Latitude":
					latitude, latprec = getLat(part)
				case "Longitude":
					longitude, longprec = getLng(part)
				case "Altitude":
					altitude = getAlt(part)
				}
			}
		}
	}

	ll = NewLatLongAlt(latitude, longitude, latprec, longprec, altitude)

	return
}
