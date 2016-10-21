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
		LenIdleGWs    int
		ExpectedAlert int
	}{
		// UT cases provided by <jwang2@linkernetworks.com>, thanks for hard work.
		{0, 0, 0, 0, alertNone},
		{50, 1, 100, 0, alertNone},
		{50, 1, 100, 1, alertIdleGw},
		{110, 1, 100, 0, alertHighGwConn},
		{110, 1, 100, 1, alertHighGwConn},
		{100, 2, 100, 0, alertNone},
		{100, 2, 100, 1, alertIdleGw},
		{200, 2, 100, 0, alertNone},
		{200, 2, 100, 1, alertIdleGw},
		{300, 2, 100, 0, alertHighGwConn},
		{300, 2, 100, 1, alertHighGwConn},
		{500, 10, 100, 0, alertNone},
		{500, 10, 100, 2, alertIdleGw},
		{1100, 10, 100, 0, alertHighGwConn},
		{1100, 10, 100, 1, alertHighGwConn},
		{800, 16, 100, 0, alertNone},
		{800, 16, 100, 1, alertIdleGw},
		{1700, 16, 100, 0, alertHighGwConn},
		{1700, 16, 100, 1, alertHighGwConn},
	}

	for _, c := range cases {
		gotAlert, gotErr := analyseAlert(c.Instances, c.Connections, c.HighThreshold, c.LenIdleGWs)
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
		ExpectedAction string
		ExpectedGwIP   string
	}{
		// UT cases provided by <jwang2@linkernetworks.com>, thanks for hard work.
		{
			// liveGws <= 2
			[]string{},
			[]string{},
			tAllGWs,
			alertNone,
			actionNone,
			"",
		},
		{
			// liveGws <= 2
			[]string{"1"},
			[]string{},
			tAllGWs,
			alertNone,
			actionNone,
			"",
		},
		{
			// liveGws <= 2
			[]string{"2"},
			[]string{},
			tAllGWs,
			alertNone,
			actionNone,
			"",
		},
		{
			// liveGws <= 2
			[]string{"1", "2"},
			[]string{"2"},
			tAllGWs,
			alertIdleGw,
			actionNone,
			"",
		},
		{
			[]string{"1", "2", "3"},
			[]string{},
			tAllGWs,
			alertNone,
			actionNone,
			"",
		},
		{
			[]string{"1", "2", "3"},
			[]string{},
			tAllGWs,
			alertHighGwConn,
			actionAdd,
			"4",
		},
		{
			[]string{"1", "2", "3"},
			[]string{"2"},
			tAllGWs,
			alertHighGwConn,
			actionAdd,
			"4",
		},
		{
			[]string{"1", "2", "3"},
			[]string{"2"},
			tAllGWs,
			alertIdleGw,
			actionDel,
			"2",
		},
		{
			[]string{"1", "2", "5"},
			[]string{},
			tAllGWs,
			alertHighGwConn,
			actionAdd,
			"3",
		},
		{
			[]string{"1", "2", "5"},
			[]string{"2"},
			tAllGWs,
			alertHighGwConn,
			actionAdd,
			"3",
		},
		{
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			[]string{},
			tAllGWs,
			alertNone,
			actionNone,
			"",
		},
		{
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			[]string{},
			tAllGWs,
			alertHighGwConn,
			actionAdd,
			"11",
		},
		{
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			[]string{"3"},
			tAllGWs,
			alertHighGwConn,
			actionAdd,
			"11",
		},
		{
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			[]string{"2", "10"},
			tAllGWs,
			alertIdleGw,
			actionDel,
			"2",
		},
		{
			// unexpected alert
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			[]string{"2", "10"},
			tAllGWs,
			alertError,
			actionNone,
			"",
		},
		{
			// unexpected idle GWs
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			[]string{},
			tAllGWs,
			alertIdleGw,
			actionNone,
			"",
		},
		{
			// unknown alert
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			[]string{},
			tAllGWs,
			10,
			actionNone,
			"",
		},
		{
			// no more usable GW
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"},
			[]string{},
			tAllGWs,
			alertHighGwConn,
			actionNone,
			"",
		},
		{
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"},
			[]string{"7", "13"},
			tAllGWs,
			alertHighGwConn,
			actionNone,
			"",
		},
		{
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"},
			[]string{"7", "13"},
			tAllGWs,
			alertIdleGw,
			actionDel,
			"7",
		},
	}

	for _, c := range cases {
		gotDecision := makeDecision(c.LiveGWs, c.IdleGWs, c.AllGWs, c.Alert)
		assert.Equal(t, c.ExpectedAction, gotDecision.Action)
		assert.Equal(t, c.ExpectedGwIP, gotDecision.GwIP)
	}
}
