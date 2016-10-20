package services

import (
	"log"
	"strings"
)

// GetInfos calls OVS API and returns processed data
func GetInfos() (instances, connNum int, monitorType string, allScaleInIPs []string, allLiveGWs []string, err error) {
	//get UDP server addresses from ENV file
	addrs, err := getAddrs()
	if err != nil {
		return
	}
	log.Println(addrs, len(addrs))

	monitorType, err = getMonitorType()
	if err != nil {
		log.Println(monitorType)
		return
	}
	log.Println("MONITOR_TYPE: ", monitorType)

	infos := make([]string, 0, len(addrs))

	//call UDP servers
	for _, address := range addrs {
		info, err := UdpCall(strings.TrimSpace(address), "hi")
		if err != nil {
			log.Println("UdpCall "+strings.TrimSpace(address)+" failed.", err)
		}
		info = strings.TrimSpace(info)
		infos = append(infos, info)
	}

	log.Println(infos)

	instances, connNum, allScaleInIPs, allLiveGWs = process(infos)
	return
}
