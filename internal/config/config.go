package config

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	UserName string `mapstructure:"user_name"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type JwtConfig struct {
	Secret string `mapstructure:"secret"` //密钥
	Expire int    `mapstructure:"expire"` //过期时间
}

type Config struct {
	Port  int         `mapstructure:"port"`
	Debug bool        `mapstructure:"debug"` // 是否调试模式
	Mysql MysqlConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	Jwt   JwtConfig   `mapstructure:"jwt"`
}
