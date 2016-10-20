package autoscaling

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestAnalyseAlert(t *testing.T) {
	var cases = []struct {
		Instances     int
		ConnNum       int
		HighThreshold int
		AllScaleInIPs []string
		ExpectAlert   int
	}{
		{2, 100, 100, []string{"192.168.1.1"}, alertIdleGw},
	}

	for _, c := range cases {
		gotAlert, gotErr := analyseAlert(c.Instances, c.ConnNum, c.HighThreshold, c.AllScaleInIPs)
		assert.Equal(t, c.ExpectAlert, gotAlert)
		assert.Equal(t, nil, gotErr)
	}
}
