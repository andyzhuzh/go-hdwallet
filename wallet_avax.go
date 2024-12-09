package hdwallet

import "github.com/ethereum/go-ethereum/crypto"

func init() {
	coins[AVAX] = newAVAX
}

type avax struct {
	name   string
	symbol string
	key    *Key

	// bnb token
	contract string
}

func newAVAX(key *Key) Wallet {
	return &avax{
		name:   "Avalanche",
		symbol: "AVAX",
		key:    key,
	}
}

func (c *avax) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *avax) GetName() string {
	return c.name
}

func (c *avax) GetSymbol() string {
	return c.symbol
}

func (c *avax) GetKey() *Key {
	return c.key
}

func (c *avax) GetAddress() (string, error) {
	return crypto.PubkeyToAddress(*c.key.PublicECDSA).Hex(), nil
}

func (c *avax) GetPrivateKey() (string, error) {
	return c.key.PrivateHex(), nil
}
