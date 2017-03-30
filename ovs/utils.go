package ovs

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	keyAddresses = "ADDRESSES"
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

	// log.Printf("Writing to server[%s], data: %s\n", addr, msg)
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
			// log.Println("Got data from address", remoteAddr, "message:", string(data[0:n]), n)
			info = string(data[0:n])
		}
	}

	return info, nil
}

func getAddrs() (addrs []string, err error) {
	strAddrs := os.Getenv(keyAddresses)
	if strings.EqualFold(strAddrs, "nil") {
		err = errors.New("getAddrs failed, find no addresses")
		return
	}
	//strAddrs looks like : "127.0.0.1:8080,127.0.0.1:8081"
	addrs = strings.Split(strAddrs, ",")
	return
}

// HostCount return count of OVS hosts
func HostCount() (n int) {
	addrs, _ := getAddrs()
	return len(addrs)
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

func parseJson(jsonstring string) (resp *Resp) {
	if len(jsonstring) == 0 {
		return
	}
	resp = &Resp{}
	err := json.Unmarshal([]byte(jsonstring), resp)
	if err != nil {
		log.Printf("unmarshal json \"%s\" error: %v\n", jsonstring, err)
		return
	}
	return
}

func process(allResp []Resp) (sumInstance int, sumConn int, allScaleInIPs []string, allLiveGWs []string) {

	sumInstance = 0
	sumConn = 0
	//get sumInstance and sumConn
	for _, resp := range allResp {
		sumInstance += resp.Instances
		sumConn += resp.ConnNum
		if strings.TrimSpace(resp.ScaleInIp) != "" {
			allScaleInIPs = append(allScaleInIPs, resp.ScaleInIp)
		}
		for _, liveGW := range resp.LiveGWs {
			allLiveGWs = append(allLiveGWs, liveGW)
		}
	}
	//monitorType, _ = getMonitorType()
	// log.Printf("sumInstance=%d, sumConn=%d, allScaleInIPs=%s, allLiveGWs=%v\n",
	// 	sumInstance, sumConn, allScaleInIPs, allLiveGWs)
	return
}
