[TOC]

# IMPORTANT
This repo is no longer maintained.

It has been moved to BitBucket, redirect to [bitbucket.org/linkernetworks/gwmonitor][1]

# gwMonitor
Monitor and autoscaling for PGW & SGW

# Clone
Clone and move this project under GOPATH.

```sh
git clone git@bitbucket.org:linkernetworks/gwmonitor.git $GOPATH/src/bitbucket.org/linkernetworks/
```

# Build

```sh
./build.sh
```

# Config
## Envs

| Key        | Example           | Meaning  |Default
| :--------- |:----------------:| :---------|:--------:
| ADDRESSES | 192.168.1.49:8000,192.168.1.50:8000 | IP addresses of OVS. | "" |
| MONITOR_TYPE | PGW | Type of gateway, PGW or SGW. | "" |
| GW_CONN_NUMBER_HIGH_THRESHOLD | 200 |High threshold of GW average connections. | 0 |
| GW_CONN_NUMBER_LOW_THRESHOLD | 100 |Low threshold of GW average connections. | 0 |
| CLIENT_ENDPOINT | 192.168.10.91:10004 | Endpoint of Linker DC/OS client. | "master.mesos:10004" |
| POLLING_SECONDS | 1 | Peroid fetching from OVS. | 1 |
| GW_OVERLOAD_TOLERANCE | 60 | Max overload alert enjured times. | 120 |
| GW_IDLE_TOLERANCE | 60 | Max idle alert enjured times. | 120 |


`GW_CONN_NUMBER_HIGH_THRESHOLD` and `GW_CONN_NUMBER_LOW_THRESHOLD` apply to PGW/SGW both.

If env `CLIENT_ENDPOINT` is set, field `client_endpoint` in `monitor.conf` will be ignored.

It is recommended to set `GW_OVERLOAD_TOLERANCE` and `GW_IDLE_TOLERANCE` above 30.

# Docker

## Build

```sh
docker build -t linkerrepository/gwmonitor:dev .
```

## Run

```sh
docker run -e MONITOR_TYPE="PGW" \
	-e CLIENT_ENDPOINT="192.168.10.91:10004" \
	-e GW_CONN_NUMBER_HIGH_THRESHOLD=200 \
	-e GW_CONN_NUMBER_LOW_THRESHOLD=100 \
	-e ADDRESSES="192.168.10.186:18080" \
	-e POLLING_SECONDS="1" \
	-e GW_OVERLOAD_TOLERANCE="60" \
	-e GW_IDLE_TOLERANCE="60" \
	--network=host \
	linkerrepository/gwmonitor:dev
```

[1]: https://bitbucket.org/linkernetworks/gwmonitor/overview
