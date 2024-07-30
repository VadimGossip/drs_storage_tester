package domain

type TargetDbConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Db       int
}

type Config struct {
	TargetDb TargetDbConfig
}
