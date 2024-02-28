package main

import (
	"fmt"
	"time"
)

type Constants struct {
	redisKeyPrefix string
	redisSummaryKey string
	timeLayout string
	dateLayout string
}

var constants Constants
func getConstants() Constants {
	constants = Constants {
		redisKeyPrefix: "wehatescrum",
		timeLayout: "2006-01-02T15:04:05Z07:00",
		dateLayout: "2006-01-02",
	}

	constants.redisSummaryKey = fmt.Sprintf("%s:summary", constants.redisKeyPrefix)

	return constants
}

func getRedisSummaryKey(opts SummaryOpts) string {
	c := getConstants()

	var redisDurationKey string
	if opts.duration > 100_00 * time.Hour {
		redisDurationKey = "ALL"
	} else {
		redisDurationKey = opts.duration.String()
	}
	
	return fmt.Sprintf("%s:%s:%s", c.redisSummaryKey, opts.date.Format(c.dateLayout), redisDurationKey)
}