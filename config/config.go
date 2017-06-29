package config

var GlobCfg = Config{}

type Config struct {
	DB_NAME    string        `toml:"db_name"`
	DB_USER    string        `toml:"db_user"`
	DB_PASS    string        `toml:"db_pass"`
	DB_CHARSET string        `toml:"db_charset"`

	PROXY_PORT  string        `toml:"proxy_port"`
	MANAGE_PORT string        `toml:"manage_port"`
}

func ParseDSN(config Config) string {
	return config.DB_USER + ":" + config.DB_PASS + "@/" + config.DB_NAME + "?charset=" + config.DB_CHARSET
}
