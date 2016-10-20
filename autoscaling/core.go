package autoscaling

import "log"

const (
	alertHighGwConn = iota
	alertIdleGw
	alertNone
	alertError

	actionAdd string = "add"
	actionDel string = "del"

	// env in template app
	keyScaleInIP = "SCALE_IN_IP"
)

// judge compares 'realtime' statistic with theshold, and throw alert if overload
func analyseAlert(instances, connections int, highThreshold int, lenIdleGWs int) (int, error) {
	if instances == 0 {
		return alertNone, nil
	}
	realtimeAvgConn := float32(connections) / float32(instances)

	// check if GW is overload
	log.Printf("I | realtimeAvgConn %f, highGwThreshold %d\n", realtimeAvgConn, highThreshold)
	if realtimeAvgConn > float32(highThreshold) {
		return alertHighGwConn, nil
	}
	// check if GW is idle
	if lenIdleGWs > 0 {
		return alertIdleGw, nil
	}
	return alertNone, nil
}

// analyseOperation decides what to do on which gw.
func analyseOperation(liveGWs, idleGWs, allGWs []string, alert int) (operation Operation) {
	switch alert {
	case alertHighGwConn:
		gwAddIP := selectAddGw(liveGWs, allGWs)
		operation.Action = actionAdd
		operation.GwIP = gwAddIP
	case alertIdleGw:
		gwDelIP := selectDelGw(idleGWs)
		operation.Action = actionDel
		operation.GwIP = gwDelIP
	case alertNone:
		// do nothing
	case alertError:
		log.Println("received alertError")
	default:
		log.Printf("unknown alert %d\n", alert)
	}
	return
}

// select gateway to add
// allGWs: array of all env ScaleIP from template
func selectAddGw(liveGWs, allGWs []string) (gwAddIP string) {
	var usableGWs []string
	for _, gw := range allGWs {
		if !stringInSlice(gw, liveGWs) {
			usableGWs = append(usableGWs, gw)
		}
	}
	if len(usableGWs) >= 1 {
		return usableGWs[0]
	}
	return
}

// select gateway to remove
func selectDelGw(idleGWs []string) (gwDelIP string) {
	if len(idleGWs) >= 1 {
		return idleGWs[0]
	}
	return
}
