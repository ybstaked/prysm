package checkpoint

import (
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/api/client/openapi"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var saveFlags = struct{
	BeaconNodeHost string
	Timeout        string
	BlockHex string
	BlockSavePath string
	StateHex string
	Epoch int
}{}

var saveCmd = &cli.Command{
	Name: "save",
	Usage: "query for the current weak subjectivity period epoch, then download the corresponding state and block. To be used for checkpoint sync.",
	Action: cliActionSave,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "beacon-node-host",
			Usage: "host:port for beacon node connection",
			Destination: &saveFlags.BeaconNodeHost,
		},
		&cli.StringFlag{
			Name: "http-timeout",
			Usage: "timeout for http requests made to beacon-node-url (uses duration format, ex: 2m31s). default: 2m",
			Destination: &saveFlags.Timeout,
			Value: "2m",
		},
		&cli.IntFlag{
			Name: "epoch",
			Usage: "instead of state-root, epoch can be used to find the BeaconState for the slot at the epoch boundary.",
			Destination: &saveFlags.Epoch,
		},
	},
}

func cliActionSave(c *cli.Context) error {
	f := saveFlags
	opts := make([]openapi.ClientOpt, 0)
	log.Printf("beacon-node-url=%s", f.BeaconNodeHost)
	timeout, err := time.ParseDuration(f.Timeout)
	if err != nil {
		return err
	}
	opts = append(opts, openapi.WithTimeout(timeout))
	client, err := openapi.NewClient(saveFlags.BeaconNodeHost, opts...)
	if err !=  nil {
		return err
	}

	if saveFlags.Epoch > 0 {
		return saveCheckpointByEpoch(client, uint64(saveFlags.Epoch))
	}

	return saveCheckpoint(client)
}

func saveCheckpoint(client *openapi.Client) error {
	epoch, err := client.GetWeakSubjectivityCheckpointEpoch()
	if err != nil {
		return err
	}

	log.Printf("Beacon node computes the current weak subjectivity checkpoint as epoch = %d", epoch)
	return saveCheckpointByEpoch(client, epoch)
}

func saveCheckpointByEpoch(client *openapi.Client, epoch uint64) error {
	state, err := client.GetStateByEpoch(epoch)
	if err != nil {
		return err
	}
	stateRoot, err := state.HashTreeRoot()
	if err != nil {
		return err
	}
	log.Printf("retrieved state for checkpoint, w/ root=%s", hexutil.Encode(stateRoot[:]))

	sb, err := state.MarshalSSZ()
	if err != nil {
		return err
	}

	statePath := fmt.Sprintf("state-%s.ssz", hexutil.Encode(stateRoot[:]))
	log.Printf("saving ssz-encoded state to to %s", statePath)
	err = os.WriteFile(statePath, sb, 0644)
	if err != nil {
		return err
	}

	blockRoot, err := state.LatestBlockHeader.HashTreeRoot()
	if err != nil {
		return err
	}
	blockRootHex := hexutil.Encode(blockRoot[:])

	block, err := client.GetBlockByRoot(blockRootHex)
	log.Printf("retrieved block by root=%s", hexutil.Encode(blockRoot[:]))

	blockPath := fmt.Sprintf("block-%s.ssz", blockRootHex)
	bb, err := block.MarshalSSZ()
	if err != nil {
		return err
	}
	log.Printf("saving ssz-encoded block to to %s", statePath)

	err = os.WriteFile(blockPath, bb, 0644)
	if err != nil {
		return err
	}

	fmt.Println("To validate that your client is using this checkpoint, specify the following flag when starting prysm:")
	fmt.Printf("--weak-subjectivity-checkpoint=%s:%d\n\n", blockRootHex, epoch)
	fmt.Println("To sync a new beacon node starting from the checkpoint state, you may specify the following flags (assuming the files are in your current working directory)")
	fmt.Printf("--checkpoint-state=%s --checkpoint-block=%s\n", statePath, blockPath)
	return nil
}