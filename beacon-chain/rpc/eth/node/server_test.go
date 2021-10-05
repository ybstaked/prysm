package node

import (
	ethpbservice "github.com/prysmaticlabs/prysm/v2/proto/eth/service"
)

var _ ethpbservice.BeaconNodeServer = (*Server)(nil)
