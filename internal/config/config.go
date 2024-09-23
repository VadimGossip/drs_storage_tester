package config

import "time"

type OracleConfig interface {
	DSN() string
}

type KdbConfig interface {
	Address() string
	Username() string
	Password() string
	DB() int
	ReadTimeoutSec() time.Duration
	WriteTimeoutSec() time.Duration
}

type TarantoolConfig interface {
	Address() string
	Username() string
	Password() string
}

type RateGrpcConfig interface {
	Address() string
}

type ImitatorConfig interface {
	RequestType() string
	AllSupplierRequestType() string
	SingleRequestType() string
	RequestPerSecond() int
	PackPerSecond() int
	TotalRequests() int
}

type ServiceProviderConfig interface {
	TestDB() string
	TarantoolTestDB() string
	KdbTestDB() string
	CacheTestDB() string
}
