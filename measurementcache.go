package sdm630

import (
	"fmt"
)

type MeasurementCache struct {
	datastream          ReadingChannel
	secsBetweenReadings int
	lastreading         Readings
}

func NewMeasurementCache(ds ReadingChannel, interval int) *MeasurementCache {
	return &MeasurementCache{datastream: ds, secsBetweenReadings: interval}
}

func (hc *MeasurementCache) ConsumeData() {
	for {
		reading := <-hc.datastream
		hc.lastreading = reading
		fmt.Printf("%s\r\n", &hc.lastreading)
	}
}

func (hc *MeasurementCache) GetLast() Readings {
	return hc.lastreading
}
