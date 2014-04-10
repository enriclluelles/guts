package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Graphite struct {
	receiver          chan MetricStore
	countersNamespace []string
	gaugesNamespace   []string
	timersNamespace   []string
	setsNamespace     []string
}

func NewGraphite(receiver chan MetricStore) *Graphite {
	g := Graphite{receiver: receiver}
	g.setupNamespaces()
	log.Println("starting graphite backend")
	go g.start()
	return &g
}

func (g *Graphite) setupNamespaces() {
	cfg := config.Graphite
	g.countersNamespace = []string{}
	g.gaugesNamespace = []string{}
	g.timersNamespace = []string{}
	g.setsNamespace = []string{}
	if cfg.LegacyNamespace {
		global := "stats"
		g.countersNamespace = append(g.countersNamespace, global)
		g.gaugesNamespace = append(g.gaugesNamespace, global, "gauges")
		g.timersNamespace = append(g.timersNamespace, global, "timers")
		g.setsNamespace = append(g.setsNamespace, global, "sets")
	} else {
		global := cfg.GlobalPrefix
		g.countersNamespace = append(g.countersNamespace, global, cfg.PrefixCounter)
		g.gaugesNamespace = append(g.gaugesNamespace, global, cfg.PrefixGauge)
		g.timersNamespace = append(g.timersNamespace, global, cfg.PrefixTimer)
		g.setsNamespace = append(g.setsNamespace, global, cfg.PrefixSet)
	}
}

func (graphite *Graphite) start() {
	for {
		store := <-graphite.receiver
		timestamp := time.Now().Unix()
		var buffer bytes.Buffer
		buffer.Write(graphite.countersData(&store, timestamp))
		buffer.Write(graphite.timersData(&store, timestamp))
		graphite.send(buffer.Bytes())
	}
}

func (graphite *Graphite) send(data []byte) {
	address := fmt.Sprintf("%s:%d", config.GraphiteHost, config.GraphitePort)
	client, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Couldn't connect to graphite")
	} else {
		client.Write(data)
		client.Close()
	}
}

func (graphite *Graphite) countersData(ms *MetricStore, timestamp int64) []byte {
	namespace := strings.Join(graphite.countersNamespace, ".")
	var buffer bytes.Buffer

	for key, counter := range ms.counters {
		buffer.WriteString(fmt.Sprintf("%s.%s %d %d\n", namespace, key, counter, timestamp))
	}

	return buffer.Bytes()
}

func (graphite *Graphite) timersData(ms *MetricStore, timestamp int64) []byte {
	namespace := strings.Join(graphite.timersNamespace, ".")
	var buffer bytes.Buffer

	for key, timer := range ms.timers {
		fmt.Fprintf(&buffer, "%s.%s.count %d %d\n", namespace, key, timer.Len(), timestamp)
		fmt.Fprintf(&buffer, "%s.%s.mean %f %d\n", namespace, key, timer.mean, timestamp)
		fmt.Fprintf(&buffer, "%s.%s.median %f %d\n", namespace, key, timer.median, timestamp)
		fmt.Fprintf(&buffer, "%s.%s.upper_%d %f %d\n", namespace, key, config.PercentThreshold, timer.percentile, timestamp)
		fmt.Fprintf(&buffer, "%s.%s.sum %d %d\n", namespace, key, timer.sum, timestamp)
	}

	return buffer.Bytes()
}
