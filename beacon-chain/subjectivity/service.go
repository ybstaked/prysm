package subjectivity

import "github.com/prysmaticlabs/prysm/shared"

type WeakSyncService struct {

}

func NewWeakSyncService() *WeakSyncService {
	return &WeakSyncService{}
}

func (wss *WeakSyncService) Start() {

}

func (wss *WeakSyncService) Stop() error {
	return nil
}

func (wss *WeakSyncService) Status() error {
	return nil
}

var _ shared.Service = &WeakSyncService{}
