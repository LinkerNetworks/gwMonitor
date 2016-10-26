package autoscaling

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/LinkerNetworks/gwMonitor/conf"
	marathon "github.com/gambol99/go-marathon"
)

var (
	gwGroup *marathon.Group
)

func initTemplate() {
	monitorType := env(keyMonitorType).Value
	switch monitorType {
	case typePGW:
		jsonPath := conf.OptionsReady.PgwJSON
		gwGroup = loadTemplate(jsonPath)
	case typeSGW:
		jsonPath := conf.OptionsReady.SgwJSON
		gwGroup = loadTemplate(jsonPath)
	default:
		log.Printf("unknow monitor type: %s\n", monitorType)
	}
	verifyJSON()
}

func getEth1Ip(scaleInIp string) (eth1Ip string) {
	app := getAppByEnv(keyScaleInIP, scaleInIp)
	envMap := *app.Env
	return envMap[keyEth1]
}

func getEth2Ip(scaleInIp string) (eth2Ip string) {
	app := getAppByEnv(keyScaleInIP, scaleInIp)
	envMap := *app.Env
	return envMap[keyEth2]
}

func getAppByEnv(key string, value string) (app *marathon.Application) {
	for _, app := range gwGroup.Apps {
		envMap := *app.Env
		if envMap[key] == value {
			return app
		}
	}
	return
}

func loadTemplate(jsonPath string) (group *marathon.Group) {
	content, err := readTextFile(jsonPath)
	if err != nil {
		return
	}
	group, err = parseJSON(content)
	if err != nil {
		return
	}
	return
}

func parseJSON(content []byte) (group *marathon.Group, err error) {
	group = &marathon.Group{}
	err = json.Unmarshal(content, group)
	if err != nil {
		log.Printf("unmarshal to json error: %v\n", err)
		return
	}
	return
}

func readTextFile(path string) (content []byte, err error) {
	if _, err = os.Stat(path); err != nil {
		log.Printf("stat file error: %v\n", err)
		return
	}
	content, err = ioutil.ReadFile(path)
	if err != nil {
		log.Printf("read file error: %v\n", err)
		return
	}
	return content, nil
}
