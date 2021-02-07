package main

import (
	"fmt"
	"math/big"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
)

type Wallet struct {
	name          string
	balance       big.Int
	addressPubKey btcutil.AddressPubKey
	wif           btcutil.WIF
}

func (wallet *Wallet) updateBalance(bc gobcy.API) {
	addr, err := bc.GetAddrBal(wallet.addressPubKey.AddressPubKeyHash().String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	wallet.balance = addr.Balance
}

func (wallet *Wallet) Create(name string) error {

	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return err
	}

	wif, err := btcutil.NewWIF(privateKey, chainParam, true)
	if err != nil {
		return err
	}

	publicKey, err := btcutil.NewAddressPubKey(privateKey.PubKey().SerializeCompressed(), chainParam)
	if err != nil {
		return err
	}
	wallet.name = name
	wallet.balance = big.Int{}
	wallet.addressPubKey = *publicKey
	wallet.wif = *wif

	fmt.Println("name: ", wallet.name)
	fmt.Println("balance: ", wallet.balance.Int64())
	fmt.Println("public key: ", wallet.addressPubKey.AddressPubKeyHash())
	fmt.Println("wif: ", wallet.wif.String())

	return nil
}

func (wallet *Wallet) Import(name, inputWIF string) error {
	wif, err := btcutil.DecodeWIF(inputWIF)
	if err != nil {
		fmt.Println(err)
		return err
	}
	publicKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), chainParam)
	if err != nil {
		return err
	}
	wallet.name = name
	wallet.balance = big.Int{}
	wallet.addressPubKey = *publicKey
	wallet.wif = *wif

	return nil
}

func (wallet *Wallet) GetInfo() {
	fmt.Println("name: ", wallet.name)
	fmt.Println("wif: ", wallet.wif.String())
	fmt.Println("pub addr: ", wallet.addressPubKey.AddressPubKeyHash().String())
}

func (wallet *Wallet) GetBalance(bc gobcy.API) {
	wallet.updateBalance(bc)
	fmt.Println("balance:", wallet.balance.Int64())
}
