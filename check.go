package hdwallet

import (
	"github.com/tyler-smith/go-bip39"
	"strings"
)

func checkWorld(mnemonic string) bool {
	for _, v := range strings.Split(mnemonic, " ") {
		if _, ok := bip39.GetWordIndex(v); !ok {
			return false
		}
	}
	return true
}
