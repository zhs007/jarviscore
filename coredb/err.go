package coredb

import "errors"

var (
	// ErrNoPrivateKey - no private key
	ErrNoPrivateKey = errors.New("no private key")
	// ErrCoreDBHasNotNode - coredb has not node
	ErrCoreDBHasNotNode = errors.New("coredb has not node")
)
