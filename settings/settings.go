package settings

import (
	"fmt"
	"github.com/spf13/pflag"
	"strconv"
)
import "github.com/spf13/viper"

var (
	Config Settings
)

var LogLevelMap = map[string]int{
	"DEBUG":  LevelDebug,
	"INFO":   LevelInfo,
	"NOTICE": LevelNotice,
	"WARN":   LevelWarn,
	"ERROR":  LevelError,
}

type Settings struct {
	Version      string
	Debug        bool
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

func (s RedisSettings) Addr() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

type LogSettings struct {
	Stdout bool
	File   string
	Level  string
}

func (ls LogSettings) LogLevel() int {
	l, ok := LogLevelMap[ls.Level]
	if !ok {
		panic("Config error: invalid log level: " + ls.Level)
	}
	return l
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

func InitializeConfig() {

	viper.SetConfigName("godns.conf")   // name of config file (without extension)
	viper.SetConfigType("toml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/godns/")  // path to look for the config file in
	viper.AddConfigPath("./etc/godns/") // path to look for the config file in
	viper.AddConfigPath("$HOME/.godns")
	viper.AddConfigPath("./etc/")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	pflag.Bool("v", false, "verbose output")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)
	viper.Unmarshal(&Config)

	if viper.GetBool("v") {
		Config.Log.Stdout = true
		Config.Log.Level = "DEBUG"
	}
}
