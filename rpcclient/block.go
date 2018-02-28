package rpcclient

import "time"

type Block struct {
	Hash              string   `json:"hash"`
	StrippedSize      uint64   `json:"strippedsize"`
	Size              uint64   `json:"size"`
	Weight            uint64   `json:"weight"`
	Height            uint64   `json:"height"`
	Version           uint64   `json:"version"`
	VersionHex        string   `json:"versionHex"`
	MerkleRoot        string   `json:"merkleroot"`
	TimeInt           int64    `json:"time"`
	MedianTimeInt     int64    `json:"mediantime"`
	Nonce             int64    `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Chainwork         string   `json:"chainwork"`
	PreviousBlockhash string   `json:"previousblockhash"`
	NextBlockhash     string   `json:"nextblockhash"`
	Confirmations     uint64   `json:"confirmations"`
	TrxIDs            []string `json:"tx"`
}

// Time converts the unix time returned by the server into a go time.Time
func (d *Block) Time() time.Time {
	t := time.Unix(int64(d.TimeInt), 0)
	return t
}

// MedianTime converts the unix time returned by the server into a go time.Time
func (d *Block) MedianTime() time.Time {
	t := time.Unix(int64(d.MedianTimeInt), 0)
	return t
}
