package rpc

import (
	pb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1/validator-client"
)

var _ pb.AuthServer = (*Server)(nil)
