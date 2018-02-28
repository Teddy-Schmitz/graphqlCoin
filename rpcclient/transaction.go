package rpcclient

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID            string   `json:"txid"`
	Hash          string   `json:"hash"`
	Version       uint64   `json:"version"`
	Size          uint64   `json:"size"`
	Vsize         uint64   `json:"vsize"`
	Locktime      uint64   `json:"locktime"`
	BlockHash     string   `json:"blockhash"`
	TimeInt       uint64   `json:"time"`
	BlockTimeInt  uint64   `json:"blocktime"`
	Confirmations uint64   `json:"confirmations"`
	VIn           []*TxIn  `json:"vin"`
	VOut          []*TxOut `json:"vout"`
}

// Time converts the unix time returned by the server into a go time.Time
func (t *Transaction) Time() time.Time {
	ti := time.Unix(int64(t.TimeInt), 0)
	return ti
}

// BlockTime converts the unix time returned by the server into a go time.Time
func (t *Transaction) BlockTime() time.Time {
	ti := time.Unix(int64(t.BlockTimeInt), 0)
	return ti
}

type TxIn struct {
	Coinbase              *string       `json:"coinbase,omitempty" db:"coinbase"`
	Sequence              int64         `json:"sequence" db:"sequence"`
	ScriptSig             PubKeyScripts `json:"scriptSig,omitempty"`
	PreviousTransactionID *string       `json:"txid,omitempty" db:"previous_transaction_id"`
	OutIndex              *int64        `json:"vout" db:"out_index"`
}

type TxOut struct {
	Value         decimal.Decimal `json:"value" db:"value"`
	Index         int64           `json:"n" db:"index"`
	PubKeyScripts PubKeyScripts   `json:"scriptPubKey,omitempty"`
}

type PubKeyScripts struct {
	ASM       *string    `json:"asm,omitempty"`
	HEX       *string    `json:"hex,omitempty"`
	Type      string     `json:"type,omitempty"`
	ReqSigs   int64      `json:"reqSigs,omitempty"`
	Addresses *[]*string `json:"addresses,omitempty"`
}
