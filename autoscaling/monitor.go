package autoscaling

import (
	"bufio"
	"errors"
	"log"
	"os"
	"time"

	"github.com/LinkerNetworks/gwMonitor/conf"
	"github.com/LinkerNetworks/gwMonitor/services"
)

const (
	typePGW = "PGW"
	typeSGW = "SGW"

	alertHighPgwConn = iota
	alertHighSgwConn
	// alertLowPgwConn
	// alertLowSgwConn

)

var (
	pollingSeconds = conf.OptionsReady.PollingTime
	pollingTime    = time.Duration(pollingSeconds) * time.Second
	pgwTolerance   = 0
	sgwTolerance   = 0
)

// StartMonitor checks if an alert exists for a period <seconds>, and tigger autoscaling if it does.
func StartMonitor() {
	if env(keyMonitorDisable).ToBool() == true {
		log.Printf("monitor not enabled, set env %s to true to enable\n", keyMonitorDisable)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	rewind()
	for {
		time.Sleep(pollingTime)
		alert, err := analyse()
		if err != nil {
			log.Printf("E | analyse error: %v\n", err)
			continue
		}
		switch alert {
		case alertHighPgwConn:
			pgwTolerance--
			log.Printf("I | will scale PGW up in %ds\n", pgwTolerance)
		case alertHighSgwConn:
			sgwTolerance--
			log.Printf("I | will scale SGW up in %ds\n", sgwTolerance)
		default:
			// acts like a timer
			rewind()
		}
		if pgwTolerance <= 0 {
			rewind()
			// pgw overload
			log.Println("I | scaling up PGW instance...")
			scalePgwUp()
		}
		if sgwTolerance <= 0 {
			rewind()
			// sgw overload
			log.Println("I | scaling up SGW instance...")
			scaleSgwUp()
		}
	}
}

func rewind() {
	pgwTolerance = conf.OptionsReady.PgwTolerance
	sgwTolerance = conf.OptionsReady.SgwTolerance
}

// judge compares 'realtime' statistic with theshold, and throw alert if overload
func analyse() (int, error) {
	instances, connNum, monitorType, _, _, err := services.GetInfos()
	if err != nil {
		log.Printf("E | call service for data error: %v\n", err)
		return -1, err
	}
	log.Printf("I | got data: instances %d, connNum %d, monitorType %s\n", instances, connNum, monitorType)

	realtimeAvgConn := float32(connNum) / float32(instances)

	switch monitorType {
	case typePGW:
		// check if PGW is overload
		highPgwThreshold := env(keyPgwHighThreshold).ToInt()
		log.Printf("I | realtimeAvgConn %v, highPgwThreshold %d\n", realtimeAvgConn, highPgwThreshold)
		if realtimeAvgConn > float32(highPgwThreshold) {
			return alertHighPgwConn, nil
		}
	case typeSGW:
		// check if SGW is overload
		highSgwThreshold := env(keySgwHighThreshold).ToInt()
		log.Printf("I | realtimeAvgConn %v, highSgwThreshold %d\n", realtimeAvgConn, highSgwThreshold)
		if realtimeAvgConn > float32(highSgwThreshold) {
			return alertHighSgwConn, nil
		}
	default:
		log.Printf("E | unknow gateway type: %s\n", monitorType)
		return -2, errors.New("unknown gateway type")
	}
	return -3, nil
}
