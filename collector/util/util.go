package util

import (
	"collector/config"
	"crypto/sha256"
	"fmt"
	"sort"
)

// Hash any type message v, using SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}

func RemoveRepeatElement(list []string) []string {
	if list == nil {
		return nil
	}
	out := make([]string, len(list))
	copy(out, list)
	sort.Strings(out)
	uniq := out[:0]
	for _, x := range out {
		if len(uniq) == 0 || uniq[len(uniq)-1] != x {
			uniq = append(uniq, x)
		}
	}
	return uniq
}

func GetLeader(round int, view int) int {
	f := config.GetF()
	n := 3*f + 1

	return (round + view) % n
}

func IsSubSet(x []string, y []string) bool {
	var result bool

	for _, valuex := range x {
		result = false

		for _, valuey := range y {
			if valuex == valuey {
				result = true
				break
			}
		}

		if !result {
			return result
		}
	}

	return true
}
