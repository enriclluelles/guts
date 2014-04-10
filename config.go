package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type gutsConfig struct {
	Port             int
	Backends         []string
	FlushInterval    int
	Graphite         graphiteConfig
	GraphiteHost     string
	GraphitePort     int
	Address          string
	DeleteCounters   bool
	DeleteSets       bool
	DeleteTimers     bool
	DeleteGauges     bool
	PercentThreshold int
}

type graphiteConfig struct {
	LegacyNamespace bool
	GlobalPrefix    string
	PrefixCounter   string
	PrefixTimer     string
	PrefixGauge     string
	PrefixSet       string
	GlobalSuffix    string
}

func parseConfig(fileName string) {
	if rawConfig, err := ioutil.ReadFile(os.Args[1]); err != nil {
		log.Fatal("Couldn't read the config file: ", err)
		os.Exit(1)
	} else {
		if err = json.Unmarshal(rawConfig, &config); err != nil {
			log.Fatal("Couldn't read the config file: ", err.Error())
			os.Exit(1)
		}
	}

	setConfigDefaults()
}

func setConfigDefaults() {
	if config.FlushInterval == 0 {
		config.FlushInterval = 10000
	}

	if config.Port == 0 {
		config.Port = 8125
	}

	if config.Address == "" {
		config.Address = "127.0.0.1"
	}

	if config.Graphite.PrefixSet == "" {
		config.Graphite.PrefixSet = "sets"
	}

	if config.Graphite.PrefixCounter == "" {
		config.Graphite.PrefixCounter = "counters"
	}

	if config.Graphite.PrefixTimer == "" {
		config.Graphite.PrefixTimer = "timers"
	}

	if config.Graphite.PrefixGauge == "" {
		config.Graphite.PrefixGauge = "gauges"
	}

	if config.Graphite.GlobalPrefix == "" {
		config.Graphite.GlobalPrefix = "stats"
	}

	if config.PercentThreshold == 0 {
		config.PercentThreshold = 90
	}

}
