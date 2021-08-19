package geolite

import (
	"log"

	"github.com/oschwald/geoip2-golang"
)

var DB *geoip2.Reader

// Setup Initialize the geolite2 city instance
func Setup() error {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	return nil
}
