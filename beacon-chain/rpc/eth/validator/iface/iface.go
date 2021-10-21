package iface

import (
	"context"

	ethpbv1 "github.com/prysmaticlabs/prysm/proto/eth/v1"
	ethpbv2 "github.com/prysmaticlabs/prysm/proto/eth/v2"
)

type ServerIface interface {
	ProduceBlockV2(ctx context.Context, req *ethpbv1.ProduceBlockRequest) (*ethpbv2.ProduceBlockResponseV2, error)
}
