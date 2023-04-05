package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// get cboard ip from config file
func GetBoardIP() string {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../boardIP.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	return configViper.GetString("ip")
}

// get class group parameter from config file
func GetGroupParameter() (int, int, int) {
	// set config file
	configViper := viper.New()
	// // configViper.SetConfigFile("../Config.yml")
	configViper.SetConfigFile("download/Config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	return configViper.GetInt("a"), configViper.GetInt("b"), configViper.GetInt("c")
}

// get time parameter from config file
func GetTimeParameter() int {
	// set config file
	configViper := viper.New()
	// configViper.SetConfigFile("../Config.yml")
	configViper.SetConfigFile("download/Config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	return configViper.GetInt("t")
}

// get public group parameter from config file
func GetPublicGroupParameter() (int, int, int, int, int, int) {
	// set config file
	configViper := viper.New()
	// configViper.SetConfigFile("../Config.yml")
	configViper.SetConfigFile("download/Config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	return configViper.GetInt("m_k_a"), configViper.GetInt("m_k_b"), configViper.GetInt("m_k_c"), configViper.GetInt("r_k_a"), configViper.GetInt("r_k_b"), configViper.GetInt("r_k_c")
}

// get public parameter proof from config file
func GetPublicParameterProof() (int, int, int) {
	// set config file
	configViper := viper.New()
	// configViper.SetConfigFile("../Config.yml")
	configViper.SetConfigFile("download/Config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	return configViper.GetInt("p_a"), configViper.GetInt("p_b"), configViper.GetInt("p_c")
}

// get fault node number
func GetF() int {
	// set config file
	configViper := viper.New()
	// configViper.SetConfigFile("../Config.yml")
	configViper.SetConfigFile("download/Config.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetF]Read config file failed:%s", err))
	}

	return configViper.GetInt("f")
}

// get peer IP from ipAdress config file
func GetPeerIP() []string {
	// set config file
	configViper := viper.New()
	// configViper.SetConfigFile("../Config.yml")
	configViper.SetConfigFile("download/IP.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	var ipList []string
	ipList = append(ipList, configViper.GetString("peer0"))
	ipList = append(ipList, configViper.GetString("peer1"))
	ipList = append(ipList, configViper.GetString("peer2"))
	ipList = append(ipList, configViper.GetString("peer3"))

	return ipList
}
