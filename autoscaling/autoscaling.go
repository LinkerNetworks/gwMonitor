package autoscaling

import (
	"log"
	"strings"
)

var (
	allGwScaleIPs []string
)

func initScaling() {
	for _, app := range gwGroup.Apps {
		env := *app.Env
		allGwScaleIPs = append(allGwScaleIPs, env[keyScaleInIP])
	}
}

func scaleGwOut(allLiveGWs []string) (err error) {
	gwAddIP := selectAddGw(allLiveGWs)
	if len(gwAddIP) == 0 {
		log.Println("gwAddIP is blank")
		return
	}

	appAdd := getAppByEnv(keyScaleInIP, gwAddIP)

	c := &MinComponent{}
	c.App = *appAdd
	c.AppsetName = strings.TrimLeft(gwGroup.ID, "/")

	err = addComponent(c)
	if err != nil {
		log.Printf("add component[appID: %s] error: %v\n", c.App.ID, err)
		return
	}

	err = startComponent(c.App.ID)
	if err != nil {
		log.Printf("start component[appID: %s] error: %v\n", c.App.ID, err)
	}
	return
}

func scaleGwIn(allScaleInIPs []string) (err error) {
	gwDelIP := selectDelGw(allScaleInIPs)
	if len(gwDelIP) == 0 {
		log.Println("gwDelIP is blank")
		return
	}

	appDel := getAppByEnv(keyScaleInIP, gwDelIP)

	err = delComponent(appDel.ID)
	if err != nil {
		log.Printf("delete component error: %v\n", err)
		return
	}
	return
}

// select gateway to add
func selectAddGw(allLiveGWs []string) (gwAddIP string) {
	var usableGWs []string
	for _, gw := range allGwScaleIPs {
		if !stringInSlice(gw, allLiveGWs) {
			usableGWs = append(usableGWs, gw)
		}
	}
	if len(usableGWs) >= 1 {
		return usableGWs[0]
	}
	return
}

// select gateway to remove
func selectDelGw(allScaleInIPs []string) (gwDelIP string) {
	if len(allScaleInIPs) >= 1 {
		return allScaleInIPs[0]
	}
	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
