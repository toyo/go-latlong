package latlong

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/golang/geo/s1"
)

type Angle struct {
	radian     s1.Angle
	radianprec s1.Angle
}

func (a Angle) preclog() (lngprec int) {
	if a.radianprec.Degrees() != 0 {
		lngprec = int(math.Ceil(-math.Log10(a.radianprec.Degrees())))
		if lngprec < 0 {
			lngprec = 0
		}
	} else {
		lngprec = 2
	}
	return
}

func (a Angle) String() (s string) {
	return strconv.FormatFloat(a.radian.Degrees(), 'f', a.preclog(), 64)
}

// MarshalJSON is a marshaler for JSON.
func (a Angle) MarshalJSON() ([]byte, error) {
	return []byte(a.String()), nil
}

// UnmarshalJSON is a unmarshaler for JSON.
func (a *Angle) UnmarshalJSON(data []byte) (err error) {
	data = bytes.TrimSpace(data)
	pos := bytes.Index(data, []byte(`.`))
	if pos == -1 {
		pos = len(data)
	}

	deg, degprec := getDeg(data, pos)

	if isErrorDeg(deg, degprec) {
		err = fmt.Errorf("Error Degree on JSON Deg %s", string(data))
	}

	a.radian = s1.Angle(deg) * s1.Degree
	a.radianprec = s1.Angle(degprec) * s1.Degree
	return
}
