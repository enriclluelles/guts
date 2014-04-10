package main

import (
	"regexp"
	"strconv"
	"strings"
)

var sanitizer *regexp.Regexp = regexp.MustCompile("[^-a-zA-Z_0-9.]")

type MetricProcessor struct {
	metricStorage *MetricStorage
}

type Metric struct {
	value int64
}

func NewMetricProcessor(ms *MetricStorage) *MetricProcessor {
	return &MetricProcessor{ms}
}

func (mp *MetricProcessor) ParseAndProcess(packet *string) (err WrongMetricFormatError) {
	measurement, err := parseMeasurement(packet)

	mp.metricStorage.Store(&measurement)

	return err
}

func parseMeasurement(body *string) (measurement Measurement, err WrongMetricFormatError) {
	structure := strings.Split(*body, ":")

	length := len(structure)

	measurement.Name = sanitizeMetricName(structure[0])

	var parts []string

	if length > 2 {
		err = WrongMetricFormatError{}
		return
	} else if length > 1 {
		parts = strings.Split(structure[1], "|")
	} else {
		parts = []string{"1"}
	}

	valueLength := len(parts)

	var errInt error

	measurement.Value = parts[0]

	if errInt != nil {
		err = WrongMetricFormatError{}
		return
	}

	if valueLength > 1 {
		measurement.Type = parts[1]
	}

	if valueLength > 2 {
		val, errFloat := strconv.ParseFloat(parts[2], 64)
		if errFloat == nil {
			measurement.SampleRate = val
		} else {
			err = WrongMetricFormatError{}
			return
		}
	} else {
		measurement.SampleRate = 1.0
	}
	return
}

func sanitizeMetricName(metricName string) string {
	metricName = strings.Replace(metricName, " ", "_", -1)
	metricName = strings.Replace(metricName, "/", "-", -1)
	return sanitizer.ReplaceAllString(metricName, "")
}
