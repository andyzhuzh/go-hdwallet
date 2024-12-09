package hdwallet

// mnemonic language
const (
	English            = "english"
	ChineseSimplified  = "chinese_simplified"
	ChineseTraditional = "chinese_traditional"
)

// zero is deafult of uint32
const (
	Zero      uint32 = 0
	ZeroQuote uint32 = 0x80000000
	BTCToken  uint32 = 0x10000000
	ETHToken  uint32 = 0x20000000
)

// wallet type from bip44
const (
	// https://github.com/satoshilabs/slips/blob/master/slip-0044.md#registered-coin-types
	BTC        = ZeroQuote + 0
	BTCTestnet = ZeroQuote + 1
	LTC        = ZeroQuote + 2
	DOGE       = ZeroQuote + 3
	DASH       = ZeroQuote + 5
	ETH        = ZeroQuote + 60
	ETC        = ZeroQuote + 61
	BCH        = ZeroQuote + 145
	SOL        = ZeroQuote + 501
	QTUM       = ZeroQuote + 2301
	TRX        = ZeroQuote + 195
	AVAX       = ZeroQuote + 713
	BNB        = ZeroQuote + 714
	MATIC      = ZeroQuote + 9999
	FIL        = ZeroQuote + 461
	SUI        = ZeroQuote + 784

	// btc token
	USDT = BTCToken + 1

	// eth token
	IOST = ETHToken + 1
	USDC = ETHToken + 2
)

// network
const (
	MAINNET = "mainnet"
	TESTNET = "testnet"
)

var coinTypes = map[uint32]uint32{
	USDT: ETH,
	IOST: ETH,
	USDC: ETH,
	TRX:  TRX,
	//BNB:  BNB,
	BNB:   ETH, // BNB与ETH地址相同
	AVAX:  ETH, // AVAX与ETH地址相同
	MATIC: ETH, // MATIC与ETH地址相同
	FIL:   FIL,
	SUI:   SUI,
	SOL:   SOL,
}
