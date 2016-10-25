package ovs

import (
	"encoding/json"
	"log"
	"strings"
)

// CallOvsUDP calls OVS API and returns processed data
func CallOvsUDP(req Req) (instances, connNum int, monitorType string, allScaleInIPs []string, allLiveGWs []string, err error) {
	//get UDP server addresses from ENV file
	addrs, err := getAddrs()
	if err != nil {
		return
	}

	monitorType, err = getMonitorType()
	if err != nil {
		log.Println(monitorType)
		return
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		log.Printf("json marshal reqData error: %v\n", err)
		return
	}

	allResp := make([]Resp, 0, len(addrs))
	//call UDP servers
	for _, address := range addrs {
		respData, err := UdpCall(strings.TrimSpace(address), string(reqData))
		if err != nil {
			log.Println("UdpCall "+strings.TrimSpace(address)+" failed.", err)
		}
		resp := parseJson(respData)
		if resp == nil {
			continue
		}
		allResp = append(allResp, *resp)
	}

	instances, connNum, allScaleInIPs, allLiveGWs = process(allResp)
	return
}
