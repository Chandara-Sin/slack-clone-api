package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config-dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Fatal error config file: %s \n", err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
