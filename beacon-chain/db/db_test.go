package db

import "github.com/prysmaticlabs/prysm/v2/beacon-chain/db/kv"

var _ Database = (*kv.Store)(nil)
