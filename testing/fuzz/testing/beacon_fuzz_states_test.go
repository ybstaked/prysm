package testing

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v2/testing/require"
)

func TestGetBeaconFuzzState(t *testing.T) {
	_, err := BeaconFuzzState(1)
	require.NoError(t, err)
}
