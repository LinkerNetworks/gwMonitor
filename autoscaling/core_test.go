package autoscaling

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestAnalyseAlert(t *testing.T) {
	var cases = []struct {
		Instances     int
		Connections   int
		HighThreshold int
		LenIdleGWs    int
		ExpectedAlert int
	}{
		{2, 100, 100, 1, alertIdleGw},
	}

	for _, c := range cases {
		gotAlert, gotErr := analyseAlert(c.Instances, c.Connections, c.HighThreshold, c.LenIdleGWs)
		assert.Equal(t, c.ExpectedAlert, gotAlert)
		assert.Equal(t, nil, gotErr)
	}
}

func TestMakeDecision(t *testing.T) {
	var cases = []struct {
		LiveGWs          []string
		IdleGWs          []string
		AllGWs           []string
		Alert            int
		ExpectedDecision Decision
	}{
		{
			[]string{"1", "2", "3"},
			[]string{"2"},
			[]string{"1", "2", "3", "4", "5"},
			alertHighGwConn,
			Operation{Action: actionAdd, GwIP: "4", Reason: ""},
		},
	}

	for _, c := range cases {
		gotDecision := makeDecision(c.LiveGWs, c.IdleGWs, c.AllGWs, c.Alert)
		assert.Equal(t, c.ExpectedDecision, gotDecision)
	}
}
