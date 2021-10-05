package operations

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v2/testing/spectest/shared/altair/operations"
)

func TestMinimal_Altair_Operations_AttesterSlashing(t *testing.T) {
	operations.RunAttesterSlashingTest(t, "minimal")
}
