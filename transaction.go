package main

import "github.com/blockcypher/gobcy"

type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

func getLastUnsignedTx(bc gobcy.API) string {
	txid, _ := bc.GetUnTX()
	return txid[0].Hash
}
