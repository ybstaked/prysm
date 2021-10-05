package fuzz

import (
	"bytes"

	"github.com/prysmaticlabs/prysm/v2/beacon-chain/p2p/encoder"
	"github.com/prysmaticlabs/prysm/v2/config/params"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
)

var buf = new(bytes.Buffer)

// SszEncoderAttestationFuzz runs network encode/decode for attestations.
func SszEncoderAttestationFuzz(b []byte) {
	params.UseMainnetConfig()
	buf.Reset()
	input := &ethpb.Attestation{}
	e := encoder.SszNetworkEncoder{}
	if err := e.DecodeGossip(b, input); err != nil {
		_ = err
		return
	}
	if _, err := e.EncodeGossip(buf, input); err != nil {
		_ = err
		return
	}
}
