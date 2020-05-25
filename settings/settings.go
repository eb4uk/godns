package settings

import (
	"fmt"
	"github.com/eb4uk/godns/models"
	"github.com/spf13/pflag"
)
import "github.com/spf13/viper"

var (
	Config models.Settings
)

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
