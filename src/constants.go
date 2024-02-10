package main

import "fmt"

type Constants struct {
	redisKeyPrefix string
	redisSummaryKey string
}

var constants Constants
func getConstants() Constants {
	constants = Constants {
		redisKeyPrefix: "wehatescrum",
	}

	constants.redisSummaryKey = fmt.Sprintf("%s:summary", constants.redisKeyPrefix)

	return constants
}