package config

var GlobCfg = Config{}

type Config struct {
	DB_NAME    string        `toml:"db_name"`
	DB_USER    string        `toml:"db_user"`
	DB_PASS    string        `toml:"db_pass"`
	DB_CHARSET string        `toml:"db_charset"`

	DOMAIN              string        `toml:"domain"`
	PROXY_PORT          string        `toml:"proxy_port"`
	MANAGE_PORT         string        `toml:"manage_port"`
	ENABLE_DEFAULT_SITE bool          `toml:"enable_default_site"`
	DEFAULT_SITE        string        `toml:"default_site"`
}

func ParseDSN(config Config) string {
	return config.DB_USER + ":" + config.DB_PASS + "@/" + config.DB_NAME + "?charset=" + config.DB_CHARSET + "&parseTime=true"
}
