package domain

type OracleConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Service  string
}

type KeyDbConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Db       int
}

type Config struct {
	DataSourceDb OracleConfig
	TargetDb     KeyDbConfig
}
