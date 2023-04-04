package main

import (
	"RB/config"
	"RB/watch"
	"encoding/hex"
	"math/rand"
	"time"
)

func main() {
	// binaryquadraticform.TestInit()
	// binaryquadraticform.TestExp()

	config.InitGroup()
	// maskedMsg, h, M_k, a1, a2, z := timedCommitment.GenerateTC()
	// fmt.Println(timedCommitment.VerifyTC(maskedMsg, h, M_k, a1, a2, z))

	// timedCommitment.ForcedOpen(maskedMsg, h)

	// fmt.Println("Signature Verification ", testSig())

	watch.WatchOutput("output")
	for {
		rand.Seed(time.Now().UnixNano())
		b := make([]byte, 6)
		rand.Read(b)
		rand_str := hex.EncodeToString(b)

		config.WriteFile("/var/www/html/output", rand_str)

		time.Sleep(5 * time.Second)
	}

	// config.InitGroup()

	// rand.Seed(time.Now().UnixNano())
	// b := make([]byte, 6)
	// rand.Read(b)
	// rand_str := hex.EncodeToString(b)

	// watch.WriteFile("../output", rand_str)
}

// func testSig() bool {
// 	signerA := blsSig.NewSigner()
// 	signerB := blsSig.NewSigner()
// 	signerC := blsSig.NewSigner()

// 	publicKeys := []blsCrypto.PublicKey{
// 		signerA.GetPublicKey(),
// 		signerB.GetPublicKey(),
// 	}

// 	message := []byte("42")

// 	signatureA, err := signerA.Sign(message)
// 	if err != nil {
// 		panic("signer A failed: " + err.Error())
// 	}

// 	signatureB, err := signerB.Sign(message)
// 	if err != nil {
// 		panic("signer B failed: " + err.Error())
// 	}

// 	aggregate, err := signerC.Aggregate(signatureA, signatureB)
// 	if err != nil {
// 		panic("aggregate failed: " + err.Error())
// 	}

// 	verifier, err := signerC.GetVerifierFactory().FromArray(publicKeys)
// 	if err != nil {
// 		panic("verifier failed: " + err.Error())
// 	}

// 	err = verifier.Verify(message, aggregate)
// 	if err != nil {
// 		return false
// 	}

// 	fmt.Println("Success")
// 	return true
// }
