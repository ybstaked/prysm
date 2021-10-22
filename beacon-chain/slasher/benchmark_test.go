package slasher

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	types "github.com/prysmaticlabs/eth2-types"
	dbtest "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	slashertypes "github.com/prysmaticlabs/prysm/beacon-chain/slasher/types"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/io/file"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/attestation"
	"github.com/prysmaticlabs/prysm/testing/require"
	"github.com/prysmaticlabs/prysm/time/slots"
)

func TestWriteAttestations(t *testing.T) {
	params.UsePraterConfig()
	params.UsePraterNetworkConfig()
	enc, err := file.ReadFileAsBytes("/tmp/blockdata/block.ssz")
	require.NoError(t, err)
	blk := &ethpb.SignedBeaconBlockAltair{}
	require.NoError(t, blk.UnmarshalSSZ(enc))
	ctx := context.Background()
	outputDirectory := "/tmp/attestations"
	for i, att := range blk.Block.Body.Attestations {
		// Read the committee for the attestation slot from disk.
		committeesFromDisk, err := readCommitteesFromDisk(att.Data.Slot)
		require.NoError(t, err)
		if uint64(att.Data.CommitteeIndex) >= uint64(len(committeesFromDisk)) {
			t.Fatalf(
				"Committee index %d bigger than the number of committees in the slot %d",
				att.Data.CommitteeIndex,
				len(committeesFromDisk),
			)
		}
		committee := committeesFromDisk[att.Data.CommitteeIndex]
		// Convert the attestation to indexed form.
		idxAtt, err := attestation.ConvertToIndexed(ctx, att, committee)
		require.NoError(t, err)
		encodedAtt, err := idxAtt.MarshalSSZ()
		require.NoError(t, err)
		f, err := os.Create(filepath.Join(outputDirectory, fmt.Sprintf("%d.ssz", i)))
		require.NoError(t, err)
		_, err = f.Write(encodedAtt)
		require.NoError(t, err)
		require.NoError(t, f.Close())
	}
}

func TestSlasherTimes(t *testing.T) {
	params.UsePraterConfig()
	params.UsePraterNetworkConfig()
	enc, err := file.ReadFileAsBytes("/tmp/blockdata/block.ssz")
	require.NoError(t, err)
	blk := &ethpb.SignedBeaconBlockAltair{}
	require.NoError(t, blk.UnmarshalSSZ(enc))
	slasherDB := dbtest.SetupSlasherDB(t)
	ctx := context.Background()

	// Converting the attestations from the block into indexed format.
	indexedAttWrappers := make([]*slashertypes.IndexedAttestationWrapper, len(blk.Block.Body.Attestations))
	for i, att := range blk.Block.Body.Attestations {
		// Read the committee for the attestation slot from disk.
		committeesFromDisk, err := readCommitteesFromDisk(att.Data.Slot)
		require.NoError(t, err)
		if uint64(att.Data.CommitteeIndex) >= uint64(len(committeesFromDisk)) {
			t.Fatalf(
				"Committee index %d bigger than the number of committees in the slot %d",
				att.Data.CommitteeIndex,
				len(committeesFromDisk),
			)
		}
		committee := committeesFromDisk[att.Data.CommitteeIndex]
		// Convert the attestation to indexed form.
		idxAtt, err := attestation.ConvertToIndexed(ctx, att, committee)
		require.NoError(t, err)
		signingRoot, err := att.Data.HashTreeRoot()
		require.NoError(t, err)
		indexedAttWrappers[i] = &slashertypes.IndexedAttestationWrapper{
			IndexedAttestation: idxAtt,
			SigningRoot:        signingRoot,
		}
	}
	t.Log("Got indexed atts", len(indexedAttWrappers))

	// Initializing the service
	srv, err := New(ctx, &ServiceConfig{
		Database: slasherDB,
	})
	require.NoError(t, err)
	require.NoError(t, srv.serviceCfg.Database.SaveAttestationRecordsForValidators(ctx, indexedAttWrappers))

	// Set the current epoch to the epoch the block was extracted from + 1.
	currentEpoch := slots.ToEpoch(1512693) + 1
	fmt.Println(currentEpoch)
	_, err = srv.checkSlashableAttestations(ctx, currentEpoch, indexedAttWrappers)
	require.NoError(t, err)
}

func readCommitteesFromDisk(slot types.Slot) ([][]types.ValidatorIndex, error) {
	f, err := os.Open(fmt.Sprintf("/tmp/blockdata/committees/%d.txt", slot))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	committees := make([][]types.ValidatorIndex, 0)
	for _, ln := range lines {
		valIndicesStr := strings.Split(ln, ",")
		valIndices := make([]types.ValidatorIndex, len(valIndicesStr))
		for i, vStr := range valIndicesStr {
			valIdx, err := strconv.ParseUint(vStr, 10, 64)
			if err != nil {
				return nil, err
			}
			valIndices[i] = types.ValidatorIndex(valIdx)
		}
		committees = append(committees, valIndices)
	}
	return committees, nil
}
