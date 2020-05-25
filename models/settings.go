package models

import "strconv"

var LogLevelMap = map[string]int{
	"DEBUG":  LevelDebug,
	"INFO":   LevelInfo,
	"NOTICE": LevelNotice,
	"WARN":   LevelWarn,
	"ERROR":  LevelError,
}

type Settings struct {
	Version string

	Debug                  bool
	TargetResponse         bool `mapstructure:"target-response"`
	TargetResponseRedisKey bool `mapstructure:"target-response-redis-key"`

	Server       DNSServerSettings `mapstructure:"server"`
	ResolvConfig ResolvSettings    `mapstructure:"resolv"`
	Redis        RedisSettings     `mapstructure:"redis"`
	Memcache     MemcacheSettings  `mapstructure:"memcache"`
	Log          LogSettings       `mapstructure:"log"`
	Cache        CacheSettings     `mapstructure:"cache"`
	Hosts        HostsSettings     `mapstructure:"hosts"`
}

type ResolvSettings struct {
	Timeout        int
	Interval       int
	SetEDNS0       bool
	ServerListFile string `mapstructure:"server-list-file"`
	ResolvFile     string `mapstructure:"resolv-file"`
}

type DNSServerSettings struct {
	Host string
	Port int
}

type RedisSettings struct {
	Host     string
	Port     int
	DB       int
	Password string
}

type MemcacheSettings struct {
	Servers []string
}

type LogSettings struct {
	Stdout bool
	File   string
	Level  string
}

type CacheSettings struct {
	Backend  string
	Expire   int
	Maxcount int
}

type HostsSettings struct {
	Enable          bool
	HostsFile       string `mapstructure:"host-file"`
	RedisEnable     bool   `mapstructure:"redis-enable"`
	RedisKey        string `mapstructure:"redis-key"`
	TTL             uint32 `mapstructure:"ttl"`
	RefreshInterval uint32 `mapstructure:"refresh-interval"`
}

func (s RedisSettings) Addr() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

func (ls LogSettings) LogLevel() int {
	l, ok := LogLevelMap[ls.Level]
	if !ok {
		panic("Config error: invalid log level: " + ls.Level)
	}
	return l
}
