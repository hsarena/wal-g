package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/manifoldco/promptui"
)

const (
	API_KEY = "a53948ef98be4aa7af84c0d96922ff0d"
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
			"transaction",
			"exit"},
	}

	transferPrompt := promptui.Select{
		Label: "transaction",
		Items: []string{"get last unsinged tx",
			"update wallet balance",
			"create transaction",
			"sign transaction",
			"push transaction",
			"track transaction"},
	}

	walletPrompt := promptui.Select{
		Label: "wallet",
		Items: []string{"new",
			"import",
			"info",
			"balance"},
	}

	var wal Wallet
	var tx Transaction

	for {
		_, name, _ := mainPrompt.Run()
		switch name {
		case "transaction":
			pos, _, _ := transferPrompt.Run()
			switch pos {
			case 0:
				fmt.Println("txid: ", GetLastUnsignedTx(bc))
			case 1:
				wal.GetBalance(bc)
			case 2:
				p := promptui.Prompt{
					Label: "amount",
				}
				q := promptui.Prompt{
					Label: "destination address",
				}
				amount, _ := p.Run()
				destAddr, _ := q.Run()
				
				i_amount, _:= strconv.Atoi(amount)

				tx.CreateTX(bc, &wal, destAddr, int64(i_amount) )
			case 3:
				fmt.Println(tx.SignTx(&wal))
			case 4:
				tx.PushTX(bc)
			case 5:
				p := promptui.Prompt{
					Label: "txid",
				}
				txid, _ := p.Run()
				GetTxInfo(bc, txid)

			}
		case "wallet":
			pos, _, _ := walletPrompt.Run()
			switch pos {
			case 0:
				p := promptui.Prompt{
					Label: "name",
				}
				name, _ := p.Run()
				wal.Create(name)
			case 1:
				p := promptui.Prompt{
					Label: "name",
				}
				q := promptui.Prompt{
					Label: "wif",
				}
				name, _ := p.Run()
				wif, _ := q.Run()
				wal.Import(name, wif)
			case 2:
				wal.GetInfo()
			case 3:
				wal.GetBalance(bc)
			}
		case "exit":
			fmt.Println("bye! =)")
			return
		}
	}

}
