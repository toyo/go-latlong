package latlong

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// CityCodeJP return city code.
// http://www.soumu.go.jp/denshijiti/code.html
// https://developer.yahoo.co.jp/webapi/map/openlocalplatform/v1/reversegeocoder.html
func (latlong *Point) CityCodeJP() (code string, err error) {

	req, err := http.NewRequest("GET", Config.YahooJPAPIURL, nil)
	if err != nil {
		return
	}

	values := url.Values{
		"appid":  {Config.YahooJPClientID},
		"lat":    {latlong.LatString()},
		"lon":    {latlong.LngString()},
		"output": {"json"},
	}

	req.URL.RawQuery = values.Encode()

	resp, err := Config.HTTPClient.Do(req)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(resp.Body)

	var v struct {
		ResultInfo struct {
			Count   int
			Latency float32
			Status  int
		}
		Feature []struct {
			Property struct {
				AddressElement []struct {
					Name  string
					Kana  string
					Level string
					Code  string
				}
			}
		}
		Error struct {
			Message string
			Code    int
		}
	}

	err = decoder.Decode(&v)

	resp.Body.Close()

	if v.ResultInfo.Status != 200 {
		err = errors.New(strconv.Itoa(v.Error.Code) + v.Error.Message)
		return
	}

	for _, f := range v.Feature {
		for _, ae := range f.Property.AddressElement {
			if ae.Level == "city" {
				code = ae.Code
				return
			}
		}
	}
	return
}
