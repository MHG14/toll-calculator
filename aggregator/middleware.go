package main

import (
	"time"

	"github.com/mhg14/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("aggregating distance")
	}(time.Now())

	err = m.next.AggregateDistance(distance)
	return
}

func (m *LogMiddleware) CalculateInvoice(id int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			totalDistance float64
			totalAmount   float64
		)
		if invoice != nil {
			totalAmount = invoice.TotalAmount
			totalDistance = invoice.TotalDistance
		}
		logrus.WithFields(logrus.Fields{
			"took":        time.Since(start),
			"err":         err,
			"obuID":       invoice.OBUID,
			"totalDist":   totalDistance,
			"totalAmount": totalAmount,
		}).Info("generating invoice")
	}(time.Now())

	invoice, err = m.next.CalculateInvoice(id)
	return
}
