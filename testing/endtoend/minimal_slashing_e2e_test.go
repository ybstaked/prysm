package endtoend

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v2/config/params"
	ev "github.com/prysmaticlabs/prysm/v2/testing/endtoend/evaluators"
	e2eParams "github.com/prysmaticlabs/prysm/v2/testing/endtoend/params"
	"github.com/prysmaticlabs/prysm/v2/testing/endtoend/types"
	"github.com/prysmaticlabs/prysm/v2/testing/require"
)

func TestEndToEnd_Slasher_MinimalConfig(t *testing.T) {
	params.UseE2EConfig()
	require.NoError(t, e2eParams.Init(e2eParams.StandardBeaconCount))

	testConfig := &types.E2EConfig{
		BeaconFlags: []string{
			"--slasher",
		},
		ValidatorFlags: []string{},
		EpochsToRun:    4,
		TestSync:       false,
		TestDeposits:   false,
		Evaluators: []types.Evaluator{
			ev.PeersConnect,
			ev.HealthzCheck,
			ev.ValidatorsSlashedAfterEpoch(4),
			ev.SlashedValidatorsLoseBalanceAfterEpoch(4),
			ev.InjectDoubleVoteOnEpoch(2),
			ev.InjectDoubleBlockOnEpoch(2),
		},
	}

	newTestRunner(t, testConfig).run()
}
