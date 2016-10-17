package autoscaling

import (
	"log"
	"os"
	"strconv"
)

const (
	keyMonitorDisable   = "MONITOR_DISABLE"
	keyMonitorType      = "MONITOR_TYPE"
	keyPgwHighThreshold = "PGW_CONN_NUMBER_HIGH_THRESHOLD"
	keySgwHighThreshold = "SGW_CONN_NUMBER_HIGH_THRESHOLD"
	// keyPgwLowThreshold  = "PGW_CONN_NUMBER_LOW_THRESHOLD"
	// keySgwLowThreshold  = "SGW_CONN_NUMBER_LOW_THRESHOLD"
	keyMarathonURL = "MARATHON_URL"
)

// Env is value of environment variable
type Env struct {
	Key   string
	Value string
}

func env(key string) (e *Env) {
	return &Env{Key: key, Value: os.Getenv(key)}
}

// ToInt parse value to int
func (e *Env) ToInt() int {
	i, err := strconv.ParseInt(e.Value, 10, 0)
	if err != nil {
		log.Printf("parse env %s to int error: %v\n", e.Key, err)
		return -1
	}
	return int(i)
}

// ToBool parse value to bool
func (e *Env) ToBool() bool {
	b, err := strconv.ParseBool(e.Value)
	if err != nil {
		log.Printf("parse env %s to bool error: %v\n", e.Key, err)
		return false
	}
	return b
}
