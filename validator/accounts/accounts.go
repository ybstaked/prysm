package accounts

import (
	"github.com/prysmaticlabs/prysm/v2/validator/accounts/wallet"
	"github.com/prysmaticlabs/prysm/v2/validator/keymanager"
)

var (
	errKeymanagerNotSupported = "keymanager kind not supported: %s"
	// ErrCouldNotInitializeKeymanager informs about failed keymanager initialization
	ErrCouldNotInitializeKeymanager = "could not initialize keymanager"
)

// Config specifies parameters for accounts commands.
type Config struct {
	Wallet           *wallet.Wallet
	Keymanager       keymanager.IKeymanager
	DeletePublicKeys [][]byte
}
