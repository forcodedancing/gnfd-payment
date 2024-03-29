package monitor

import (
	"github.com/avast/retry-go/v4"
	"strings"
	"time"
)

const (
	RPCTimeout                   = 10 * time.Second
	SleepSecondForUpdateClient   = 10 * time.Second
	DataSeedDenyServiceThreshold = 60
	FallBehindThreshold          = 5
)

var (
	RtyAttNum = uint(5)
	RtyAttem  = retry.Attempts(RtyAttNum)
	RtyDelay  = retry.Delay(time.Millisecond * 500)
	RtyErr    = retry.LastErrorOnly(true)
)

func escape(raw string) string {
	return strings.ReplaceAll(raw, "'", "\\'")
}
