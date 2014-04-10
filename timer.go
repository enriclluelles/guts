package main

import (
	"github.com/grd/stat"
	"sort"
)

type Timer struct {
	samples    []int
	mean       float64
	stdev      float64
	median     float64
	percentile float64
	sum        int
}

func NewTimer() *Timer {
	return &Timer{samples: make([]int, 0, 100)}
}

func (t *Timer) Add(measurement int) {
	t.samples = append(t.samples, measurement)
}

func (t *Timer) Get(i int) float64 {
	return float64(t.samples[i])
}

func (t *Timer) Len() int {
	return len(t.samples)
}

func (t *Timer) setStats() {
	sort.Sort(sort.IntSlice(t.samples))
	t.mean = stat.Mean(t)
	t.stdev = stat.SdMean(t, t.mean)
	t.median = stat.MedianFromSortedData(t)
	t.percentile = stat.QuantileFromSortedData(t, float64(config.PercentThreshold)/100)
	for _, sample := range t.samples {
		t.sum += sample
	}
}
