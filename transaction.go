package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Transaction struct {
	TxId               string 		`json:"txid"`
	SourceAddress      string 		`json:"source_address"`
	DestinationAddress string 		`json:"destination_address"`
	Amount             int64  		`json:"amount"`
	UnsignedTx         string 		`json:"unsignedtx"`
	SignedTx           string 		`json:"signedtx"`
	RedeemTx 		   *wire.MsgTx  `json:"redeemtx"`
}

func GetLastUnsignedTx(bc gobcy.API) string {
	txid, _ := bc.GetUnTX()
	return txid[0].Hash
}

func GetTxInfo(bc gobcy.API, txid string) {
	fmt.Println(bc.GetTX(txid,nil))
	fmt.Printf("https://www.blockchain.com/btc-testnet/tx/%#v", txid)
}

func (transaction *Transaction) CreateTX(bc gobcy.API, wallet *Wallet, destination string, amount int64) error {
	wallet.updateBalance(bc)
	srcAddr := wallet.addressPubKey.AddressPubKeyHash()
	destAddr, err := btcutil.DecodeAddress(destination, chainParam)

	balance := wallet.balance.Int64()
	fmt.Println(balance)
	if amount > balance {
		fmt.Println("insufficent balance ", amount - balance)
	}

	txHash := GetLastUnsignedTx(bc)
	srcTx := wire.NewMsgTx(wire.TxVersion)
	srcUtxoHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return err
	}

	srcUtxo := wire.NewOutPoint(srcUtxoHash, 0)
	srcTxIn := wire.NewTxIn(srcUtxo, nil, nil)

	//prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0)) //utxo
	//txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	//fmt.Println(emptyMsg, txIn)
	//var tx = emptyMsg.AddTxIn(txIn)

	destPkScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	srcPkScript, err := txscript.PayToAddrScript(wallet.addressPubKey.AddressPubKeyHash())
	if err != nil {
		fmt.Println(err)
		return err
	}

	srcTxOut := wire.NewTxOut(amount, srcPkScript)
	srcTx.AddTxIn(srcTxIn)
	srcTx.AddTxOut(srcTxOut)
	srcTxHash := srcTx.TxHash()

	redeemTx := wire.NewMsgTx(wire.TxVersion)

	destUtxo := wire.NewOutPoint(&srcTxHash, 0)

	redeemTxIn := wire.NewTxIn(destUtxo, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)
	redeemTxOut := wire.NewTxOut(amount, destPkScript)
	redeemTx.AddTxOut(redeemTxOut)

	bigInt := big.NewInt(amount)
	//Post New TXSkeleton
	skel, err := bc.NewTX(gobcy.TempNewTX(srcAddr.String(), destAddr.String(), *bigInt ), false)
	//Sign it locally
	err = skel.Sign([]string {hex.EncodeToString(wallet.wif.PrivKey.Serialize())})
	if err != nil {
	  fmt.Println(err)
	}
	//Send TXSkeleton
	skel, err = bc.SendTX(skel)
	if err != nil {
	  fmt.Println(err)
	}

	transaction.TxId = srcTxHash.String()
	transaction.SourceAddress = srcAddr.EncodeAddress()
	transaction.DestinationAddress = destAddr.EncodeAddress()
	transaction.Amount = amount
	transaction.UnsignedTx = hex.EncodeToString(srcPkScript)
	transaction.RedeemTx = redeemTx

	

	fmt.Println("tx hash: ", skel.Trans.Hash)
	fmt.Println(skel.Trans.Total.MarshalJSON())
	fmt.Println("block hash", skel.Trans.BlockHash)

	return nil
}

func (transaction *Transaction) SignTx(wallet *Wallet) (string, error) {
	srcPkScript, err := hex.DecodeString(transaction.UnsignedTx)
	if err != nil {
		return "", err
	}
	fmt.Println(transaction.RedeemTx)
	sign, err := txscript.SignatureScript(transaction.RedeemTx, 0, srcPkScript, txscript.SigHashAll, wallet.wif.PrivKey, false)
	if err != nil {
		return "", err
	}

	transaction.RedeemTx.TxIn[0].SignatureScript = sign

	var signedTx bytes.Buffer

	transaction.RedeemTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	transaction.SignedTx = hexSignedTx
	return hexSignedTx, nil
}

func (transaction *Transaction) PushTX(bc gobcy.API) error {
  skel, err := bc.PushTX(transaction.RedeemTx.TxHash().String())
  if err != nil {
    fmt.Println(err)
  }
  fmt.Printf("%+v\n", skel)
  return nil
}