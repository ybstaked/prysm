package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/config/params"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/time/slots"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	outputFlag = flag.String("output-dir", "/tmp/blockdata", "output dir path for ssz contents")
	slotFlag   = flag.Uint64("slot", 0, "slot to fetch")
)

func main() {
	flag.Parse()
	params.UsePraterNetworkConfig()
	params.UsePraterConfig()
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
	enc, err := altairBlk.AltairBlock.MarshalSSZ()
	if err != nil {
		log.Fatal(err)
	}
	blockPath := filepath.Join(*outputFlag, "block.ssz")
	f, err := os.Create(blockPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	n, err := f.Write(enc)
	if err != nil {
		log.Fatal(err)
	}
	if n != len(enc) {
		log.Fatal("Did not write all bytes to file")
	}
	committeesBySlot := committees.Committees
	log.Info("Writing committee files")
	uniqueSlots := make(map[types.Slot]bool)
	// Check every attestation has a committee associated.
	for _, att := range altairBlk.AltairBlock.Block.Body.Attestations {
		comm, ok := committeesBySlot[uint64(att.Data.Slot)]
		if !ok {
			log.Fatalf("Committee not found for attestation with slot %d", att.Data.Slot)
		}
		uniqueSlots[att.Data.Slot] = true
		if err := writeCommitteeFile(att.Data.Slot, comm.Committees); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("num atts", len(altairBlk.AltairBlock.Block.Body.Attestations))
	for slt := range uniqueSlots {
		fmt.Println(slt)
	}
}

func writeCommitteeFile(slot types.Slot, comms []*ethpb.BeaconCommittees_CommitteeItem) error {
	commPath := filepath.Join(*outputFlag, "committees", fmt.Sprintf("%d.txt", slot))
	// Create a writer
	f, err := os.Create(commPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	w := bufio.NewWriter(f)
	for _, comm := range comms {
		if _, err := w.WriteString(arrayToString(comm.ValidatorIndices, ",") + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	return w.Flush()
}

func arrayToString(a []types.ValidatorIndex, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
