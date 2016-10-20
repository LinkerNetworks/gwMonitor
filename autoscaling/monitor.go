package autoscaling

import (
	"log"
	"os"
	"time"

	"github.com/LinkerNetworks/gwMonitor/conf"
	"github.com/LinkerNetworks/gwMonitor/ovs"
)

const (
	typePGW = "PGW"
	typeSGW = "SGW"
)

var (
	pollingSeconds      = conf.OptionsReady.PollingTime
	pollingTime         = time.Duration(pollingSeconds) * time.Second
	gwOverloadTolerance = 0
	gwIdleTolerance     = 0
)

// StartMonitor checks if GW is overload or idle for a period, and trigger scaling if it is.
func StartMonitor() {
	monitorType := env(keyMonitorType).Value
	switch monitorType {
	case typePGW:
		log.Println("starting PGW monitor daemon...")
		highThreshold := env(keyPgwHighThreshold).ToInt()
		if highThreshold <= 0 {
			log.Printf("invalid threshold, must set env %s\n", keyPgwHighThreshold)
			os.Exit(1)
		}
		startGwMonitorDaemon(highThreshold)
	case typeSGW:
		log.Println("starting SGW monitor daemon...")
		highThreshold := env(keySgwHighThreshold).ToInt()
		if highThreshold <= 0 {
			log.Printf("invalid threshold, must set env %s\n", keySgwHighThreshold)
			os.Exit(1)
		}
		startGwMonitorDaemon(highThreshold)
	default:
		log.Printf("unknown monitor type \"%s\", must set env %s\n", monitorType, keyMonitorType)
		os.Exit(1)
	}

}

func startGwMonitorDaemon(highGwThreshold int) {
	initDaemon()
	reqData := services.ReqData{}
	reqData.HighThreshold = string(highGwThreshold)
	for {
		time.Sleep(pollingTime)
		instances, connNum, gwType, allIdleGWs, allLiveGWs, err := services.CallOvsUDP(reqData)
		log.Printf("I | got data: instances %d, connNum %d, gwType %s, allIdleGWs %v, allLiveGWs %v\n",
			instances, connNum, gwType, allIdleGWs, allLiveGWs)
		if err != nil {
			log.Printf("E | call service for data error: %v\n", err)
			continue
		}
		alert, err := analyseAlert(instances, connNum, highGwThreshold, len(allIdleGWs))
		if err != nil {
			log.Printf("E | analyse error: %v\n", err)
			continue
		}
		switch alert {
		case alertHighGwConn:
			gwOverloadTolerance--
			log.Printf("I | will scale out GW in %ds\n", gwOverloadTolerance*pollingSeconds)
		case alertIdleGw:
			gwIdleTolerance--
			log.Printf("I | will scale in GW in %ds\n", gwIdleTolerance*pollingSeconds)
		default:
			rewindGwOverloadTimer()
			rewindGwIdleTimer()
		}
		if gwOverloadTolerance <= 0 {
			rewindGwOverloadTimer()
			// gateway overload for 60s(default)
			log.Println("I | scaling out GW instance...")
			operation := analyseOperation(allLiveGWs, allIdleGWs, allGwScaleIPs, alert)
			scaleGw(operation)
		}
		if gwIdleTolerance <= 0 {
			rewindGwIdleTimer()
			// gateway idle for 300s(default)
			log.Println("I | scaling in GW instance...")
			operation := analyseOperation(allLiveGWs, allIdleGWs, allGwScaleIPs, alert)
			scaleGw(operation)
			go notifyOvs(operation.GwIP)
		}
	}
}

func rewindGwOverloadTimer() {
	gwOverloadTolerance = conf.OptionsReady.GwOverloadTolerance
}

func rewindGwIdleTimer() {
	gwIdleTolerance = conf.OptionsReady.GwIdleTolerance
}

func initDaemon() {
	initTemplate()
	initScaling()
	rewindGwOverloadTimer()
	rewindGwIdleTimer()
}
