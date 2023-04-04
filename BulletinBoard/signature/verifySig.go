package signature

import (
	"board/config"
	"board/util"
	"encoding/hex"
	"math/big"
	"strconv"

	blsSig "go.dedis.ch/dela/crypto/bls"
)

func VerifySig(msgtype int64, round string, view string, sender string, tc []string, sig string, from string) bool {
	tHash := new(big.Int).SetBytes(util.Digest((msgtype)))
	rHash := new(big.Int).SetBytes(util.Digest((round)))
	vHash := new(big.Int).SetBytes(util.Digest((view)))
	sHash := new(big.Int).SetBytes(util.Digest((sender)))
	tcHash := new(big.Int).SetBytes(util.Digest(tc))

	e := big.NewInt(0)
	e.Xor(e, tHash)
	e.Xor(e, rHash)
	e.Xor(e, vHash)
	e.Xor(e, sHash)
	e.Xor(e, tcHash)

	id, _ := strconv.Atoi(from)
	pk := config.GetKey(id)

	sigByte, _ := hex.DecodeString(sig)
	signatureRecover := blsSig.NewSignature(sigByte)
	err := pk.Verify(e.Bytes(), signatureRecover)

	return err == nil
}

func VerifyNewLeaderSig(msgtype int64, round string, view string, sender string, sig string, from string) bool {
	tHash := new(big.Int).SetBytes(util.Digest((msgtype)))
	rHash := new(big.Int).SetBytes(util.Digest((round)))
	vHash := new(big.Int).SetBytes(util.Digest((view)))
	sHash := new(big.Int).SetBytes(util.Digest((sender)))

	e := big.NewInt(0)
	e.Xor(e, tHash)
	e.Xor(e, rHash)
	e.Xor(e, vHash)
	e.Xor(e, sHash)

	id, _ := strconv.Atoi(from)
	pk := config.GetKey(id)

	sigByte, _ := hex.DecodeString(sig)
	signatureRecover := blsSig.NewSignature(sigByte)
	err := pk.Verify(e.Bytes(), signatureRecover)

	return err == nil
}

func VerifyOutputSig(msgtype int64, round string, randomNumber string, sig string, from string) bool {
	tHash := new(big.Int).SetBytes(util.Digest((msgtype)))
	rHash := new(big.Int).SetBytes(util.Digest((round)))
	rnHash := new(big.Int).SetBytes(util.Digest((randomNumber)))

	e := big.NewInt(0)
	e.Xor(e, tHash)
	e.Xor(e, rHash)
	e.Xor(e, rnHash)

	id, _ := strconv.Atoi(from)
	pk := config.GetKey(id)

	sigByte, _ := hex.DecodeString(sig)
	signatureRecover := blsSig.NewSignature(sigByte)
	err := pk.Verify(e.Bytes(), signatureRecover)

	return err == nil
}
