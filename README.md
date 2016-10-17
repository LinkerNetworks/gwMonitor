[![Build Status](https://travis-ci.org/LinkerNetworks/gwMonitor.svg)](https://travis-ci.org/LinkerNetworks/gwMonitor)
[![Go Report](https://goreportcard.com/badge/github.com/LinkerNetworks/gwMonitor)](https://goreportcard.com/report/github.com/LinkerNetworks/gwMonitor)

# gwMonitor
Monitor and autoscaling for PGW & SGW

Clone and move this project under $GOPATH/src/github.com/LinkerNetworks/ to start your work.

# Envs

| Key        | Example           | Meaning  |
| ------------- |:-------------:| -----:|
| MONITOR_DISABLE | false | `true` to disable monitor. |
| ADDRESSES | 192.168.1.49:8000,192.168.1.49:8000 | IP addresses of OVS. |
| MONITOR_TYPE | PGW | Type of gateway, PGW or SGW. |
| PGW_CONN_NUMBER_HIGH_THRESHOLD | 200 |Threshold of PGW average connections. |
| SGW_CONN_NUMBER_HIGH_THRESHOLD | 300 | Threshold of SGW average connections. |
| MARATHON_URL | http://192.168.10.91:8080 | Endpoint of Marathon on which PGW/SGW rely. |

If `MONITOR_TYPE` is set to **PGW**, setting only `PGW_CONN_NUMBER_HIGH_THRESHOLD`.

Similarly, **PGW** for **SGW_CONN_NUMBER_HIGH_THRESHOLD**.

If env `MARATHON_URL` is set, field `marathon_url` in `monitor.conf` will be ignored.
