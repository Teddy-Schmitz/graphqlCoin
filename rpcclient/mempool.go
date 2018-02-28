package rpcclient

import "time"

type MemPoolTrx struct {
	ID              string
	Size            uint64   `json:"size"`
	Fee             float64  `json:"fee"`
	ModifiedFee     float64  `json:"modifiedfee"`
	TimeInt         int64    `json:"time"`
	Height          uint64   `json:"height"`
	DescendantCount uint64   `json:"descendantcount"`
	DescendantSize  uint64   `json:"descendantsize"`
	DescendantFees  uint64   `json:"descendantfees"`
	AncestorCount   uint64   `json:"ancestorcount"`
	AncestorSize    uint64   `json:"ancestorsize"`
	AncestorFees    uint64   `json:"ancestorfees"`
	Depends         []string `json:"depends"`
}

// Time converts the unix time returned by the server into a go time.Time
func (m *MemPoolTrx) Time() time.Time {
	return time.Unix(m.TimeInt, 0)
}
