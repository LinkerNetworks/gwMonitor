package conf

var defaultOptions = Options{
	PollingTime:         1,
	GwOverloadTolerance: 60,
	GwIdleTolerance:     300,
	ClientEndpoint:      "master.mesos:10004",
	PgwJSON:             "pgw.json",
	SgwJSON:             "sgw.json",
}
