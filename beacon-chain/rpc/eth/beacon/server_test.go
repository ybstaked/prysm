package beacon

import ethpbservice "github.com/prysmaticlabs/prysm/v2/proto/eth/service"

var _ ethpbservice.BeaconChainServer = (*Server)(nil)
