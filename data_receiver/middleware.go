package main

import (
	"time"
	types "tolling/Types"

	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat":   data.Lat,
			"lon":   data.Lon,
			"took":  time.Since(start),
		}).Info("Producing to kafka")
	}(time.Now())

	return l.next.ProduceData(data)
}
