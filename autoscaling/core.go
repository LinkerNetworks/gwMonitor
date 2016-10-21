package autoscaling

import "log"

const (
	alertHighGwConn = iota
	alertIdleGw
	alertNone
	alertError

	actionAdd  string = "add"
	actionDel  string = "del"
	actionNone string = "none"

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

// makeDecision decides what to do on which gw.
func makeDecision(liveGWs, idleGWs, allGWs []string, alert int) (decision Decision) {
	if len(liveGWs) < 2 {
		decision.Action = actionNone
		decision.Reason = "liveGws < 2"
		return
	}

	switch alert {
	case alertHighGwConn:
		gwAddIP := selectAddGw(liveGWs, allGWs)
		if len(gwAddIP) == 0 {
			decision.Action = actionNone
			decision.Reason = "no more usable GW"
			return
		}
		decision.Action = actionAdd
		decision.GwIP = gwAddIP
	case alertIdleGw:
		gwDelIP := selectDelGw(idleGWs)
		if len(gwDelIP) == 0 {
			decision.Action = actionNone
			decision.Reason = "unexpected idle GWs"
			return
		}
		decision.Action = actionDel
		decision.GwIP = gwDelIP
	case alertNone, alertError:
		decision.Action = actionNone
		decision.Reason = "unexpected alert"
	default:
		decision.Action = actionNone
		decision.Reason = "unknown alert"
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
