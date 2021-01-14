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

func (wallet *Wallet) updateWalletBalance(bc gobcy.API) {
	addr, _ := bc.GetAddrBal(wallet.addressPubKey.String(), nil)
	wallet.balance = addr.Balance
}

func (wallet *Wallet) createWallet(name string) error {

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

	return nil
}

func (wallet *Wallet) improtWallet(name, inputWIF string) error {
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

func (wallet *Wallet) getWalletInfo() {
	fmt.Println("name: ", wallet.name)
	fmt.Println("wif: ", wallet.wif.String())
	fmt.Println("pub addr: ", wallet.addressPubKey.AddressPubKeyHash().String())
}

func (wallet *Wallet) getBalance() {
	wallet.updateWalletBalance(bc)
	println("blance:", wallet.balance.Int64())
}
