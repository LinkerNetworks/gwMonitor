package autoscaling

import marathon "github.com/gambol99/go-marathon"

// MinAppset is simplified wrapper for group .
type MinAppset struct {
	Name          string         `json:"name"`
	CreatedByJson bool           `json:"created_by_json"`
	Status        string         `json:"status"`
	Group         marathon.Group `json:"group"`
}

type Resp struct {
	Success bool  `json:"success"`
	Error   Error `json:"error"`
}

type Error struct {
	Code     string `code`
	ErrorMsg string `errormsg`
}
