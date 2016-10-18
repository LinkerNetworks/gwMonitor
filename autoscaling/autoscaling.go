package autoscaling

import (
	"log"

	"github.com/LinkerNetworks/gwMonitor/conf"
)

func scalePgwUp() {
	pgwAppset := MinAppset{}
	// group, err := getPgwGroup()
	//
	// if err != nil {
	// 	log.Printf("get pgw group error: %v\n", err)
	// 	return
	// }

	scaleto := getPgwInstances() + conf.OptionsReady.PgwScaleStep

	if scaleto > lenTemplateApps() {
		log.Println("all template apps has started up, wont scale up")
		return
	}

	pgwAppset.Name = pgwGroupID
	pgwAppset.CreatedByJson = true
	pgwAppset.Group.ID = pgwGroupID
	pgwAppset.Group.Apps = getFirstNApps(scaleto)

	err := putAppset(pgwAppset)
	if err != nil {
		log.Printf("update pgw group error: %v\n", err)
	}
}

func scaleSgwUp() {
	sgwAppset := MinAppset{}
	// group, err := getSgwGroup()
	// if err != nil {
	// 	log.Printf("get sgw group error: %v\n", err)
	// 	return
	// }

	scaleto := getSgwInstances() + conf.OptionsReady.SgwScaleStep

	if scaleto > lenTemplateApps() {
		log.Println("all template apps has started up, wont scale up")
		return
	}

	sgwAppset.Name = sgwGroupID
	sgwAppset.CreatedByJson = true
	sgwAppset.Group.ID = sgwGroupID
	sgwAppset.Group.Apps = getFirstNApps(scaleto)

	err := putAppset(sgwAppset)
	if err != nil {
		log.Printf("update sgw group error: %v\n", err)
	}
}
