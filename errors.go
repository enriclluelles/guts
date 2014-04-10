package main

import (
	"fmt"
)

type WrongMetricFormatError struct {
}

func (e WrongMetricFormatError) Error() string {
	return fmt.Sprintf("the metric has the wrong format")
}
