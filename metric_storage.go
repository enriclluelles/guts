package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var addRegex *regexp.Regexp = regexp.MustCompile(`^(\+|-)\d`)

type MetricStore struct {
	gauges   map[string]int
	counters map[string]int
	timers   map[string]*Timer
	sets     map[string]*Set
}

type MetricStorage struct {
	MetricStore
	receivers    []chan MetricStore
	countersChan chan *Measurement
	gaugesChan   chan *Measurement
	timersChan   chan *Measurement
	setsChan     chan *Measurement
	inChannel    chan *Measurement
	storer       chan *Measurement
	pause        chan interface{}
}

func NewMetricStorage(receivers []chan MetricStore) *MetricStorage {
	m := new(MetricStorage)
	m.receivers = receivers
	m.counters = make(map[string]int)
	m.gauges = make(map[string]int)
	m.timers = make(map[string]*Timer)
	m.sets = make(map[string]*Set)
	m.countersChan = make(chan *Measurement)
	m.gaugesChan = make(chan *Measurement)
	m.timersChan = make(chan *Measurement)
	m.setsChan = make(chan *Measurement)
	m.storer = make(chan *Measurement)
	m.pause = make(chan interface{})
	go m.runCounters()
	go m.runSets()
	go m.runTimers()
	go m.runGauges()
	go m.store()
	go m.emit()
	return m
}

func (ms *MetricStore) reset() {
	oldCounters := ms.counters
	ms.counters = make(map[string]int)
	if !config.DeleteCounters {
		for i, _ := range oldCounters {
			ms.counters[i] = 0
		}
	}

	if config.DeleteGauges {
		ms.gauges = make(map[string]int)
	}

	oldTimers := ms.timers
	ms.timers = make(map[string]*Timer)
	if !config.DeleteTimers {
		for i, _ := range oldTimers {
			ms.timers[i] = NewTimer()
		}
	}

	ms.sets = make(map[string]*Set)
}

func (ms *MetricStorage) Store(m *Measurement) {
	ms.storer <- m
}

func (ms *MetricStorage) store() {
	for {
		select {
		case <-ms.pause:
			<-ms.pause
		case m := <-ms.storer:
			switch m.Type {
			case "c":
				ms.countersChan <- m
			case "g":
				ms.gaugesChan <- m
			case "ms":
				ms.timersChan <- m
			case "s":
				ms.setsChan <- m
			}
		}
	}
}

func (ms *MetricStorage) runCounters() {
	var m *Measurement
	for {
		m = <-ms.countersChan
		ms.counters[m.Name] = ms.counters[m.Name] + 1
	}
}

func (ms *MetricStorage) runSets() {
	var m *Measurement
	for {
		m = <-ms.setsChan
		_, ok := ms.sets[m.Name]

		if !ok {
			ms.sets[m.Name] = NewSet()
		}

		set := ms.sets[m.Name]
		(&set).Add(m)
	}
}

func (ms *MetricStorage) runTimers() {
	var m *Measurement
	for {
		m = <-ms.timersChan
		intValue, err := strconv.Atoi(m.Value)
		if err != nil {
			continue
		}
		timer, ok := ms.timers[m.Name]
		if !ok {
			fmt.Println("not ok")
			timer = NewTimer()
			ms.timers[m.Name] = timer
		}
		timer.Add(intValue)
	}
}

func (ms *MetricStorage) runGauges() {
	var m *Measurement
	for {
		m = <-ms.gaugesChan
		intValue, err := strconv.Atoi(m.Value)
		if err != nil {
			continue
		}
		if addRegex.MatchString(m.Value) {
			ms.gauges[m.Name] = ms.gauges[m.Name] + intValue
		} else {
			ms.gauges[m.Name] = intValue
		}
	}
}

func (ms *MetricStorage) reporter() {
	for {
		fmt.Printf("counters:\n")
		for name, value := range ms.counters {
			fmt.Printf("%s: %d\n", name, value)
		}
		fmt.Printf("gauges:\n")
		for name, value := range ms.gauges {
			fmt.Printf("%s: %d\n", name, value)
		}
		fmt.Printf("timers:\n")
		for name, value := range ms.timers {
			fmt.Println(name, ":", value)
		}
		fmt.Printf("sets:\n")
		for name, value := range ms.sets {
			fmt.Printf("%s: %#v\n", name, value)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func (ms *MetricStorage) emit() {
	time.Sleep(1 * time.Second)
	for {

		store := ms.MetricStore

		go func() {
			for _, timer := range store.timers {
				timer.setStats()
			}

			for _, receiver := range ms.receivers {
				receiver <- store
			}
		}()

		ms.MetricStore.reset()
		time.Sleep(time.Duration(config.FlushInterval) * time.Millisecond)
	}
}
