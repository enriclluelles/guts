package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
)

var config gutsConfig

type MetricServer struct {
	processor *MetricProcessor
	address   string
}

type Measurement struct {
	Name       string
	Type       string
	Value      string
	SampleRate float64
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) < 2 {
		log.Fatal("Please provide the config file as the first argument")
		os.Exit(1)
	} else {
		parseConfig(os.Args[1])
	}

	var receivers []chan MetricStore
	for _, backendType := range config.Backends {
		switch backendType {
		case "graphite":
			{
				receiver := make(chan MetricStore)
				receivers = append(receivers, receiver)
				NewGraphite(receiver)
			}
		}
	}

	metricStorage := NewMetricStorage(receivers)

	metricProcessor := NewMetricProcessor(metricStorage)
	server := NewMetricServer(metricProcessor, fmt.Sprintf("%s:%d", config.Address, config.Port))
	server.Start()
}

func NewMetricServer(mp *MetricProcessor, address string) *MetricServer {
	return &MetricServer{mp, address}
}

func (server *MetricServer) Start() {
	conn, err := net.ListenPacket("udp", server.address)
	fmt.Println("listening")
	buffer := make([]byte, 1024*8)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	for {
		n, _, err := conn.ReadFrom(buffer[0:])
		if err == nil {
			go server.ProcessPayload(string(buffer[0:n]))
		} else {
			log.Fatal(err)
		}
	}
}

func (server *MetricServer) ProcessPayload(payload string) {
	packets := strings.Split(payload, "\n")
	for _, packet := range packets {
		server.processor.ParseAndProcess(&packet)
	}
}
