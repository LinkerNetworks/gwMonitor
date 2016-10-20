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

func scaleGw(operation Operation) {
	switch operation.Action {
	case actionAdd:
		scaleGwOut(operation.GwIP)
	case actionDel:
		scaleGwIn(operation.GwIP)
	default:
		log.Printf("unknown action \"%s\"\n", operation.Action)
	}
}

func scaleGwOut(gwAddIP string) (err error) {
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

func scaleGwIn(gwDelIP string) (err error) {
	appDel := getAppByEnv(keyScaleInIP, gwDelIP)

	err = delComponent(appDel.ID)
	if err != nil {
		log.Printf("delete component error: %v\n", err)
		return
	}
	return
}
