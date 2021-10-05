package testing

import (
	"github.com/prysmaticlabs/prysm/v2/time/slots"
)

var _ slots.Ticker = (*MockTicker)(nil)
