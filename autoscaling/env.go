package autoscaling

import (
	"log"
	"os"
	"strconv"
)

const (
	// env inside gwMonitor container
	keyMonitorType         = "MONITOR_TYPE"
	keyGwHighThreshold     = "GW_CONN_NUMBER_HIGH_THRESHOLD"
	keyGwLowThreshold      = "GW_CONN_NUMBER_LOW_THRESHOLD"
	keyClientEndpoint      = "CLIENT_ENDPOINT"
	keyAddresses           = "ADDRESSES"
	keyPollingSeconds      = "POLLING_SECONDS"
	keyGwOverloadTolerance = "GW_OVERLOAD_TOLERANCE"
	keyGwIdleTolerance     = "GW_IDLE_TOLERANCE"
)

// Env is struct of environment variable
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
