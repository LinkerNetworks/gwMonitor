package autoscaling

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/LinkerNetworks/gwMonitor/conf"
	marathon "github.com/gambol99/go-marathon"
)

var (
	client       marathon.Marathon
	onceMarathon sync.Once
)

func init() {
	onceMarathon.Do(func() {
		config := marathon.NewDefaultConfig()
		if env(keyMarathonURL).Value != "" {
			config.URL = env(keyMarathonURL).Value
		} else {
			config.URL = conf.OptionsReady.MarathonURL
		}
		config.PollingWaitTime = 500 * time.Millisecond
		config.EventsPort = 50001

		c, err := marathon.NewClient(config)
		if err != nil {
			log.Printf("new marathon client error: %v\n", err)
		}
		client = c
	})
}

func getPgwGroup() (group *marathon.Group, err error) {
	return getGroup(pgwGroupID)
}

func getSgwGroup() (group *marathon.Group, err error) {
	return getGroup(sgwGroupID)
}

func getGroup(id string) (group *marathon.Group, err error) {
	group, err = client.Group(id)
	if err != nil {
		fmt.Printf("get group error: %v\n", err)
		return
	}
	return
}

func updateGroup(group *marathon.Group) (err error) {
	deployment, err := client.UpdateGroup(group.ID, group, false)
	if err != nil {
		log.Printf("update group error: %v\n", err)
		return
	}

	err = client.WaitOnDeployment(deployment.DeploymentID, 60)
	if err != nil {
		log.Printf("wait on deployment error: %v\n", err)
		return
	}
	return
}
