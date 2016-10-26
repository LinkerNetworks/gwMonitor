package autoscaling

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestAnalyseAlert(t *testing.T) {
	var cases = []struct {
		Connections   int
		Instances     int
		HighThreshold int
		LowThreshold  int
		LenIdleGWs    int
		ExpectedAlert int
	}{
		// UT cases provided by <jwang2@linkernetworks.com>, thanks for hard work.
		{0, 0, 0, 0, 0, alertNone},
		{30, 1, 100, 50, 0, alertNone},
		{30, 1, 100, 50, 1, alertIdleGw},
		{50, 1, 100, 50, 0, alertNone},
		{50, 1, 100, 50, 1, alertNone},
		{80, 1, 100, 50, 0, alertNone},
		{80, 1, 100, 50, 1, alertNone},
		{100, 1, 100, 50, 0, alertNone},
		{100, 1, 100, 50, 1, alertNone},
		{120, 1, 100, 50, 0, alertHighGwConn},
		{120, 1, 100, 50, 1, alertHighGwConn},
		{701, 7, 100, 50, 0, alertHighGwConn},
		{701, 7, 100, 50, 1, alertHighGwConn},
		{349, 7, 100, 50, 0, alertNone},
		{349, 7, 100, 50, 1, alertIdleGw},
	}

	for _, c := range cases {
		gotAlert, gotErr := analyseAlert(c.Instances, c.Connections, c.HighThreshold, c.LowThreshold, c.LenIdleGWs)
		assert.Equal(t, c.ExpectedAlert, gotAlert)
		assert.Equal(t, nil, gotErr)
	}
}

// To make testing cases readable, use simplified string array of ["1","2","3"] instead of
// IPs like ["192.168.1.46","192.168.1.47","192.168.1.48"].
func TestMakeDecision(t *testing.T) {
	tAllGWs := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"}
	var cases = []struct {
		LiveGWs        []string
		IdleGWs        []string
		AllGWs         []string
		Alert          int
		HostCount      int
		ExpectedAction string
		ExpectedGwIP   string
	}{
		// UT cases provided by <jwang2@linkernetworks.com>, thanks for hard work.
		{[]string{}, []string{}, tAllGWs, alertNone, 0, actionNone, ""},
		{[]string{"1"}, []string{}, tAllGWs, alertNone, 1, actionNone, ""},
		{[]string{"1"}, []string{}, tAllGWs, alertHighGwConn, 1, actionAdd, "2"},
		{[]string{"1"}, []string{"1"}, tAllGWs, alertHighGwConn, 1, actionAdd, "2"},
		{[]string{"1"}, []string{}, tAllGWs, alertHighGwConn, 2, actionAdd, "2"},
		{[]string{"1"}, []string{"1"}, tAllGWs, alertIdleGw, 2, actionNone, ""},
		{[]string{"1", "2"}, []string{"1"}, tAllGWs, alertIdleGw, 2, actionNone, ""},
		{[]string{"1", "2"}, []string{}, tAllGWs, alertHighGwConn, 2, actionAdd, "3"},
		{[]string{"1", "2"}, []string{"2"}, tAllGWs, alertHighGwConn, 2, actionAdd, "3"},
		{[]string{"1", "2", "5"}, []string{}, tAllGWs, alertHighGwConn, 2, actionAdd, "3"},
		{[]string{"5", "2", "6"}, []string{}, tAllGWs, alertHighGwConn, 2, actionAdd, "1"},
		{[]string{"1", "2", "5"}, []string{"1"}, tAllGWs, alertIdleGw, 2, actionDel, "1"},
		{[]string{"1", "2", "5"}, []string{""}, tAllGWs, alertIdleGw, 2, actionNone, ""},
		{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, []string{"2", "5"}, tAllGWs, alertIdleGw, 2, actionDel, "2"},
		{[]string{"6", "8", "2", "3", "4", "9", "10"}, []string{"2", "9"}, tAllGWs, alertIdleGw, 2, actionDel, "2"},
		{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"}, []string{"2", "9"}, tAllGWs, alertIdleGw, 2, actionDel, "2"},
		{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"}, []string{}, tAllGWs, alertHighGwConn, 2, actionNone, ""},
		{[]string{"1", "2", "5"}, []string{""}, tAllGWs, 9999, 2, actionNone, ""},
		{[]string{"1", "2", "5"}, []string{""}, tAllGWs, alertError, 2, actionNone, ""},
		{[]string{"3", "2", "5"}, []string{"2"}, tAllGWs, alertIdleGw, 5, actionNone, ""},
		{[]string{"3", "2", "5", "6", "7"}, []string{"2"}, tAllGWs, alertIdleGw, 5, actionNone, ""},
		{[]string{"3", "2", "5", "6", "7"}, []string{}, tAllGWs, alertHighGwConn, 5, actionAdd, "1"},
		{[]string{"3", "2", "5", "6", "7", "4"}, []string{"4", "7"}, tAllGWs, alertIdleGw, 5, actionDel, "4"},
	}

	for _, c := range cases {
		gotDecision := makeDecision(c.LiveGWs, c.IdleGWs, c.AllGWs, c.Alert, c.HostCount)
		assert.Equal(t, c.ExpectedAction, gotDecision.Action)
		assert.Equal(t, c.ExpectedGwIP, gotDecision.GwIP)
	}
}
