package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oschwald/geoip2-golang"
)

var h bool

// init global database variables for GeoIP
var GeoLite2City *geoip2.Reader

// Databases path (download from https://dev.maxmind.com/geoip/geoip2/geolite2/)
var DBCityPath = "./GeoLite2-City.mmdb"

func main() {
	var err error

	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&DBCityPath, "mmdb", "./GeoLite2-City.mmdb", "Databases path (download from https://dev.maxmind.com/geoip/geoip2/geolite2/)")
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	GeoLite2City, err = geoip2.Open(DBCityPath)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())

	/* pprof */
	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	/* get / */
	e.GET("/", func(c echo.Context) error {

		i, err := ip_info(c.RealIP())
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}

		i_be_f, err := json.MarshalIndent(i, "", "    ")
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}

		return c.JSONBlob(http.StatusOK, i_be_f)
	})

	err = e.Start(":80")
	if err != nil {
		panic(err)
	}
}

// ip_info get the detailed information of the ip address
func ip_info(ip string) (map[string]interface{}, error) {
	r := make(map[string]interface{})

	net_ip := net.ParseIP(ip)

	city_record, err := GeoLite2City.City(net_ip)
	if err != nil {
		return nil, err
	}

	r["IP Address"] = ip
	r["Portuguese (BR) city name"] = city_record.City.Names["en"]
	if len(city_record.Subdivisions) > 0 {
		r["Subdivision name"] = city_record.Subdivisions[0].Names["en"]
	}
	r["Country name"] = city_record.Country.Names["en"]
	r["ISO country code"] = city_record.Country.IsoCode
	r["Time zone"] = city_record.Location.TimeZone
	r["Coordinates"] = fmt.Sprintf("%v, %v", city_record.Location.Latitude, city_record.Location.Longitude)

	return r, nil
}
