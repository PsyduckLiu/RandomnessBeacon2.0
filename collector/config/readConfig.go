package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	blsCrypto "go.dedis.ch/dela/crypto"
	blsSig "go.dedis.ch/dela/crypto/bls"
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
	// configViper.SetConfigFile("../Config.yml")
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

// get public key
func GetKey(id int) blsCrypto.PublicKey {
	// set config file
	configViper := viper.New()
	// configViper.SetConfigFile("../Key.yml")
	configViper.SetConfigFile("download/Key.yml")

	// read config and keep origin settings
	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetKey]Read config file failed:%s", err))
	}

	tag := "pk" + strconv.Itoa(id)
	pkByte := configViper.GetString(tag)

	pk, err := blsSig.NewPublicKey([]byte(pkByte))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GetKey]Recover key failed:%s", err))
	}

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	// fmt.Println("===>[GetKey]Get public key success", pk)
	return pk
}

// get total TC from TC file
func GetTC() []string {
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

// get board IP from ipAdress config file
func GetBoardLisIP() string {
	// set config file
	configViper := viper.New()
	// configViper.SetConfigFile("../Config.yml")
	configViper.SetConfigFile("download/IP.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetGroupParameter]Read config file failed:%s", err))
	}

	return configViper.GetString("board")
}
