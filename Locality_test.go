package latlong_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	latlong "github.com/toyo/go-latlong"
)

func TestLocality(t *testing.T) {
	response := `{
   "results" : [
      {
         "address_components" : [
            {
               "long_name" : "Unnamed Road",
               "short_name" : "Unnamed Road",
               "types" : [ "route" ]
            },
            {
               "long_name" : "伊勢市",
               "short_name" : "伊勢市",
               "types" : [ "locality", "political" ]
            },
            {
               "long_name" : "三重県",
               "short_name" : "三重県",
               "types" : [ "administrative_area_level_1", "political" ]
            },
            {
               "long_name" : "日本",
               "short_name" : "JP",
               "types" : [ "country", "political" ]
            },
            {
               "long_name" : "516-0023",
               "short_name" : "516-0023",
               "types" : [ "postal_code" ]
            }
         ],
         "formatted_address" : "Unnamed Road, 伊勢市 三重県 516-0023 日本",
         "geometry" : {
            "bounds" : {
               "northeast" : {
                  "lat" : 34.4564137,
                  "lng" : 136.725903
               },
               "southwest" : {
                  "lat" : 34.4554832,
                  "lng" : 136.7245453
               }
            },
            "location" : {
               "lat" : 34.4559758,
               "lng" : 136.725824
            },
            "location_type" : "GEOMETRIC_CENTER",
            "viewport" : {
               "northeast" : {
                  "lat" : 34.4572974302915,
                  "lng" : 136.7265731302915
               },
               "southwest" : {
                  "lat" : 34.4545994697085,
                  "lng" : 136.7238751697085
               }
            }
         },
         "place_id" : "ChIJadYR7OhQBGAR4Eu2euhqXWQ",
         "types" : [ "route" ]
      },
      {
         "address_components" : [
            {
               "long_name" : "宇治館町",
               "short_name" : "宇治館町",
               "types" : [ "political", "sublocality", "sublocality_level_1" ]
            },
            {
               "long_name" : "伊勢市",
               "short_name" : "伊勢市",
               "types" : [ "locality", "political" ]
            },
            {
               "long_name" : "三重県",
               "short_name" : "三重県",
               "types" : [ "administrative_area_level_1", "political" ]
            },
            {
               "long_name" : "日本",
               "short_name" : "JP",
               "types" : [ "country", "political" ]
            },
            {
               "long_name" : "516-0023",
               "short_name" : "516-0023",
               "types" : [ "postal_code" ]
            }
         ],
         "formatted_address" : "日本、〒516-0023 三重県伊勢市宇治館町",
         "geometry" : {
            "bounds" : {
               "northeast" : {
                  "lat" : 34.4686677,
                  "lng" : 136.7875221
               },
               "southwest" : {
                  "lat" : 34.4118037,
                  "lng" : 136.7209235
               }
            },
            "location" : {
               "lat" : 34.4620235,
               "lng" : 136.725732
            },
            "location_type" : "APPROXIMATE",
            "viewport" : {
               "northeast" : {
                  "lat" : 34.4686677,
                  "lng" : 136.7875221
               },
               "southwest" : {
                  "lat" : 34.4118037,
                  "lng" : 136.7209235
               }
            }
         },
         "place_id" : "ChIJ_-Wu2MlWBGARvKowXVXbZpI",
         "types" : [ "political", "sublocality", "sublocality_level_1" ]
      },
      {
         "address_components" : [
            {
               "long_name" : "伊勢市",
               "short_name" : "伊勢市",
               "types" : [ "locality", "political" ]
            },
            {
               "long_name" : "三重県",
               "short_name" : "三重県",
               "types" : [ "administrative_area_level_1", "political" ]
            },
            {
               "long_name" : "日本",
               "short_name" : "JP",
               "types" : [ "country", "political" ]
            }
         ],
         "formatted_address" : "日本、三重県伊勢市",
         "geometry" : {
            "bounds" : {
               "northeast" : {
                  "lat" : 34.566524,
                  "lng" : 136.8177107
               },
               "southwest" : {
                  "lat" : 34.3847593,
                  "lng" : 136.633115
               }
            },
            "location" : {
               "lat" : 34.48751439999999,
               "lng" : 136.7093359
            },
            "location_type" : "APPROXIMATE",
            "viewport" : {
               "northeast" : {
                  "lat" : 34.566524,
                  "lng" : 136.8177107
               },
               "southwest" : {
                  "lat" : 34.3847593,
                  "lng" : 136.633115
               }
            }
         },
         "place_id" : "ChIJxW5CEJlQBGARGRH-87cfV_o",
         "types" : [ "locality", "political" ]
      },
      {
         "address_components" : [
            {
               "long_name" : "516-0023",
               "short_name" : "516-0023",
               "types" : [ "postal_code" ]
            },
            {
               "long_name" : "宇治館町",
               "short_name" : "宇治館町",
               "types" : [ "political", "sublocality", "sublocality_level_1" ]
            },
            {
               "long_name" : "伊勢市",
               "short_name" : "伊勢市",
               "types" : [ "locality", "political" ]
            },
            {
               "long_name" : "三重県",
               "short_name" : "三重県",
               "types" : [ "administrative_area_level_1", "political" ]
            },
            {
               "long_name" : "日本",
               "short_name" : "JP",
               "types" : [ "country", "political" ]
            }
         ],
         "formatted_address" : "〒516-0023, 日本",
         "geometry" : {
            "bounds" : {
               "northeast" : {
                  "lat" : 34.4686677,
                  "lng" : 136.7875221
               },
               "southwest" : {
                  "lat" : 34.4118037,
                  "lng" : 136.7209235
               }
            },
            "location" : {
               "lat" : 34.4440555,
               "lng" : 136.7526854
            },
            "location_type" : "APPROXIMATE",
            "viewport" : {
               "northeast" : {
                  "lat" : 34.4686677,
                  "lng" : 136.7875221
               },
               "southwest" : {
                  "lat" : 34.4118037,
                  "lng" : 136.7209235
               }
            }
         },
         "place_id" : "ChIJn0bqG-tQBGARyV-M0vnjAnw",
         "types" : [ "postal_code" ]
      },
      {
         "address_components" : [
            {
               "long_name" : "三重県",
               "short_name" : "三重県",
               "types" : [ "administrative_area_level_1", "political" ]
            },
            {
               "long_name" : "日本",
               "short_name" : "JP",
               "types" : [ "country", "political" ]
            }
         ],
         "formatted_address" : "日本、三重県",
         "geometry" : {
            "bounds" : {
               "northeast" : {
                  "lat" : 35.2576844,
                  "lng" : 136.9877523
               },
               "southwest" : {
                  "lat" : 33.7226405,
                  "lng" : 135.8528648
               }
            },
            "location" : {
               "lat" : 34.7302829,
               "lng" : 136.5085883
            },
            "location_type" : "APPROXIMATE",
            "viewport" : {
               "northeast" : {
                  "lat" : 35.2576844,
                  "lng" : 136.9877523
               },
               "southwest" : {
                  "lat" : 33.7226405,
                  "lng" : 135.8528648
               }
            }
         },
         "place_id" : "ChIJY0FHSvU9BGARKIc4ff5ynkI",
         "types" : [ "administrative_area_level_1", "political" ]
      },
      {
         "address_components" : [
            {
               "long_name" : "日本",
               "short_name" : "JP",
               "types" : [ "country", "political" ]
            }
         ],
         "formatted_address" : "日本",
         "geometry" : {
            "bounds" : {
               "northeast" : {
                  "lat" : 45.6412626,
                  "lng" : 154.0031455
               },
               "southwest" : {
                  "lat" : 20.3585295,
                  "lng" : 122.8554688
               }
            },
            "location" : {
               "lat" : 36.204824,
               "lng" : 138.252924
            },
            "location_type" : "APPROXIMATE",
            "viewport" : {
               "northeast" : {
                  "lat" : 45.6412626,
                  "lng" : 154.0031455
               },
               "southwest" : {
                  "lat" : 20.3585295,
                  "lng" : 122.8554688
               }
            }
         },
         "place_id" : "ChIJLxl_1w9OZzQRRFJmfNR1QvU",
         "types" : [ "country", "political" ]
      }
   ],
   "status" : "OK"
}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("access!")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	ctx := context.Background()
	latlong.Config.GoogleMapsAPIURL = server.URL

	l := latlong.NewPointISO6709([]byte("+34.455846+136.725739/"))

	latlong.Config.Lang = "ja"
	latlong.Config.GoogleAPIKey = "AIzaNotReallyAnAPIKey"
	latlong.Config.YahooJPClientID = "NotReallyAnAPIKey"

	s, err := l.Locality(ctx)
	if err != nil {
		t.Errorf("Locality returned non nil error: %v", err)
	}

	correctResponse := "伊勢市"
	if !reflect.DeepEqual(s, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, s)
	}
}
