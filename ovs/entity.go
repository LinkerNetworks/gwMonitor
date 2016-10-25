package ovs

// Resp is structure of response from OVS by UDP
// * ScaleInIp is an unique mark of gateway, value is env of gateway when its' connNum==0
type Resp struct {
	Instances   int      `json:"instances"`
	ConnNum     int      `json:"connNum"`
	MonitorType string   `json:"monitorType"`
	OvsId       int      `json:"ovsId"`
	ScaleInIp   string   `json:"ScaleInIp"`
	LiveGWs     []string `json:"LiveGWs"`
}

// Req is structure of request sent to OVS by UPD
// * ScaleInIp is an unique mark of gateway, env from template
type Req struct {
	HighThreshold string `json:"HighThreshold,omitempty"`
	ScaleInIp     string `json:"ScaleInIp,omitempty"`
}
