package model

type DB struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"BD_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Database string `mapstructure:"DB_DATABASE"`
	Driver   string `mapstructure:"DB_DRIVER"`
}
