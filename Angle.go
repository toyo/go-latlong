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

func NewAngle(degree, degreeprec float64) (a Angle) {
	a.radian = s1.Angle(degree) * s1.Degree
	a.radianprec = s1.Angle(degreeprec) * s1.Degree
	return
}

func NewAngleFromS1Angle(angle, angleprec s1.Angle) (a Angle) {
	a.radian = angle
	a.radianprec = angleprec
	return
}

func (a Angle) Equal(a1 Angle) bool {
	return a == a1
}

func (a Angle) Degrees() float64 {
	return a.radian.Degrees()
}

func (a Angle) PrecDegrees() float64 {
	return a.radianprec.Degrees()
}

func (a Angle) PrecS1Angle() s1.Angle {
	return a.radianprec
}

func (a Angle) Radians() s1.Angle {
	return a.radian
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
	*a = AngleFromBytes(data)

	if isErrorDeg(*a) {
		err = fmt.Errorf("Error Degree on JSON Deg %s", string(data))
	}
	return
}

// AngleFromBytes creates Angle from []byte to unmarshal.
func AngleFromBytes(part []byte) (a Angle) {
	part = bytes.TrimSpace(part)
	pos := bytes.Index(part, []byte(`.`))
	if pos == -1 {
		pos = len(part)
	}

	if pos < 3 && false {
		a = getErrorDeg()
	} else if pos < 5 {
		a = getDeg(part, pos)
	} else if pos < 7 {
		a = getDegMin(part, pos)
	} else if pos < 9 {
		a = getDegMinSec(part, pos)
	} else {
		a = getErrorDeg()
	}
	return
}

func isErrorDeg(a Angle) bool {
	erra := getErrorDeg()
	if a.radian == erra.radian && a.radianprec == erra.radianprec {
		return true
	}
	return false
}

func getErrorDeg() (a Angle) {
	a.radian = 0
	a.radianprec = 360
	return
}

func getDeg(part []byte, pos int) Angle {
	var deg, degprec float64
	var err error
	deg, err = strconv.ParseFloat(string(part), 64)
	if err != nil {
		return getErrorDeg()
	}

	if l := len(part); l == pos {
		degprec = 1
	} else {
		degprec = math.Pow10(pos - l + 1)
	}
	return Angle{radian: s1.Angle(deg) * s1.Degree, radianprec: s1.Angle(degprec) * s1.Degree}
}

func getDegMin(part []byte, pos int) Angle {
	var err error
	var deg, degprec float64
	if deg, err = strconv.ParseFloat(string(part[1:pos-2]), 64); err != nil {
		return getErrorDeg()
	}

	var min float64
	if min, err = strconv.ParseFloat(string(part[pos-2:]), 64); err != nil {
		return getErrorDeg()

	}
	deg += min / 60

	switch part[0] {
	case '-':
		deg = -deg
	case '+':
		break
	default:
		return getErrorDeg()
	}

	if l := len(part); l == pos {
		degprec = float64(1) / 60
	} else {
		degprec = math.Pow10(pos-l+1) / 60
	}

	return Angle{radian: s1.Angle(deg) * s1.Degree, radianprec: s1.Angle(degprec) * s1.Degree}
}

func getDegMinSec(part []byte, pos int) Angle {
	var err error
	var deg, degprec float64
	if deg, err = strconv.ParseFloat(string(part[1:pos-4]), 64); err != nil {
		return getErrorDeg()
	}

	var min float64
	if min, err = strconv.ParseFloat(string(part[pos-4:pos-2]), 64); err != nil {
		return getErrorDeg()
	}
	deg += min / 60

	var sec float64
	if sec, err = strconv.ParseFloat(string(part[pos-2:]), 64); err != nil {
		return getErrorDeg()
	}
	deg += sec / 3600

	switch part[0] {
	case '-':
		deg = -deg
	case '+':
		break
	default:
		return getErrorDeg()
	}

	if l := len(part); l == pos {
		degprec = float64(1) / 3600
	} else {
		degprec = math.Pow10(pos-l+1) / 3600
	}

	return Angle{radian: s1.Angle(deg) * s1.Degree, radianprec: s1.Angle(degprec) * s1.Degree}
}
