package autoscaling

import (
	"encoding/json"
	"fmt"
	"os"
)

//for debugging
func printPretty(v interface{}, mark string) (err error) {
	fmt.Printf("*********%s\n", mark)
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return
	}
	data = append(data, '\n')
	os.Stdout.Write(data)
	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
