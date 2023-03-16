package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// get class group parameter from config file
func GetGroupParameter() (int, int, int) {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("Config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	return configViper.GetInt("a"), configViper.GetInt("b"), configViper.GetInt("c")
}

// get time parameter from config file
func GetTimeParameter() int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("Config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	return configViper.GetInt("t")
}

// get public group parameter from config file
func GetPublicGroupParameter() (int, int, int, int, int, int, int, int, int) {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("Config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	return configViper.GetInt("m_k_a"), configViper.GetInt("m_k_b"), configViper.GetInt("m_k_c"), configViper.GetInt("r_k_a"), configViper.GetInt("r_k_b"), configViper.GetInt("r_k_c"), configViper.GetInt("m_kSub_a"), configViper.GetInt("m_kSub_b"), configViper.GetInt("m_kSub_c")
}
