package autoscaling

import marathon "github.com/gambol99/go-marathon"

// MinAppset is simplified wrapper for group .
type MinAppset struct {
	Name          string         `json:"name"`
	CreatedByJson bool           `json:"created_by_json"`
	Status        string         `json:"status"`
	Group         marathon.Group `json:"group"`
}

type AppsetResp struct {
	Success bool       `json:"success"`
	Data    *MinAppset `json:"data"`
}

type MinComponent struct {
	AppsetName string               `json:"appset_name"`
	App        marathon.Application `json:"app"`
}

type ErrResp struct {
	Success bool  `json:"success"`
	Error   Error `json:"error"`
}

type Error struct {
	Code     string `code`
	ErrorMsg string `errormsg`
}
