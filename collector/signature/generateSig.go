package signature

import (
	"collector/util"
	"encoding/hex"
	"fmt"
	"math/big"

	blsSig "go.dedis.ch/dela/crypto/bls"
)

func GenerateSig(msgtype int64, round string, sender string, tc []string, signer blsSig.Signer) string {
	tHash := new(big.Int).SetBytes(util.Digest((msgtype)))
	rHash := new(big.Int).SetBytes(util.Digest((round)))
	sHash := new(big.Int).SetBytes(util.Digest((sender)))
	tcHash := new(big.Int).SetBytes(util.Digest(tc))

	e := big.NewInt(0)
	e.Xor(e, tHash)
	e.Xor(e, rHash)
	e.Xor(e, sHash)
	e.Xor(e, tcHash)

	signature, err := signer.Sign(e.Bytes())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateSig]Failed to generate signature: %s", err))
	}

	result, err := signature.MarshalBinary()
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateSig]Failed to generate signature: %s", err))
	}

	return hex.EncodeToString(result)
}
