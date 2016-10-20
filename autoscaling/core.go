package autoscaling

import "log"

const (
	alertHighGwConn = iota
	alertIdleGw
	alertNone
	alertError
	operationScaleOut
	operationScaleIn
	operationDoNothing

	// env in template app
	keyScaleInIP = "SCALE_IN_IP"
)

// judge compares 'realtime' statistic with theshold, and throw alert if overload
func analyseAlert(instances, connNum int, highThreshold int, allScaleInIPs []string) (int, error) {
	if instances == 0 {
		return alertNone, nil
	}
	realtimeAvgConn := float32(connNum) / float32(instances)

	// check if GW is overload
	log.Printf("I | realtimeAvgConn %f, highGwThreshold %d\n", realtimeAvgConn, highThreshold)
	if realtimeAvgConn > float32(highThreshold) {
		return alertHighGwConn, nil
	}
	// check if GW is idle
	if len(allScaleInIPs) > 0 {
		return alertIdleGw, nil
	}
	return alertNone, nil
}
