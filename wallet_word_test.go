package hdwallet

import (
	"fmt"
	"github.com/block-vision/sui-go-sdk/signer"
	//"golang.org/x/crypto/blake2b"
	"testing"
)

func TestSui(t *testing.T) {
	wd := "range sheriff try enroll deer over ten level bring display stamp recycle"
	signerAccount, err := signer.NewSignertWithMnemonic(wd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("addr", signerAccount.Address)

	master, err := NewKey(
		Mnemonic(wd),
	)
	if err != nil {
		panic(err)
	}
	wallet, _ := master.GetWallet(CoinType(SUI))
	//wallet, _ := master.GetWallet(Purpose(hdwallet.ZeroQuote+44), hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(0))
	address, _ := wallet.GetAddress() // 1AwEPfoojHnKrhgt1vfuZAhrvPrmz7Rh44
	fmt.Println("addr", address)
}
