package autoscaling

import (
	"errors"
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

func scaleGw(action string, gwIP string) {
	switch action {
	case actionAdd:
		log.Printf("====> adding GW with scale ip: %s\n", gwIP)
		scaleGwOut(gwIP)
	case actionDel:
		log.Printf("====> deleting GW with scale ip: %s\n", gwIP)
		scaleGwIn(gwIP)
	default:
		log.Printf("unknown action \"%s\"\n", action)
	}
}

func scaleGwOut(gwAddIP string) (err error) {
	appAdd := getAppByEnv(keyScaleInIP, gwAddIP)
	if appAdd == nil {
		err = errors.New("app not found")
		log.Printf("app not found by env key %s, value %s\n", keyScaleInIP, gwAddIP)
		return
	}

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

func scaleGwIn(gwDelIP string) (err error) {
	appDel := getAppByEnv(keyScaleInIP, gwDelIP)

	if appDel == nil {
		err = errors.New("app not found")
		log.Printf("app not found by env key %s, value %s\n", keyScaleInIP, gwDelIP)
		return
	}

	err = delComponent(appDel.ID)
	if err != nil {
		log.Printf("delete component error: %v\n", err)
		return
	}
	return
}
