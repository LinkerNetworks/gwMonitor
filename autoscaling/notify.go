package autoscaling

import (
	"log"
	"time"

	"github.com/LinkerNetworks/gwMonitor/services"
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

func send(scaleInIp string) (aliveGWs []string, err error) {
	reqData := services.ReqData{}
	reqData.ScaleInIp = scaleInIp
	_, _, _, _, aliveGWs, err = services.CallOvsUDP(reqData)
	if err != nil {
		log.Printf("call ovs udp error: %v\n", err)
	}
	return
}
