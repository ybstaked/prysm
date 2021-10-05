package attestations

import (
	"github.com/prysmaticlabs/prysm/v2/beacon-chain/operations/attestations/kv"
)

var _ Pool = (*kv.AttCaches)(nil)
