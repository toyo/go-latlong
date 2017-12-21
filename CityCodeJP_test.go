package latlong

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCityCodeJP(t *testing.T) {
	response := `{"ResultInfo":{"Count":1,"Total":1,"Start":1,"Latency":0.0062708854675293,"Status":200,"Description":"\u6307\u5b9a\u306e\u5730\u70b9\u306e\u4f4f\u6240\u60c5\u5831\u3092\u53d6\u5f97\u3059\u308b\u6a5f\u80fd\u3092\u63d0\u4f9b\u3057\u307e\u3059\u3002","Copyright":"Copyright (C) 2017 Yahoo Japan Corporation. All Rights Reserved.","CompressType":""},"Feature":[{"Property":{"Country":{"Code":"JP","Name":"\u65e5\u672c"},"Address":"\u4e09\u91cd\u770c\u4f0a\u52e2\u5e02\u5b87\u6cbb\u9928\u753a","AddressElement":[{"Name":"\u4e09\u91cd\u770c","Kana":"\u307f\u3048\u3051\u3093","Level":"prefecture","Code":"24"},{"Name":"\u4f0a\u52e2\u5e02","Kana":"\u3044\u305b\u3057","Level":"city","Code":"24203"},{"Name":"\u5b87\u6cbb\u9928\u753a","Kana":"\u3046\u3058\u305f\u3061\u3061\u3087\u3046","Level":"oaza"},{"Name":"","Kana":"","Level":"aza"}]},"Geometry":{"Type":"point","Coordinates":"136.725739,34.455846"}}]}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("access!")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	Config.YahooJPAPIURL = server.URL

	l := NewLatLongISO6709("+34.455846+136.725739/")

	Config.Lang = "ja"
	Config.GoogleAPIKey = "AIzaNotReallyAnAPIKey"
	Config.YahooJPClientID = "NotReallyAnAPIKey"

	s, err := l.CityCodeJP()
	if err != nil {
		t.Errorf("Locality returned non nil error: %v", err)
	}

	correctResponse := "24203"
	if !reflect.DeepEqual(s, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, s)
	}
}
