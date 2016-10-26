package autoscaling

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	verifyEnv()
	highThreshold := env(keyGwHighThreshold).ToInt()
	lowThreshold := env(keyGwLowThreshold).ToInt()
	log.Println("I | starting GW monitor daemon...")
	if highThreshold <= 0 {
		log.Printf("E | invalid threshold, check env %s\n", keyGwHighThreshold)
		os.Exit(1)
	}
	if lowThreshold <= 0 {
		log.Printf("E | invalid threshold, check env %s\n", keyGwLowThreshold)
		os.Exit(1)
	}
	startGwMonitorDaemon(highThreshold, lowThreshold)
}

func startGwMonitorDaemon(highGwThreshold, lowThreshold int) {
	initDaemon()
	req := ovs.Req{}
	req.HighThreshold = strconv.Itoa(highGwThreshold)
	for {
		time.Sleep(pollingTime)
		instances, connNum, _, allIdleGWs, allLiveGWs, err := ovs.CallOvsUDP(req)
		log.Printf("I | got data: instances %d, connNum %d, allIdleGWs %v, allLiveGWs %v\n",
			instances, connNum, allIdleGWs, allLiveGWs)
		if err != nil {
			log.Printf("E | call service for data error: %v\n", err)
			continue
		}
		alert, err := analyseAlert(instances, connNum, highGwThreshold, lowThreshold, len(allIdleGWs))
		if err != nil {
			log.Printf("E | analyse error: %v\n", err)
			continue
		}
		switch alert {
		case alertHighGwConn:
			gwOverloadTolerance--
			log.Printf("I | will consider scaling out GW in %ds\n", gwOverloadTolerance*pollingSeconds)
		case alertIdleGw:
			gwIdleTolerance--
			log.Printf("I | will consider scaling in GW in %ds\n", gwIdleTolerance*pollingSeconds)
		default:
			rewindGwOverloadTimer()
			rewindGwIdleTimer()
		}
		if gwOverloadTolerance <= 0 {
			rewindGwOverloadTimer()
			// gateway overload for 60s(default)
			decision := makeDecision(allLiveGWs, allIdleGWs, allGwScaleIPs, alert, ovs.HostCount())
			log.Printf("I | figured out decision %v...\n", decision)
			if decision.Action == actionNone {
				log.Printf("I | wont scale out because \"%s\"\n", decision.Reason)
				continue
			}
			scaleGw(decision.Action, decision.GwIP)
		}
		if gwIdleTolerance <= 0 {
			rewindGwIdleTimer()
			// gateway idle for 300s(default)
			decision := makeDecision(allLiveGWs, allIdleGWs, allGwScaleIPs, alert, ovs.HostCount())
			log.Printf("I | figured out operation %v...\n", decision)
			if decision.Action == actionNone {
				log.Printf("I | wont scale in because \"%s\"\n", decision.Reason)
				continue
			}
			scaleGw(decision.Action, decision.GwIP)
			go notifyOvs(decision.GwIP)
		}
		fmt.Println("")
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
