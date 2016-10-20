package autoscaling

import (
	"bufio"
	"fmt"
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

// StartMonitor checks if an alert exists for a period <seconds>, and tigger autoscaling if it does.
func StartMonitor() {
	if env(keyMonitorDisable).ToBool() == true {
		log.Printf("monitor not enabled, set env %s to true to enable\n", keyMonitorDisable)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	monitorType := env(keyMonitorType).Value
	switch monitorType {
	case typePGW:
		log.Println("starting PGW monitor daemon...")
		highThreshold := env(keyPgwHighThreshold).ToInt()
		startGwMonitorDaemon(highThreshold)
	case typeSGW:
		log.Println("starting SGW monitor daemon...")
		highThreshold := env(keySgwHighThreshold).ToInt()
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
		instances, connNum, gwType, allScaleInIPs, allLiveGWs, err := services.CallOvsUDP(reqData)
		log.Printf("I | got data: instances %d, connNum %d, gwType %s, allScaleInIPs %v, allLiveGWs %v\n",
			instances, connNum, gwType, allScaleInIPs, allLiveGWs)
		if err != nil {
			log.Printf("E | call service for data error: %v\n", err)
			continue
		}
		alert, err := analyse(instances, connNum, highGwThreshold, allScaleInIPs)
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
			// acts like a timer
			rewindGwOverloadTimer()
			rewindGwIdleTimer()
		}
		if gwOverloadTolerance <= 0 {
			rewindGwOverloadTimer()
			// gateway overload for 60s
			log.Println("I | scaling out GW instance...")
			gwAddIP := selectAddGw(allLiveGWs)
			if len(gwAddIP) == 0 {
				log.Println("gwAddIP is blank")
				continue
			}
			scaleGwOut(gwAddIP)
		}
		if gwIdleTolerance <= 0 {
			rewindGwIdleTimer()
			// gateway idle for 300s
			log.Println("I | scaling in GW instance...")
			gwDelIP := selectDelGw(allScaleInIPs)
			if len(gwDelIP) == 0 {
				log.Println("gwDelIP is blank")
				return
			}
			scaleGwIn(gwDelIP)
			go notifyOvs(gwDelIP)
		}
	}
}

func rewindGwOverloadTimer() {
	gwOverloadTolerance = conf.OptionsReady.GwOverloadTolerance
	fmt.Println(conf.OptionsReady.GwOverloadTolerance)
}

func rewindGwIdleTimer() {
	gwIdleTolerance = conf.OptionsReady.GwIdleTolerance
	fmt.Println(conf.OptionsReady.GwIdleTolerance)
}

func initDaemon() {
	initTemplate()
	initScaling()
	rewindGwOverloadTimer()
	rewindGwIdleTimer()
}
