package config

type DbConfig struct {
	User     string `default:"postgres"`
	Password string `default:"1234"`
	Host     string `default:"localhost"`
	Path     string `default:"postgres"`
	SslMode  bool   `default:"false"`
}
