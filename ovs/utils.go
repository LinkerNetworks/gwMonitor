package services

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func UdpCall(server, msg string) (info string, err error) {
	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		log.Println("error: ", "Can't resolve address: ", err)
		return "", err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Println("error: ", "Can't dial: ", err)
		return "", err
	}
	conn.SetDeadline(time.Now().Add(3 * time.Second))

	defer conn.Close()

	log.Printf("Writing to server[%s], data: %s\n", addr, msg)
	_, err = conn.Write([]byte(msg))
	if err != nil {
		log.Println("error: ", "failed:", err)
		return "", err
	}

	data := make([]byte, 1024)

	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		log.Println("error: ", "failed to read UDP msg because of ", err)
		return "", err
	}

	if remoteAddr != nil {
		if n > 0 {
			log.Println("Got data from address", remoteAddr, "message:", string(data[0:n]), n)
			info = string(data[0:n])
		}
	}

	return info, nil
}

func getAddrs() (addrs []string, err error) {
	strAddrs := os.Getenv("ADDRESSES")
	if strings.EqualFold(strAddrs, "nil") {
		err = errors.New("getAddrs failed, find no addresses")
		return
	}
	//strAddrs looks like : "127.0.0.1:8080,127.0.0.1:8081"
	addrs = strings.Split(strAddrs, ",")
	return
}

func getMonitorType() (mtype string, err error) {

	//PGW or SGW
	mtype = os.Getenv("MONITOR_TYPE")

	if !strings.EqualFold(mtype, "PGW") && !strings.EqualFold(mtype, "SGW") {
		err = errors.New("getMonitorType failed, invalid MONITOR_TYPE")
		return
	}
	return
}

func parseJson(jsonstring string) (instances, connNum, ovsId int, scaleInIp string, liveGWs []string) {
	if len(jsonstring) == 0 {
		return
	}
	var respData = &RespData{}
	err := json.Unmarshal([]byte(jsonstring), respData)
	if err != nil {
		log.Printf("unmarshal json \"%s\" error: %v\n", jsonstring, err)
		return
	}
	return respData.Instances, respData.ConnNum, respData.OvsId, respData.ScaleInIp, respData.LiveGWs
}

func process(infos []string) (sumInstance int, sumConn int, allScaleInIPs []string, allLiveGWs []string) {

	sumInstance = 0
	sumConn = 0
	//get sumInstance and sumConn
	for _, info := range infos {
		instances, connNum, _, scaleInIp, liveGWs := parseJson(info)
		sumInstance = sumInstance + instances
		sumConn = sumConn + connNum
		if strings.TrimSpace(scaleInIp) != "" {
			allScaleInIPs = append(allScaleInIPs, scaleInIp)
		}
		for _, liveGW := range liveGWs {
			allLiveGWs = append(allLiveGWs, liveGW)
		}
	}
	//monitorType, _ = getMonitorType()
	log.Printf("sumInstance=%d, sumConn=%d, allScaleInIPs=%s, allLiveGWs=%v\n",
		sumInstance, sumConn, allScaleInIPs, allLiveGWs)
	return
}
