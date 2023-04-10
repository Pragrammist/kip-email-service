package config

type Config struct {
	Smtp Smtp
}

type Smtp struct {
	Host     string
	Port     int
	User     string
	Password string
}
