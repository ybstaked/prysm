package openapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	types "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
)

const GET_WEAK_SUBJECTIVITY_CHECKPOINT_PATH = "/eth/v1alpha1/beacon/weak_subjectivity_checkpoint"

type ClientOpt func(*Client)

func WithTimeout(timeout time.Duration) ClientOpt {
	return func(c *Client) {
		c.c.Timeout = timeout
	}
}

type Client struct {
	c *http.Client
	host string
	scheme string
}

func (c *Client) urlForPath(methodPath string) *url.URL {
	u := &url.URL{
		Scheme: c.scheme,
		Host: c.host,
	}
	u.Path = path.Join(u.Path, methodPath)
	return u
}

func NewClient(host string, opts ...ClientOpt) *Client {
	c := &Client{
		c: &http.Client{},
		scheme: "http",
		host: host,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type WSCResponse struct {
	BlockRoot string
	StateRoot string
	Epoch string
}

func (c *Client) GetWeakSubjectivityCheckpoint() (*ethpb.WeakSubjectivityCheckpoint, error) {
	u := c.urlForPath(GET_WEAK_SUBJECTIVITY_CHECKPOINT_PATH)
	r, err := c.c.Get(u.String())
	if err != nil {
		return nil, err
	}
	v := &WSCResponse{}
	b := bytes.NewBuffer(nil)
	bodyReader := io.TeeReader(r.Body, b)
	log.Printf("%s\n", b.String())
	err = json.NewDecoder(bodyReader).Decode(v)
	if err != nil {
		return nil, err
	}
	epoch, err := strconv.ParseUint(v.Epoch, 10, 64)
	if err != nil {
		return nil, err
	}
	blockRoot, err := base64.StdEncoding.DecodeString(v.BlockRoot)
	if err != nil {
		return nil, err
	}
	stateRoot, err := base64.StdEncoding.DecodeString(v.StateRoot)
	if err != nil {
		return nil, err
	}
	return &ethpb.WeakSubjectivityCheckpoint{
		Epoch: types.Epoch(epoch),
		BlockRoot: blockRoot,
		StateRoot: stateRoot,
	}, nil
}
