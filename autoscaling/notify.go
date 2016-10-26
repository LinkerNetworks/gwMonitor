package autoscaling

import (
	"log"
	"time"

	"github.com/LinkerNetworks/gwMonitor/ovs"
)

func notifyOvs(removedGwIP string) {
	aliveGWs, _ := send(removedGwIP)
	// retry if package lost
	for i := 0; i < 4; i++ {
		if stringInSlice(removedGwIP, aliveGWs) {
			time.Sleep(1 * time.Second)
			aliveGWs, _ = send(removedGwIP)
		}
	}
}

func send(scaleInIP string) (aliveGWs []string, err error) {
	req := ovs.Req{}
	eth1Ip := getEth1Ip(scaleInIP)

	monitorType := env(keyMonitorType).Value
	switch monitorType {
	case typePGW:
		req.ScaleInIp = eth1Ip
	case typeSGW:
		eth2Ip := getEth2Ip(scaleInIP)
		req.ScaleInIp = eth1Ip + "," + eth2Ip
	default:
		log.Printf("unknow monitor type: %s\n", monitorType)
	}

	_, _, _, _, aliveGWs, err = ovs.CallOvsUDP(req)
	if err != nil {
		log.Printf("call ovs udp error: %v\n", err)
	}
	return
}
