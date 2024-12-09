package hdwallet

import (
	"crypto/ecdsa"
	"encoding/hex"

	// "github.com/bnb-chain/go-sdk/common/bech32"
	// "github.com/bnb-chain/go-sdk/common/types"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cpacia/bchutil"
	filaddr "github.com/filecoin-project/go-address"
	"github.com/myxtype/filecoin-client/local"
	filtypes "github.com/myxtype/filecoin-client/types"
	// "github.com/tendermint/tendermint/crypto/secp256k1"
)

// Key struct
type Key struct {
	Opt      *Options
	Extended *hdkeychain.ExtendedKey

	// for btc
	Private *btcec.PrivateKey
	Public  *btcec.PublicKey

	// for eth
	PrivateECDSA *ecdsa.PrivateKey
	PublicECDSA  *ecdsa.PublicKey
}

// NewKey creates a master key
// params: [Mnemonic], [Password], [Language], [Seed]
func NewKey(opts ...Option) (*Key, error) {
	var (
		err error
		o   = newOptions(opts...)
	)

	if len(o.Seed) <= 0 {
		o.Seed, err = NewSeed(o.Mnemonic, o.Password, o.Language)
	}

	if err != nil {
		return nil, err
	}

	extended, err := hdkeychain.NewMaster(o.Seed, o.Params)
	if err != nil {
		return nil, err
	}

	key := &Key{
		Opt:      o,
		Extended: extended,
	}

	err = key.init()
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (k *Key) init() error {
	var err error

	k.Private, err = k.Extended.ECPrivKey()
	if err != nil {
		return err
	}

	k.Public, err = k.Extended.ECPubKey()
	if err != nil {
		return err
	}

	k.PrivateECDSA = k.Private.ToECDSA()
	k.PublicECDSA = &k.PrivateECDSA.PublicKey
	return nil
}

// GetChildKey return a key from master key
// params: [Purpose], [CoinType], [Account], [Change], [AddressIndex], [Path]
func (k *Key) GetChildKey(opts ...Option) (*Key, error) {
	var (
		err error
		o   = newOptions(opts...)
		no  = o
	)

	typ, ok := coinTypes[o.CoinType]
	if ok {
		no = newOptions(append(opts, CoinType(typ))...)
	}

	extended := k.Extended
	for _, i := range no.GetPath() {
		extended, err = extended.Derive(i)
		if err != nil {
			return nil, err
		}
	}

	key := &Key{
		Opt:      o,
		Extended: extended,
	}

	err = key.init()
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GetWallet return wallet from master key
// params: [Purpose], [CoinType], [Account], [Change], [AddressIndex], [Path]
func (k *Key) GetWallet(opts ...Option) (Wallet, error) {
	opts = append(opts, Mnemonic(k.Opt.Mnemonic))
	key, err := k.GetChildKey(opts...)
	if err != nil {
		return nil, err
	}

	coin, ok := coins[key.Opt.CoinType]
	if !ok {
		return nil, ErrCoinTypeUnknow
	}

	return coin(key), nil
}

// PrivateHex generate private key to string by hex
func (k *Key) PrivateHex() string {
	return hex.EncodeToString(k.Private.Serialize())
}

// PrivateWIF generate private key to string by wif
func (k *Key) PrivateWIF(compress bool) (string, error) {
	wif, err := btcutil.NewWIF(k.Private, k.Opt.Params, compress)
	if err != nil {
		return "", err
	}

	return wif.String(), nil
}

// PublicHex generate public key to string by hex
func (k *Key) PublicHex(compress bool) string {
	if compress {
		return hex.EncodeToString(k.Public.SerializeCompressed())
	}

	return hex.EncodeToString(k.Public.SerializeUncompressed())
}

// PublicHash generate public key by hash160
func (k *Key) PublicHash() ([]byte, error) {
	address, err := k.Extended.Address(k.Opt.Params)
	if err != nil {
		return nil, err
	}

	return address.ScriptAddress(), nil
}

// AddressBTC generate public key to btc style address
func (k *Key) AddressBTC() (string, error) {
	address, err := k.Extended.Address(k.Opt.Params)
	if err != nil {
		return "", err
	}

	return address.EncodeAddress(), nil
}

// AddressBCH generate public key to bch style address
func (k *Key) AddressBCH() (string, error) {
	address, err := k.Extended.Address(k.Opt.Params)
	if err != nil {
		return "", err
	}

	addr, err := bchutil.NewCashAddressPubKeyHash(address.ScriptAddress(), k.Opt.Params)
	if err != nil {
		return "", err
	}

	data := addr.EncodeAddress()
	prefix := bchutil.Prefixes[k.Opt.Params.Name]
	return prefix + ":" + data, nil
}

// AddressP2WPKH generate public key to p2wpkh style address
func (k *Key) AddressP2WPKH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}

	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubHash, k.Opt.Params)
	if err != nil {
		return "", err
	}

	return addr.EncodeAddress(), nil
}

// AddressP2WPKHInP2SH generate public key to p2wpkh nested within p2sh style address
func (k *Key) AddressP2WPKHInP2SH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}

	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubHash, k.Opt.Params)
	if err != nil {
		return "", err
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", err
	}

	addr1, err := btcutil.NewAddressScriptHash(script, k.Opt.Params)
	if err != nil {
		return "", err
	}

	return addr1.EncodeAddress(), nil
}

// AddressBNB 生成bnb地址
// network mainnet主网地址 testnet测试网地址
// func (k *Key) AddressBNB(network string) (string, error) {
// 	hrp := "bnb"

// 	switch network {
// 	case MAINNET:
// 		hrp = "bnb"
// 	case TESTNET:
// 		hrp = "tbnb"
// 	}

// 	priBytes, err := hex.DecodeString(k.PrivateHex())
// 	if err != nil {
// 		return "", err
// 	}

// 	if len(priBytes) != 32 {
// 		return "", fmt.Errorf("Len of Keybytes is not equal to 32 ")
// 	}

// 	var keyBytesArray [32]byte
// 	copy(keyBytesArray[:], priBytes[:32])
// 	priKey := secp256k1..PrivKeySecp256k1(keyBytesArray)
// 	addr := types.AccAddress(priKey.PubKey().Address())

// 	bech32Addr, err := bech32.ConvertAndEncode(hrp, addr.Bytes())
// 	if err != nil {
// 		return "", err
// 	}

// 	return bech32Addr, nil
// }

func (k *Key) AddressFIL(network string) (string, error) {

	switch network {
	case MAINNET:
		filaddr.CurrentNetwork = filaddr.Mainnet
	case TESTNET:
		filaddr.CurrentNetwork = filaddr.Testnet
	default:
		filaddr.CurrentNetwork = filaddr.Mainnet
	}

	addr, err := local.WalletPrivateToAddress(local.ActSigType(filtypes.KTSecp256k1), k.Private.Serialize())
	if err != nil {
		return "", err
	}

	return addr.String(), nil
}

func (k *Key) AddressSOL() (string, error) {
	newCoin := coins[k.Opt.Change]
	return newCoin(k).GetAddress()
}
