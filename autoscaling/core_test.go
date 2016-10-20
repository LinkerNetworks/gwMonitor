package autoscaling

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestAnalyseAlert(t *testing.T) {
	var cases = []struct {
		Instances        int
		ConnNum          int
		HighThreshold    int
		AllScaleInIPsLen int
		ExpectAlert      int
	}{
		{2, 100, 100, 1, alertIdleGw},
	}

	for _, c := range cases {
		gotAlert, gotErr := analyseAlert(c.Instances, c.ConnNum, c.HighThreshold, c.AllScaleInIPsLen)
		assert.Equal(t, c.ExpectAlert, gotAlert)
		assert.Equal(t, nil, gotErr)
	}
}
