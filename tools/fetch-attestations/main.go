package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	types "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/attestation"
	"github.com/prysmaticlabs/prysm/time/slots"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	outputFlag = flag.String("output-dir", "/tmp/attestations", "output dir path for ssz contents")
	slotFlag   = flag.Uint64("slot", 0, "slot to fetch")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := ethpb.NewBeaconChainClient(conn)
	ctx := context.Background()
	slot := types.Slot(*slotFlag)
	if slot == 0 {
		log.Fatal("Please specify a -slot")
	}
	epoch := slots.ToEpoch(slot)
	log.Info("Fetching beacon block")
	blks, err := client.ListBeaconBlocks(ctx, &ethpb.ListBlocksRequest{
		QueryFilter: &ethpb.ListBlocksRequest_Slot{Slot: slot},
		PageSize:    1,
		PageToken:   "",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Fetching committees")
	committees, err := client.ListBeaconCommittees(ctx, &ethpb.ListCommitteesRequest{
		QueryFilter: &ethpb.ListCommitteesRequest_Epoch{Epoch: epoch},
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(blks.BlockContainers) != 1 {
		log.Fatal("Did not receive single block")
	}
	altairBlk, ok := blks.BlockContainers[0].Block.(*ethpb.BeaconBlockContainer_AltairBlock)
	if !ok {
		log.Fatal("Not altair")
	}
	committeesBySlot := committees.Committees
	log.Info("Writing indexed attestations")
	for i, att := range altairBlk.AltairBlock.Block.Body.Attestations {
		committees, ok := committeesBySlot[uint64(att.Data.Slot)]
		if !ok {
			log.Fatalf("Committee not found for attestation with slot %d", att.Data.Slot)
		}
		if uint64(att.Data.CommitteeIndex) >= uint64(len(committees.Committees)) {
			log.Fatal("Committee with index does not exist")
		}
		comm := committees.Committees[att.Data.CommitteeIndex]
		idxAtt, err := attestation.ConvertToIndexed(ctx, att, comm.ValidatorIndices)
		if err != nil {
			log.Fatal(err)
		}
		if err := writeIndexedAttestation(i, idxAtt); err != nil {
			log.Fatal(err)
		}
	}
	log.Info("Num atts", len(altairBlk.AltairBlock.Block.Body.Attestations))
}

func writeIndexedAttestation(i int, att *ethpb.IndexedAttestation) error {
	attPath := filepath.Join(*outputFlag, fmt.Sprintf("%d.ssz", i))
	f, err := os.Create(attPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error(err)
		}
	}()
	enc, err := att.MarshalSSZ()
	if err != nil {
		return err
	}
	if _, err = f.Write(enc); err != nil {
		return err
	}
	return nil
}
