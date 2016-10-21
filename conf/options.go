package conf

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/hcl"
)

const configFile = "monitor.conf"

var (
	once sync.Once

	// OptionsReady is a loaded object containing fields in config file.
	// Using defaults if not set.
	OptionsReady = &Options{}
)

// Options is the structure for config file.
type Options struct {
	PollingTime         int    `hcl:"polling_time"`
	GwOverloadTolerance int    `hcl:"gw_overload_tolerance"`
	GwIdleTolerance     int    `hcl:"gw_idle_tolerance"`
	ClientEndpoint      string `hcl:"client_endpoint"`
	PgwJSON             string `hcl:"pgw_json"`
	SgwJSON             string `hcl:"sgw_json"`
}

func init() {
	once.Do(func() {
		options := newOptions()
		if err := options.load(); err != nil {
			log.Printf("load options error: %v, using default", err)
			OptionsReady = &defaultOptions
			return
		}
		OptionsReady = options
	})
}

// NewOptions returns new config
func newOptions() *Options {
	return &Options{}
}

// Load loads options from config file
func (p *Options) load() error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("stat config file error: %v", err)
		return err
	}

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("read config file error: %v", err)
		return err
	}

	if err := hcl.Decode(p, string(content)); err != nil {
		log.Printf("decode hcl error: %v", err)
		return err
	}
	return nil
}
