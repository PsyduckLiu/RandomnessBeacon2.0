package main

import (
	"RB/config"
	"RB/crypto/binaryquadraticform"
	"RB/crypto/timedCommitment"
	"fmt"
)

func main() {
	binaryquadraticform.TestInit()
	binaryquadraticform.TestExp()

	config.Init()
	maskedMsg, h, M_kSub, M_k, a1, a2, a3, z := timedCommitment.GenerateTC()
	fmt.Println(timedCommitment.VerifyTC(maskedMsg, h, M_kSub, M_k, a1, a2, a3, z))

	// a, b, d := config.GetGroupParameter()
	// fmt.Println(a, b, d)
}
