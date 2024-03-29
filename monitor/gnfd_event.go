package monitor

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type Extra struct {
	Category string          `json:"category"`
	Desc     string          `json:"desc"`
	Url      string          `json:"url"`
	Price    decimal.Decimal `json:"price"`
}

// If there is wrong format extra, we just ignore the error and use default values.
// https://gnfd-testnet-fullnode-tendermint-us.bnbchain.org/block_results?height=1658808
func parseExtra(str string) (*Extra, error) {
	var extra Extra
	err := json.Unmarshal([]byte(str), &extra)
	if err != nil {
		return &Extra{
			Desc:  "",
			Url:   "",
			Price: decimal.Zero,
		}, nil
	}

	return &extra, nil
}
