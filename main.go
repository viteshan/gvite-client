package main

import (
	"fmt"
	"math/big"
	"os/user"
	"path"

	"github.com/vitelabs/go-vite/ledger"

	gviteclient "github.com/vitelabs/go-vite/client"

	"github.com/vitelabs/go-vite/wallet/entropystore"

	"github.com/vitelabs/go-vite/common/types"

	"github.com/vitelabs/go-vite/wallet"
)

var RawUrl = "http://127.0.0.1:48132"

var WalletDir string

var Wallet *entropystore.Manager

func pre() {
	current, _ := user.Current()
	home := current.HomeDir
	WalletDir = path.Join(home, "Library/GVite/devdata/wallet")

	w := wallet.New(&wallet.Config{
		DataDir:        WalletDir,
		MaxSearchIndex: 100000,
	})
	w.Start()

	w2, err := w.RecoverEntropyStoreFromMnemonic("alarm canal scheme actor left length bracket slush tuna garage prepare scout school pizza invest rose fork scorpion make enact false kidney mixed vast", "en", "123456", nil)

	if err != nil {
		panic(err)
	}
	err = w2.Unlock("123456")
	if err != nil {

		panic(err)
	}

	Wallet = w2

}

func main() {
	pre()

	rpc, err := gviteclient.NewRpcClient(RawUrl)
	if err != nil {
		panic(err)
	}

	client, e := gviteclient.NewClient(rpc)
	if e != nil {
		panic(e)
	}
	self, err := types.HexToAddress("vite_ab24ef68b84e642c0ddca06beec81c9acb1977bbd7da27a87a")
	if err != nil {
		panic(err)
	}
	to, err := types.HexToAddress("vite_73c6b08e401608bca17272e7c59508f2e549c221ae7efccd53")
	if err != nil {
		panic(err)
	}
	// send tx to server
	err = client.SubmitRequestTx(gviteclient.RequestTxParams{
		ToAddr:       to,
		SelfAddr:     self,
		Amount:       big.NewInt(10000),
		TokenId:      ledger.ViteTokenId,
		SnapshotHash: nil,
		Data:         []byte("hello pow"),
	}, func(addr types.Address, data []byte) (signedData, pubkey []byte, err error) {
		return Wallet.SignData(addr, data, nil, nil)
	})
	if err != nil {
		panic(err)
	}

	// query on road tx for `to` address
	blocks, err := client.QueryOnroad(gviteclient.OnroadQuery{
		Address: to,
		Index:   1,
		Cnt:     100,
	})
	if err != nil {
		panic(err)
	}

	for _, b := range blocks {
		fmt.Println(b.Height, b.Hash, b.PrevHash, b.Amount, b.AccountAddress)
	}

	// more examples -> client_test.go
}
