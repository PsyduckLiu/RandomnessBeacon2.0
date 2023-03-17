package main

import (
	"RB/config"
	"RB/crypto/timedCommitment"
	"fmt"
)

func main() {
	// binaryquadraticform.TestInit()
	// binaryquadraticform.TestExp()

	config.Init()
	maskedMsg, h, M_k, a1, a2, z := timedCommitment.GenerateTC()
	fmt.Println(timedCommitment.VerifyTC(maskedMsg, h, M_k, a1, a2, z))

}
