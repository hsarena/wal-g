package main

import (
	"flag"
	"fmt"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/manifoldco/promptui"
)

const (
	API_KEY = "c952f5cfef004662bc446d6f0da7913e"
	COIN    = "btc"
)

var (
	testnet    = flag.Bool("testnet", true, "using tesnet network")
	offline    = flag.Bool("offline", false, "signed tx offline")
	chainParam = &chaincfg.TestNet3Params
	bc         gobcy.API
)

func main() {
	flag.Parse()

	bc := gobcy.API{API_KEY, COIN, "test3"}
	if !*testnet {
		chainParam = &chaincfg.MainNetParams
		bc = gobcy.API{API_KEY, COIN, "main"}
	}

	if *offline {
		fmt.Println("Using offline mode")
	}

	mainPrompt := promptui.Select{
		Label: "wal-g",
		Items: []string{"wallet",
			"transfer",
			"transaction",
			"exit"},
	}

	transferPrompt := promptui.Select{
		Label: "transfer",
		Items: []string{"get last unsinged tx",
			"update wallet balance"},
	}

	walletPrompt := promptui.Select{
		Label: "wallet",
		Items: []string{"new",
			"import",
			"info",
			"balance"},
	}

	var wal Wallet

	for {
		_, name, _ := mainPrompt.Run()
		switch name {
		case "transaction":
			pos, _, _ := transferPrompt.Run()
			switch pos {
			case 0:
				fmt.Println("txid: ", getLastUnsignedTx(bc))
			case 1:
				fmt.Println("balance: ")
			}
		case "wallet":
			pos, _, _ := walletPrompt.Run()
			switch pos {
			case 0:
				p := promptui.Prompt{
					Label: "name",
				}
				name, _ := p.Run()
				wal.createWallet(name)
			case 1:
				p := promptui.Prompt{
					Label: "name",
				}
				q := promptui.Prompt{
					Label: "wif",
				}
				name, _ := p.Run()
				wif, _ := q.Run()
				wal.improtWallet(name, wif)
			case 2:
				wal.getWalletInfo()
			case 3:
				wal.getBalance()
			}
		case "exit":
			fmt.Println("bye! =)")
			return
		}
	}

}
