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
	StateSavePath string
	Epoch int
}{}

var saveCmd = &cli.Command{
	Name: "save",
	Usage: "Given a block_root+epoch, save the corresponding block and state ssz-encoded values to disk. To be used for checkpoint sync.",
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
		&cli.StringFlag{
			Name: "state-root",
			Usage: "instead of epoch, state root (in 0x hex string format) can be used to retrieve from the beacon-node and save locally.",
			Destination: &saveFlags.StateHex,
		},
		&cli.StringFlag{
			Name: "state-save-path",
			Usage: "path to file where state root should be saved if specified. defaults to `state-<state_root>.ssz`",
			Destination: &saveFlags.StateSavePath,
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


	if f.StateHex != "" {
		if err := saveStateByRoot(client, saveFlags.StateHex, saveFlags.StateSavePath); err != nil {
			return err
		}
	} else if f.Epoch > 0 {
		if err := saveStateByEpoch(client, saveFlags.Epoch, saveFlags.StateSavePath); err != nil {
			return err
		}
	}

	return nil
}

func saveStateByRoot(client *openapi.Client, root, path string) error {
	state, err := client.GetStateByRoot(root)
	if err != nil {
		return err
	}
	stateRoot, err := state.HashTreeRoot()
	if err != nil {
		return err
	}
	log.Printf("retrieved state for checkpoint, w/ root=%s", hexutil.Encode(stateRoot[:]))
	if path == "" {
		path = fmt.Sprintf("state-%s.ssz", root)
	}
	log.Printf("saving to %s...", path)
	blockBytes, err := state.MarshalSSZ()
	if err != nil {
		return err
	}
	return os.WriteFile(path, blockBytes, 0644)
}

func saveStateByEpoch(client *openapi.Client, epoch int, path string) error {
	state, err := client.GetStateByEpoch(epoch)
	if err != nil {
		return err
	}
	stateRoot, err := state.HashTreeRoot()
	if err != nil {
		return err
	}
	log.Printf("retrieved state for checkpoint, w/ root=%s", hexutil.Encode(stateRoot[:]))
	if path == "" {
		path = fmt.Sprintf("state-%s.ssz", hexutil.Encode(stateRoot[:]))
	}
	log.Printf("saving to %s...", path)
	blockBytes, err := state.MarshalSSZ()
	if err != nil {
		return err
	}
	return os.WriteFile(path, blockBytes, 0644)
}