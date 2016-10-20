// linker dcos client

package autoscaling

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// call client to add an app, not started after creation
func addComponent(component *MinComponent) (err error) {
	endpoint := env(keyClientEndpoint).Value
	if strings.TrimSpace(endpoint) == "" {
		log.Printf("client endpoint not set, check env %s\n", keyClientEndpoint)
		return errors.New("invalid client endpoint")
	}

	url := "http://" + endpoint + "/v1/components"

	data, err := json.Marshal(component)
	if err != nil {
		log.Printf("json marshal component error: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, ioutil.NopCloser(bytes.NewReader(data)))
	if err != nil {
		log.Printf("new request error: %v\n", err)
		return
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Printf("http post component %s error: %v\n", url, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		errResp := &ErrResp{}
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, errResp)
		log.Println("call client to add component error: ")
		printPretty(errResp, "resp")
		return errors.New(errResp.Error.ErrorMsg)
	}
	return
}

// appId is absolute marathon app ID
func startComponent(appID string) (err error) {
	endpoint := env(keyClientEndpoint).Value
	if strings.TrimSpace(endpoint) == "" {
		log.Printf("client endpoint not set, check env %s\n", keyClientEndpoint)
		return errors.New("invalid client endpoint")
	}

	url := "http://" + endpoint + "/v1/components/start"

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		log.Printf("new request error: %v\n", err)
		return
	}

	q := req.URL.Query()
	q.Add("name", appID)
	req.URL.RawQuery = q.Encode()

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Printf("http put component %s error: %v\n", url, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		errResp := &ErrResp{}
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, errResp)
		log.Println("call client to start component error: ")
		printPretty(errResp, "resp")
		return errors.New(errResp.Error.ErrorMsg)
	}
	return
}

// call client to delete an app
func delComponent(appID string) (err error) {
	endpoint := env(keyClientEndpoint).Value
	if strings.TrimSpace(endpoint) == "" {
		log.Printf("client endpoint not set, check env %s\n", keyClientEndpoint)
		return errors.New("invalid client endpoint")
	}

	url := "http://" + endpoint + "/v1/components"

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Printf("new request error: %v\n", err)
		return
	}

	q := req.URL.Query()
	q.Add("name", appID)
	req.URL.RawQuery = q.Encode()

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Printf("http delete component %s error: %v\n", url, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		errResp := &ErrResp{}
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, errResp)
		log.Println("call client to delete component error: ")
		printPretty(errResp, "resp")
		return errors.New(errResp.Error.ErrorMsg)
	}
	return
}
