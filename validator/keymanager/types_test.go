package keymanager_test

import (
	"github.com/prysmaticlabs/prysm/v2/validator/keymanager"
	"github.com/prysmaticlabs/prysm/v2/validator/keymanager/derived"
	"github.com/prysmaticlabs/prysm/v2/validator/keymanager/imported"
	"github.com/prysmaticlabs/prysm/v2/validator/keymanager/remote"
)

var (
	_ = keymanager.IKeymanager(&imported.Keymanager{})
	_ = keymanager.IKeymanager(&derived.Keymanager{})
	_ = keymanager.IKeymanager(&remote.Keymanager{})
)
