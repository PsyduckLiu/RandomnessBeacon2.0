package signature

import (
	"collector/config"
	"collector/util"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	blsSig "go.dedis.ch/dela/crypto/bls"
)

func VerifySig(msgtype int64, round string, sender string, tc []string, sig string) bool {
	tHash := new(big.Int).SetBytes(util.Digest((msgtype)))
	rHash := new(big.Int).SetBytes(util.Digest((round)))
	tcHash := new(big.Int).SetBytes(util.Digest(tc))

	e := big.NewInt(0)
	e.Xor(e, tHash)
	e.Xor(e, rHash)
	e.Xor(e, tcHash)

	id, _ := strconv.Atoi(sender)
	pk := config.GetKey(id)

	sigByte, _ := hex.DecodeString(sig)
	signatureRecover := blsSig.NewSignature(sigByte)

	fmt.Println("VerifySig", signatureRecover)
	err := pk.Verify(e.Bytes(), signatureRecover)

	return err == nil
}
