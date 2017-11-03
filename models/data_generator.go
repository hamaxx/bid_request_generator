package models

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	winDelayDist = 100
	winDelayMax  = 5 * time.Minute

	clickDelayDist = 100
	clickDelayMax  = time.Hour

	bidPriceDist = 10
	bidPriceMax  = 20.0

	adUrlDist = 10
	adUrlMax  = 10000

	publisherDomainDist = 10
	publisherDomainMax  = 1000
	publisherPageDist   = 10
	publisherPageMax    = 1000

	geoDist = 100
	zipMax  = 10000
)

func GetWinNoticeTimeDiff() time.Duration {
	return time.Duration(getExpRand(winDelayDist, float64(winDelayMax)))
}

func GetClickTimeDiff() time.Duration {
	return time.Duration(getExpRand(clickDelayDist, float64(clickDelayMax)))
}

func GetBidPrice() float64 {
	return math.Floor(getExpRand(bidPriceDist, bidPriceMax)*1000.0) / 1000.0
}

func GetAdUrl() string {
	adId := int(math.Floor(getExpRand(adUrlDist, adUrlMax)))
	return fmt.Sprintf("https://ad.zemanta.com/%d", adId)
}

func GetPublisherUrl() string {
	domain := int(math.Floor(getExpRand(publisherDomainDist, publisherDomainMax)))
	page := int(math.Floor(getExpRand(publisherPageDist, publisherPageMax)))

	return fmt.Sprintf("https://%d.example.com/%d.html", domain, page)
}

func GetGeoTargeting() *TargetingGeo {
	g := &TargetingGeo{}

	countryIdx := int(math.Floor(getExpRand(geoDist, float64(len(Countries)))))
	g.Country = Countries[countryIdx]

	regions := Regions[g.Country]
	if regions != nil {
		regionIdx := int(math.Floor(getExpRand(geoDist, float64(len(regions)))))
		g.Region = regions[regionIdx]
	}

	if g.Country == "US" {
		g.Zip = fmt.Sprint(int(math.Floor(getExpRand(geoDist, zipMax))))
	}

	return g
}

func GetDeviceTargeting() *TargetingDevice {
	d := &TargetingDevice{}

	d.Type = DeviceTypes[rand.Intn(len(DeviceTypes))]
	d.Os = DeviceOs[d.Type][rand.Intn(len(DeviceOs[d.Type]))]

	return d
}

func getExpRand(rate float64, max float64) float64 {
	for i := 0; i < 10; i++ {
		r := rand.ExpFloat64() * max / rate
		if r > max {
			continue
		}
		return r
	}
	return max
}
