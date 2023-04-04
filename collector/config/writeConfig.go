package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	blsCrypto "go.dedis.ch/dela/crypto"
)

// write public key
func WriteKey(id int, pk blsCrypto.PublicKey) {
	// set config file
	outputViper := viper.New()
	// outputViper.SetConfigFile("../Key.yml")
	outputViper.SetConfigFile("download/Key.yml")

	// read config and keep origin settings
	if err := outputViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteKey]Read config file failed:%s", err))
	}

	tag := "pk" + strconv.Itoa(id)
	pkByte, _ := pk.MarshalBinary()
	outputViper.Set(tag, string(pkByte))

	// write new settings
	if err := outputViper.WriteConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteKey]Write config file failed:%s", err))
	}

	CopyFile("download/Key.yml", "/var/www/html/Key.yml")

	fmt.Println("===>[WriteKey]Write public key success")
}
