package hdwallet

func init() {
	coins[USDT] = newUSDT
}

type usdt struct {
	*eth
}

func newUSDT(key *Key) Wallet {
	token := newETH(key).(*eth)
	token.name = "Tether"
	token.symbol = "USDT"

	return &usdt{eth: token}
}
