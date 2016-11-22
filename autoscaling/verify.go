package autoscaling

import (
	"log"
	"os"
	"strings"
)

const (
	keyEth1 = "LINKER_ETH1_IP"
	keyEth2 = "LINKER_ETH2_IP"
)

// call this function after template loaded
func verifyJSON() {
	switch env(keyMonitorType).Value {
	case typePGW:
		verifyPgwJSON()
	case typeSGW:
		verifySgwJSON()
	default:
	}
}

func verifyPgwJSON() {
	for _, app := range gwGroup.Apps {
		log.Printf("verifying app with id %s...\n", app.ID)
		envMap := *app.Env
		if envMap[keyEth1] == "" {
			log.Printf("verify pgw json failed, env %s not set. ", keyEth1)
			os.Exit(1)
		}
		if envMap[keyScaleInIP] == "" {
			log.Printf("verify pgw json failed, env %s not set", keyScaleInIP)
			os.Exit(1)
		}
		var expectedScaleInIP string
		arr := strings.Split(envMap[keyEth1], "/")
		if len(arr) >= 1 {
			expectedScaleInIP = arr[0]
		}
		if envMap[keyScaleInIP] != expectedScaleInIP {
			log.Printf("%s:%s, %s:%s\n", keyEth1, envMap[keyEth1], keyScaleInIP, envMap[keyScaleInIP])
			log.Printf("verify pgw json failed, env %s != %s", keyScaleInIP, keyEth1)
			os.Exit(1)
		}
	}
}

func verifySgwJSON() {
	for _, app := range gwGroup.Apps {
		log.Printf("verifying app with id %s...\n", app.ID)
		envMap := *app.Env
		if envMap[keyEth1] == "" {
			log.Printf("verify sgw json failed, env %s not set. ", keyEth1)
			os.Exit(1)
		}
		if envMap[keyScaleInIP] == "" {
			log.Printf("verify sgw json failed, env %s not set", keyScaleInIP)
			os.Exit(1)
		}
		var expectedScaleInIP string
		arr := strings.Split(envMap[keyEth1], "/")
		if len(arr) >= 1 {
			expectedScaleInIP = arr[0]
		}
		if envMap[keyScaleInIP] != expectedScaleInIP {
			log.Printf("%s:%s, %s:%s\n", keyEth1, envMap[keyEth1], keyScaleInIP, envMap[keyScaleInIP])
			log.Printf("verify sgw json failed, env %s != %s", keyScaleInIP, keyEth1)
			os.Exit(1)
		}
	}

}

func verifyEnv() {
	log.Println("verifying env...")
	mustSet(keyMonitorType)
	mustSet(keyAddresses)
	mustSet(keyGwHighThreshold)
	mustSet(keyGwLowThreshold)
}

func mustSet(key string) {
	if strings.TrimSpace(env(key).Value) == "" {
		log.Printf("must set env %s\n", key)
		os.Exit(1)
	}
}
