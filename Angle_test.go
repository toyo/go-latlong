package latlong_test

import (
	"bytes"
	"encoding/json"
	"testing"

	latlong "github.com/toyo/go-latlong"
)

func TestAngle(t *testing.T) {
	var a latlong.Angle
	b := []byte(`270`)
	err := json.Unmarshal(b, &a)
	if err != nil {
		t.Error(err.Error())
	}

	expct := `270`
	if a.String() != expct {
		t.Errorf("Got %s expect %s", a.String(), expct)
	}

	bb, err := json.Marshal(&a)
	if err != nil {
		t.Error(err.Error())
	}
	if !bytes.Equal(b, bb) {
		t.Errorf(`Got "%s" expect "%s"`, string(b), string(bb))
	}
}
