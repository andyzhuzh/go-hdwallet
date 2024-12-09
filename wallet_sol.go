package hdwallet

import (
	"crypto/ed25519"
	"fmt"

	"github.com/blocto/solana-go-sdk/pkg/hdwallet"
	"github.com/btcsuite/btcutil/base58"
	"github.com/tyler-smith/go-bip39"
)

func init() {
	coins[SOL] = newSOL
}

type sol struct {
	name   string
	symbol string
	priKey ed25519.PrivateKey
	pubKey ed25519.PrivateKey
	key    *Key
}

func newSOL(key *Key) Wallet {
	return &sol{
		name:   "Solane",
		symbol: "SOL",
		key:    key,
	}
}

func (c *sol) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *sol) GetName() string {
	return c.name
}

func (c *sol) GetSymbol() string {
	return c.symbol
}

func (c *sol) GetKey() *Key {
	return c.key
}

func (c *sol) GetAddress() (string, error) {
	if c.priKey == nil {
		if _, err := c.GetPrivateKey(); err != nil {
			return "", err
		}
	}
	return base58.Encode(c.pubKey), nil
}

func (c *sol) GetPrivateKey() (string, error) {
	seed, err := bip39.NewSeedWithErrorChecking(c.key.Opt.Mnemonic, "")
	if err != nil {
		return "", err
	}
	k, err := hdwallet.Derived(fmt.Sprintf("m/44'/%d'/%d'/%d'", c.key.Opt.CoinType-ZeroQuote, 0, c.key.Opt.AddressIndex), seed)
	if err != nil {
		return "", err
	}
	c.priKey = ed25519.NewKeyFromSeed(k.PrivateKey)
	c.pubKey = make([]byte, ed25519.PublicKeySize)
	copy(c.pubKey, c.priKey[32:])
	return base58.Encode(c.priKey), nil
}
