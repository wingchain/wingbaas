/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fileledger

import (
	"sync"

	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/common/ledger/blkstorage"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/common/ledger/blkstorage/fsblkstorage"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/common/ledger/blockledger"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/common/metrics"
)

type fileLedgerFactory struct {
	blkstorageProvider blkstorage.BlockStoreProvider
	ledgers            map[string]blockledger.ReadWriter
	mutex              sync.Mutex
}

// GetOrCreate gets an existing ledger (if it exists) or creates it if it does not
func (flf *fileLedgerFactory) GetOrCreate(chainID string) (blockledger.ReadWriter, error) {
	flf.mutex.Lock()
	defer flf.mutex.Unlock()

	key := chainID
	// check cache
	ledger, ok := flf.ledgers[key]
	if ok {
		return ledger, nil
	}
	// open fresh
	blockStore, err := flf.blkstorageProvider.OpenBlockStore(key)
	if err != nil {
		return nil, err
	}
	ledger = NewFileLedger(blockStore)
	flf.ledgers[key] = ledger
	return ledger, nil
}

// ChannelIDs returns the channel IDs the factory is aware of
func (flf *fileLedgerFactory) ChannelIDs() []string {
	channelIDs, err := flf.blkstorageProvider.List()
	if err != nil {
		logger.Panic(err)
	}
	return channelIDs
}

// Close releases all resources acquired by the factory
func (flf *fileLedgerFactory) Close() {
	flf.blkstorageProvider.Close()
}

// New creates a new ledger factory
func New(directory string, metricsProvider metrics.Provider) (blockledger.Factory, error) {
	p, err := fsblkstorage.NewProvider(
		fsblkstorage.NewConf(directory, -1),
		&blkstorage.IndexConfig{
			AttrsToIndex: []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockNum}},
		metricsProvider,
	)
	if err != nil {
		return nil, err
	}
	return &fileLedgerFactory{
		blkstorageProvider: p,
		ledgers:            make(map[string]blockledger.ReadWriter),
	}, nil
}