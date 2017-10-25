// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package blockchain

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/HcashOrg/hcashd/chaincfg"
	"github.com/HcashOrg/hcashd/chaincfg/chainhash"
	"github.com/HcashOrg/hcashd/database"
	"github.com/HcashOrg/hcashd/txscript"
	"github.com/HcashOrg/hcashd/wire"
	"github.com/HcashOrg/hcashutil"

	_ "github.com/HcashOrg/hcashd/database/ffldb"
)

// recalculateMsgBlockMerkleRootsSize recalculates the merkle roots for a msgBlock,
// then stores them in the msgBlock's header. It also updates the block size.
func recalculateMsgBlockMerkleRootsSize(msgBlock *wire.MsgBlock) {
	tempBlock := hcashutil.NewBlock(msgBlock)

	// adapt for new version
	merkles := BuildMerkleTreeStore(tempBlock.Transactions(), false)
	merklesStake := BuildMerkleTreeStore(tempBlock.STransactions(), false)

	msgBlock.Header.MerkleRoot = *merkles[len(merkles)-1]
	msgBlock.Header.StakeRoot = *merklesStake[len(merklesStake)-1]
	msgBlock.Header.Size = uint32(msgBlock.SerializeSize())
}

// filesExists returns whether or not the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// isSupportedDbType returns whether or not the passed database type is
// currently supported.
func isSupportedDbType(dbType string) bool {
	supportedDrivers := database.SupportedDrivers()
	for _, driver := range supportedDrivers {
		if dbType == driver {
			return true
		}
	}

	return false
}

//  SetupTestChain is used to create a new db and chain instance with the genesis
// block already inserted.  In addition to the new chain instance, it returns
// a teardown function the caller should invoke when done testing to clean up.
func SetupTestChain(dbName string, params *chaincfg.Params) (*BlockChain, func(), error) {
	if !isSupportedDbType(testDbType) {
		return nil, nil, fmt.Errorf("unsupported db type %v", testDbType)
	}

	// Handle memory database specially since it doesn't need the disk
	// specific handling.
	var db database.DB
	var teardown func()
	if testDbType == "memdb" {
		ndb, err := database.Create(testDbType)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating db: %v", err)
		}
		db = ndb

		// Setup a teardown function for cleaning up.  This function is
		// returned to the caller to be invoked when it is done testing.
		teardown = func() {
			db.Close()
		}
	} else {
		// Create the root directory for test databases.
		if !fileExists(testDbRoot) {
			if err := os.MkdirAll(testDbRoot, 0700); err != nil {
				err := fmt.Errorf("unable to create test db "+
					"root: %v", err)
				return nil, nil, err
			}
		}

		// Create a new database to store the accepted blocks into.
		dbPath := filepath.Join(testDbRoot, dbName)
		_ = os.RemoveAll(dbPath)
		ndb, err := database.Create(testDbType, dbPath, blockDataNet)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating db: %v", err)
		}
		db = ndb

		// Setup a teardown function for cleaning up.  This function is
		// returned to the caller to be invoked when it is done testing.
		teardown = func() {
			db.Close()
			os.RemoveAll(dbPath)
			os.RemoveAll(testDbRoot)
		}
	}

	// Copy the chain params to ensure any modifications the tests do to
	// the chain parameters do not affect the global instance.
	paramsCopy := *params

	// Create the main chain instance.
	chain, err := New(&Config{
		DB:          db,
		ChainParams: &paramsCopy,
		TimeSource:  NewMedianTime(),
		SigCache:    txscript.NewSigCache(1000),
	})

	if err != nil {
		teardown()
		err := fmt.Errorf("failed to create chain instance: %v", err)
		return nil, nil, err
	}

	return chain, teardown, nil
}

// FakeBlock fakes a block of specific version, height and stake version
func FakeBlock(blockVersion int32, height int64, stakeVersion uint32) *hcashutil.Block {
	// Make up a header.
	header := wire.BlockHeader{
		Version:      blockVersion,
		Height:       uint32(height),
		Nonce:        0,
		StakeVersion: stakeVersion,
	}

	msgBlock := &wire.MsgBlock{
		Header: header,
	}

	return hcashutil.NewBlock(msgBlock)
}

func FakeBlockFromHeader(header *wire.BlockHeader) *hcashutil.Block {
	msgBlock := &wire.MsgBlock{
		Header: *header,
	}

	return hcashutil.NewBlock(msgBlock)
}

// CheckBlockHeaderContextEx makes the internal checkBlockHeaderContext
// function available to the test package.
func (b *BlockChain) CheckBlockHeaderContextEx(header *wire.BlockHeader, prevNode *blockNode, flags BehaviorFlags) error {
	return b.checkBlockHeaderContext(header, prevNode, flags)
}

// NewBlockNodeEx makes the internal newBlockNode function available to the
// test package.
func NewBlockNodeEx(blockHeader *wire.BlockHeader, ticketsSpent []chainhash.Hash, ticketsRevoked []chainhash.Hash, voteBits []VoteVersionTuple) *blockNode {
	return newBlockNode(FakeBlockFromHeader(blockHeader), ticketsSpent, ticketsRevoked, voteBits)
}

// ToTimeSorter converts a timestamp array to a sortable type
func ToTimeSorter(time_series []time.Time) sort.Interface {
	return timeSorter(time_series)
}

// SetMaxMedianTimeEntries exports the ability to set the max number
// of median time entries available to the blockchain_test package
func SetMaxMedianTimeEntries(val int) {
	maxMedianTimeEntries = val
}
