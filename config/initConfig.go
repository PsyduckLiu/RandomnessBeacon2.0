package config

import (
	"RB/crypto/binaryquadraticform"
	"fmt"
	"math/big"
)

func Init() {
	a, b, c := GetGroupParameter()
	fmt.Println(a, b, c)
	t := GetTimeParameter()
	fmt.Println(t)

	bigOne := big.NewInt(1)
	bigTwo := big.NewInt(2)

	g, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(int64(a)), big.NewInt(int64(b)), big.NewInt(int64(c)))
	fmt.Printf("===>[InitConfig]The group element g is (a=%v,b=%v,c=%v,d=%v)\n", g.GetA(), g.GetB(), g.GetC(), g.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from InitConfig]Generate new BQuadratic Form failed: %s", err))
	}

	tPower := new(big.Int)
	tPowerSub := new(big.Int)
	tSubPower := new(big.Int)
	tPower.Exp(bigTwo, big.NewInt(int64(t)), nil)
	tPowerSub.Exp(bigTwo, big.NewInt(int64(t-1)), nil)
	tSubPower.Sub(tPower, bigOne)
	fmt.Printf("===>[InitConfig] 2^t is:%v\n", tPower)

	m_k, err := g.Exp(tPower)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from InitConfig]Generate new BQuadratic Form failed: %s", err))
	}
	m_kSub, err := g.Exp(tPowerSub)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from InitConfig]Generate new BQuadratic Form failed: %s", err))
	}
	r_k, err := g.Exp(tSubPower)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from InitConfig]Generate new BQuadratic Form failed: %s", err))
	}
	fmt.Printf("===>[InitConfig] Mk is (a=%v,b=%v,c=%v,d=%v)\n", m_k.GetA(), m_k.GetB(), m_k.GetC(), m_k.GetDiscriminant())
	fmt.Printf("===>[InitConfig] MkSub is (a=%v,b=%v,c=%v,d=%v)\n", m_kSub.GetA(), m_kSub.GetB(), m_kSub.GetC(), m_kSub.GetDiscriminant())
	fmt.Printf("===>[InitConfig] Rk is (a=%v,b=%v,c=%v,d=%v)\n", r_k.GetA(), r_k.GetB(), r_k.GetC(), r_k.GetDiscriminant())

	WriteSetup(m_k, m_kSub, r_k)
}
