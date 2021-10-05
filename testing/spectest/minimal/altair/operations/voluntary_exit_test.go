package operations

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v2/testing/spectest/shared/altair/operations"
)

func TestMinimal_Altair_Operations_VoluntaryExit(t *testing.T) {
	operations.RunVoluntaryExitTest(t, "minimal")
}
