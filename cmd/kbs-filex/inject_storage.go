package main

import (
	"github.com/google/wire"
	"github.com/kabaliserv/filex/storage"
)

var storageSet = wire.NewSet(
	storage.New,
)
