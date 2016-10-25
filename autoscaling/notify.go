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
	req.ScaleInIp = scaleInIP
	_, _, _, _, aliveGWs, err = ovs.CallOvsUDP(req)
	if err != nil {
		log.Printf("call ovs udp error: %v\n", err)
	}
	return
}
