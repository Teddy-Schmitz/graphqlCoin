package rpcclient

type FeeEstimate struct {
	FeeRate float64 `json:"feerate"`
	Blocks  int     `json:"blocks"`
}
