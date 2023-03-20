package signature

import (
	"collector/config"
	"collector/util"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	blsCrypto "go.dedis.ch/dela/crypto"
	blsSig "go.dedis.ch/dela/crypto/bls"
)

func AggSig(sig []string) string {
	aggSigner := blsSig.NewSigner()
	var err error
	var singleSig []blsCrypto.Signature

	for _, s := range sig {
		sigByte, _ := hex.DecodeString(s)
		signatureRecover := blsSig.NewSignature(sigByte)
		singleSig = append(singleSig, signatureRecover)
	}

	aggSig := singleSig[0]
	for i := 1; i < len(singleSig); i++ {
		aggSig, err = aggSigner.Aggregate(aggSig, singleSig[i])
		if err != nil {
			panic("aggregate failed: " + err.Error())
		}
	}

	result, err := aggSig.MarshalBinary()
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from AggSig]Failed to generate signature: %s", err))
	}

	return hex.EncodeToString(result)
}

func VerifyAggSig(msgtype int64, round string, sender string, tc []string, aggsig string, from []string) bool {
	aggSigner := blsSig.NewSigner()
	var publicKeys []blsCrypto.PublicKey

	for _, f := range from {
		id, _ := strconv.Atoi(f)
		pk := config.GetKey(id)
		publicKeys = append(publicKeys, pk)
	}

	sigByte, _ := hex.DecodeString(aggsig)
	signatureRecover := blsSig.NewSignature(sigByte)

	verifier, err := aggSigner.GetVerifierFactory().FromArray(publicKeys)
	if err != nil {
		panic("verifier failed: " + err.Error())
	}

	tHash := new(big.Int).SetBytes(util.Digest((msgtype)))
	rHash := new(big.Int).SetBytes(util.Digest((round)))
	sHash := new(big.Int).SetBytes(util.Digest((sender)))
	tcHash := new(big.Int).SetBytes(util.Digest(tc))

	e := big.NewInt(0)
	e.Xor(e, tHash)
	e.Xor(e, rHash)
	e.Xor(e, sHash)
	e.Xor(e, tcHash)

	err = verifier.Verify(e.Bytes(), signatureRecover)
	fmt.Println(err == nil)

	return err == nil
}
