package main

import (
	"RB/config"
	"RB/crypto/timedCommitment"
	"RB/watch"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// binaryquadraticform.TestInit()
	// binaryquadraticform.TestExp()

	config.Init()
	maskedMsg, h, M_k, a1, a2, z := timedCommitment.GenerateTC()
	fmt.Println(timedCommitment.VerifyTC(maskedMsg, h, M_k, a1, a2, z))

	timedCommitment.ForcedOpen(maskedMsg, h)

	for {
		rand.Seed(time.Now().UnixNano())
		b := make([]byte, 6)
		rand.Read(b)
		rand_str := hex.EncodeToString(b)

		watch.WriteFile("../output", rand_str)

		time.Sleep(5 * time.Second)
	}
}
