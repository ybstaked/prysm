package checkpoint

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"net"
	"net/url"
	"time"

	"github.com/prysmaticlabs/prysm/api/client/openapi"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var flags = struct{
	SaveFlags struct{
		BeaconNodeHost string
		Timeout        string
	}
}{}

var saveCmd *cli.Command = &cli.Command{
	Name: "save",
	Usage: "Connect to a beacon-node server and save checkpoint data",
	Action: cliActionSave,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "beacon-node-host",
			Usage: "host:port for beacon node to query for checkpoint state",
			Destination: &flags.SaveFlags.BeaconNodeHost,
		},
		&cli.StringFlag{
			Name: "http-timeout",
			Usage: "timeout for http requests made to beacon-node-url (uses duration format, ex: 2m31s). default: 2m",
			Destination: &flags.SaveFlags.Timeout,
			Value: "2m",
		},
	},
}

func cliActionSave(c *cli.Context) error {
	f := flags.SaveFlags
	opts := make([]openapi.ClientOpt, 0)
	log.Printf("beacon-node-url=%s", f.BeaconNodeHost)
	timeout, err := time.ParseDuration(f.Timeout)
	if err != nil {
		return err
	}
	opts = append(opts, openapi.WithTimeout(timeout))
	validatedHost, err := validHostname(flags.SaveFlags.BeaconNodeHost)
	if err !=  nil {
		return err
	}
	log.Printf("host-port=%s", validatedHost)
	client := openapi.NewClient(validatedHost, opts...)
	wsc, err := client.GetWeakSubjectivityCheckpoint()
	if err != nil {
		return err
	}
	log.Printf("epoch:%d\nblock_root:%s\nstate_root:%s\n", int(wsc.Epoch), hexutil.Encode(wsc.BlockRoot), hexutil.Encode(wsc.StateRoot))
	return nil
}

func validHostname(h string) (string, error){
	// try to parse as url (being permissive)
	u, err := url.Parse(h)
	if err == nil && u.Host != "" {
		return u.Host, nil
	}
	// try to parse as host:port
	host, port, err := net.SplitHostPort(h)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", host, port), nil
}