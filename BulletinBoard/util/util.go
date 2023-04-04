package util

import (
	"board/config"
	"crypto/sha256"
	"fmt"
)

// Hash any type message v, using SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}

func GetLeader(round int, view int) int {
	f := config.GetF()
	n := 3*f + 1

	return (round + view) % n
}
