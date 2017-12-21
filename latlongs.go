package latlong

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

// LatLongs is slice of LatLong
type LatLongs []LatLong

// NewLatLongsISO6709 is from ISO6709 latlongs.
func NewLatLongsISO6709(str string) (ll LatLongs, err error) {
	for _, s := range strings.Split(str, "/") {
		if s != "" {
			l := NewLatLongISO6709(s)
			if err == nil {
				ll = append(ll, l)
			} else {
				return
			}
		}
	}
	err = nil
	return
}

func (a *LatLongs) unset(i int) {
	l := *a
	if i >= len(l) {
		return
	}
	l = append(l[:i], l[i+1:]...)
	*a = l
}

// Uniq merge same element
func (a *LatLongs) Uniq() {
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

func (a LatLongs) String() string {
	a.Uniq()

	var ss []string
	for _, l := range a {
		ss = append(ss, l.String())
	}
	return strings.Join(ss, ",")
}

// UnmarshalXML is Unmarshal function but NOT WORK.
func (a *LatLongs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
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
		*a, err = NewLatLongsISO6709(string(t))
		return
	default:
		err = errors.New("Unexpected Token on LatLongs")
		return
	}
}
