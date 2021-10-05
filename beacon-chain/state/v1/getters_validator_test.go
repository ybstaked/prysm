package v1_test

import (
	"testing"

	v1 "github.com/prysmaticlabs/prysm/v2/beacon-chain/state/v1"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v2/testing/assert"
	"github.com/prysmaticlabs/prysm/v2/testing/require"
)

func TestBeaconState_ValidatorAtIndexReadOnly_HandlesNilSlice(t *testing.T) {
	st, err := v1.InitializeFromProtoUnsafe(&ethpb.BeaconState{
		Validators: nil,
	})
	require.NoError(t, err)

	_, err = st.ValidatorAtIndexReadOnly(0)
	assert.Equal(t, v1.ErrNilValidatorsInState, err)
}
