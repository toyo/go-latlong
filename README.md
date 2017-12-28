# go-latlong
[![Job Status](https://inspecode.rocro.com/badges/github.com/toyo/go-latlong/status?token=8_mXyA6LM_SqRuG8dCZ1ll-hW5zwNTIF5HetzaIZbOc)](https://inspecode.rocro.com/jobs/github.com/toyo/go-latlong/latest?completed=true)
[![Report](https://inspecode.rocro.com/badges/github.com/toyo/go-latlong/report?token=8_mXyA6LM_SqRuG8dCZ1ll-hW5zwNTIF5HetzaIZbOc&branch=master)](https://inspecode.rocro.com/reports/github.com/toyo/go-latlong/branch/master/summary)
[![Build Status](https://travis-ci.org/toyo/go-latlong.svg?branch=master)](https://travis-ci.org/toyo/go-latlong) [![GoDoc](https://godoc.org/github.com/toyo/go-latlong?status.svg)](https://godoc.org/github.com/toyo/go-latlong)

## Usage

This is an golang struct for store latitude, longitude and altitude by numerical, GridLocator and GeoHash.

To use the method of this struct, you can get GridLocator and GeoHash which length is considered in precision.


### Application
```go
package main

import (
	"fmt"

	"github.com/toyo/go-latlong"
)

func main() {
	l := latlong.NewLatLong(35, 135, 0.1, 0.1) // N35+-0.05 Deg. E135+-0.05 Deg.

	fmt.Println(l.GridLocator()) // shows GridLocator. https://en.wikipedia.org/wiki/Maidenhead_Locator_System
	fmt.Println(l.GeoHash())     // shows GeoHash. http://geohash.org/
	fmt.Println(l.String())      // shows lat/long in string.
}

```
