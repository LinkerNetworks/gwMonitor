package services

import (
	"encoding/json"
	"log"
	"strings"
)

// CallOvsUDP calls OVS API and returns processed data
func CallOvsUDP(reqData ReqData) (instances, connNum int, monitorType string, allScaleInIPs []string, allLiveGWs []string, err error) {
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

	data, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("json marshal reqData error: %v\n", err)
		return
	}

	//call UDP servers
	for _, address := range addrs {
		info, err := UdpCall(strings.TrimSpace(address), string(data))
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
