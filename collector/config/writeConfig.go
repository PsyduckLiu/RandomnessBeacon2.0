package config

import (
	"collector/crypto/binaryquadraticform"
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	blsCrypto "go.dedis.ch/dela/crypto"
)

// write new m_k and r_k
func WriteSetup(m_k *binaryquadraticform.BQuadraticForm, r_k *binaryquadraticform.BQuadraticForm, proof *binaryquadraticform.BQuadraticForm) {
	// set config file
	outputViper := viper.New()
	outputViper.SetConfigFile("../Config.yml")

	// read config and keep origin settings
	if err := outputViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteSetup]Read config file failed:%s", err))
	}

	outputViper.Set("m_k_a", m_k.GetA())
	outputViper.Set("m_k_b", m_k.GetB())
	outputViper.Set("m_k_c", m_k.GetC())
	outputViper.Set("r_k_a", r_k.GetA())
	outputViper.Set("r_k_b", r_k.GetB())
	outputViper.Set("r_k_c", r_k.GetC())
	outputViper.Set("p_a", proof.GetA())
	outputViper.Set("p_b", proof.GetB())
	outputViper.Set("p_c", proof.GetC())

	// write new settings
	if err := outputViper.WriteConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteSetup]Write config file failed:%s", err))
	}
	// outputViper.Debug()

	fmt.Println("===>[WriteSetup]Write output success")
}

// write public key
func WriteKey(id int, pk blsCrypto.PublicKey) {
	// set config file
	outputViper := viper.New()
	outputViper.SetConfigFile("../Key.yml")

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

	fmt.Println("===>[WriteKey]Write public key success")
}
