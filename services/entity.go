package services

type RespStruct struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Err     string      `json:"err"`
}

type RespData struct {
	Instances   int    `json:"instances"`
	ConnNum     int    `json:"connNum"`
	MonitorType string `json:"monitorType"`
	OvsId       int    `json:"ovsId"`
}
