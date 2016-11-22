package conf

var defaultOptions = Options{
	PollingSeconds:      1,
	GwOverloadTolerance: 120,
	GwIdleTolerance:     120,
	ClientEndpoint:      "master.mesos:10004",
	PgwJSON:             "pgw.json",
	SgwJSON:             "sgw.json",
}
