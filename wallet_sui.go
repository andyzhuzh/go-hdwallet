package hdwallet

import (
	"encoding/hex"
	"github.com/block-vision/sui-go-sdk/signer"
)

func init() {
	coins[SUI] = newSUI
}

type sui struct {
	name          string
	symbol        string
	key           *Key
	signerAccount *signer.Signer
}

func newSUI(key *Key) Wallet {
	return &sui{
		name:   "SUI",
		symbol: "SUI",
		key:    key,
	}
}

func (c *sui) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *sui) GetName() string {
	return c.name
}

func (c *sui) GetSymbol() string {
	return c.symbol
}

func (c *sui) GetKey() *Key {
	return c.key
}

func (c *sui) GetAddress() (string, error) {
	if c.signerAccount == nil {
		if _, err := c.GetPrivateKey(); err != nil {
			return "", err
		}
	}
	return c.signerAccount.Address, nil
}

func (c *sui) GetPrivateKey() (string, error) {
	var err error
	if c.signerAccount == nil {
		c.signerAccount, err = signer.NewSignertWithMnemonic(c.key.Opt.Mnemonic)
		if err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(c.signerAccount.PriKey), nil
}
