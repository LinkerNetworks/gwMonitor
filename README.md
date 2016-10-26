[![Build Status](https://travis-ci.org/LinkerNetworks/gwMonitor.svg)](https://travis-ci.org/LinkerNetworks/gwMonitor)
[![Go Report](https://goreportcard.com/badge/github.com/LinkerNetworks/gwMonitor)](https://goreportcard.com/report/github.com/LinkerNetworks/gwMonitor)

# gwMonitor
Monitor and autoscaling for PGW & SGW

Clone and move this project under $GOPATH/src/github.com/LinkerNetworks/ to start your work.

# Envs

| Key        | Example           | Meaning  |Default
| :--------- |:----------------:| :---------|:--------:
| ADDRESSES | 192.168.1.49:8000,192.168.1.50:8000 | IP addresses of OVS. | "" |
| MONITOR_TYPE | PGW | Type of gateway, PGW or SGW. | "" |
| GW_CONN_NUMBER_HIGH_THRESHOLD | 200 |High threshold of GW average connections. | 0 |
| GW_CONN_NUMBER_LOW_THRESHOLD | 100 |Low threshold of GW average connections. | 0 |
| CLIENT_ENDPOINT | 192.168.10.91:10004 | Endpoint of Linker DC/OS client. | "" |

`GW_CONN_NUMBER_HIGH_THRESHOLD` and `GW_CONN_NUMBER_LOW_THRESHOLD` apply to PGW/SGW both.

If env `CLIENT_ENDPOINT` is set, field `client_endpoint` in `monitor.conf` will be ignored.

# Docker

```sh
./build.sh
docker build -t linkerrepository/gwmonitor:dev .
docker run -e MONITOR_TYPE="PGW" -e CLIENT_ENDPOINT="192.168.10.91:10004" -e GW_CONN_NUMBER_HIGH_THRESHOLD=200 -e GW_CONN_NUMBER_LOW_THRESHOLD=100 -e ADDRESSES="192.168.10.186:18080" --network=host linkerrepository/gwmonitor:dev
```
