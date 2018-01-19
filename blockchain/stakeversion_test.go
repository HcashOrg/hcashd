// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"fmt"
	"testing"
	"time"

	//	"github.com/HcashOrg/hcashd/blockchain"
	"github.com/HcashOrg/hcashd/chaincfg"
	"github.com/HcashOrg/hcashd/chaincfg/chainhash"
	"github.com/HcashOrg/hcashutil"
)

// newFakeChain returns a chain that is usable for syntetic tests.
func newFakeChain(params *chaincfg.Params) *BlockChain {
	return &BlockChain{
		chainParams:      params,
		deploymentCaches: newThresholdCaches(params),
		index:            make(map[chainhash.Hash]*blockNode),
		isVoterMajorityVersionCache:   make(map[[stakeMajorityCacheKeySize]byte]bool),
		isStakeMajorityVersionCache:   make(map[[stakeMajorityCacheKeySize]byte]bool),
		calcPriorStakeVersionCache:    make(map[[chainhash.HashSize]byte]uint32),
		calcVoterVersionIntervalCache: make(map[[chainhash.HashSize]byte]uint32),
		calcStakeVersionCache:         make(map[[chainhash.HashSize]byte]uint32),
	}
}

// genesisBlockNode creates a fake chain of blockNodes.  It is used for testing
// the mechanical properties of the version code.
func genesisBlockNode(params *chaincfg.Params) *blockNode {
	// Create a new node from the genesis block.
	genesisBlock := hcashutil.NewBlock(params.GenesisBlock)

	node := newBlockNode(genesisBlock, nil, nil, nil)
	node.inMainChain = true

	return node
}

func TestCalcWantHeight(t *testing.T) {
	// For example, if StakeVersionInterval = 11 and StakeValidationHeight = 13 the
	// windows start at 13 + (11 * 2) 25 and are as follows: 24-34, 35-45, 46-56 ...
	// If height comes in at 35 we use the 24-34 window, up to height 45.
	// If height comes in at 46 we use the 35-45 window, up to height 56 etc.
	tests := []struct {
		name       string
		skip       int64
		interval   int64
		multiplier int64
		negative   int64
	}{
		{
			name:       "13 11 10000",
			skip:       13,
			interval:   11,
			multiplier: 10000,
		},
		{
			name:       "27 33 10000",
			skip:       27,
			interval:   33,
			multiplier: 10000,
		},
		{
			name:       "mainnet params",
			skip:       chaincfg.MainNetParams.StakeValidationHeight,
			interval:   chaincfg.MainNetParams.StakeVersionInterval,
			multiplier: 5000,
		},
		{
			name:       "testnet2 params",
			skip:       chaincfg.TestNet2Params.StakeValidationHeight,
			interval:   chaincfg.TestNet2Params.StakeVersionInterval,
			multiplier: 1000,
		},
		{
			name:       "simnet params",
			skip:       chaincfg.SimNetParams.StakeValidationHeight,
			interval:   chaincfg.SimNetParams.StakeVersionInterval,
			multiplier: 10000,
		},
		{
			name:       "negative mainnet params",
			skip:       chaincfg.MainNetParams.StakeValidationHeight,
			interval:   chaincfg.MainNetParams.StakeVersionInterval,
			multiplier: 1000,
			negative:   1,
		},
	}

	for _, test := range tests {
		t.Logf("running: %v skip: %v interval: %v",
			test.name, test.skip, test.interval)

		start := int64(test.skip + test.interval*2)
		expectedHeight := start - 1 // zero based
		x := int64(0) + test.negative
		for i := start; i < test.multiplier*test.interval; i++ {
			if x%test.interval == 0 && i != start {
				expectedHeight += test.interval
			}
			wantHeight := calcWantHeight(test.skip, test.interval, i)

			if wantHeight != expectedHeight {
				if test.negative == 0 {
					t.Fatalf("%v: i %v x %v -> wantHeight %v expectedHeight %v\n",
						test.name, i, x, wantHeight, expectedHeight)
				}
			}

			x++
		}
	}
}

// newFakeNode creates a fake blockNode and sets pertinent internals.
func newFakeNode(blockVersion int32, height int64, currentNode *blockNode) *blockNode {
	// Make up a header.
	/*
		header := &wire.BlockHeader{
			Version: blockVersion,
			Height:  uint32(height),
			Nonce:   0,
		}
	*/

	//node := newBlockNode(header, nil, nil, nil)
	node := newBlockNode(FakeBlock(blockVersion, height, 0), nil, nil, nil)
	node.height = height
	node.parent = currentNode

	return node
}

// TestCalcStakeVersionCorners doesn't work yet
func DNWTestCalcStakeVersionCorners(t *testing.T) {
	params := &chaincfg.SimNetParams
	currentNode := genesisBlockNode(params)

	bc := newFakeChain(params)

	svh := params.StakeValidationHeight
	interval := params.StakeVersionInterval

	height := int64(0)
	for i := int64(1); i <= svh; i++ {
		node := newFakeNode(0, i, currentNode)

		// Don't set stake versions.

		currentNode = node
		bc.bestNode = currentNode
		height = i
	}
	if height != svh {
		t.Fatalf("invalid height got %v expected %v", height,
			params.StakeValidationHeight)
	}

	// Generate 3 intervals with v2 votes and calculate StakeVersion.
	runCount := interval * 3
	for i := int64(0); i < runCount; i++ {
		node := newFakeNode(3, height+i, currentNode)

		// Set stake versions.
		for x := uint16(0); x < params.TicketsPerBlock; x++ {
			node.votes = append(node.votes,
				VoteVersionTuple{Version: 2})
		}

		sv, err := bc.calcStakeVersionByNode(currentNode)
		if err != nil {
			t.Fatalf("calcStakeVersionByNode: unexpected error: %v", err)
		}
		node.header.StakeVersion = sv

		currentNode = node
		bc.bestNode = currentNode

	}
	height += runCount

	if !bc.isStakeMajorityVersion(0, currentNode) {
		t.Fatalf("invalid StakeVersion expected 0 -> true")
	}
	if !bc.isStakeMajorityVersion(2, currentNode) {
		t.Fatalf("invalid StakeVersion expected 2 -> true")
	}
	if bc.isStakeMajorityVersion(4, currentNode) {
		t.Fatalf("invalid StakeVersion expected 4 -> false")
	}

	// Generate 3 intervals with v4 votes and calculate StakeVersion.
	runCount = interval * 3
	for i := int64(0); i < runCount; i++ {
		node := newFakeNode(3, height+i, currentNode)

		// Set stake versions.
		for x := uint16(0); x < params.TicketsPerBlock; x++ {
			node.votes = append(node.votes,
				VoteVersionTuple{Version: 4})
		}

		sv, err := bc.calcStakeVersionByNode(currentNode)
		if err != nil {
			t.Fatalf("calcStakeVersionByNode: unexpected error: %v", err)
		}
		node.header.StakeVersion = sv

		currentNode = node
		bc.bestNode = currentNode
	}
	height += runCount

	if !bc.isStakeMajorityVersion(0, currentNode) {
		t.Fatalf("invalid StakeVersion expected 0 -> true")
	}
	if !bc.isStakeMajorityVersion(2, currentNode) {
		t.Fatalf("invalid StakeVersion expected 2 -> true")
	}
	if !bc.isStakeMajorityVersion(4, currentNode) {
		t.Fatalf("invalid StakeVersion expected 4 -> true")
	}
	if bc.isStakeMajorityVersion(5, currentNode) {
		t.Fatalf("invalid StakeVersion expected 5 -> false")
	}

	// Generate 3 intervals with v2 votes and calculate StakeVersion.
	runCount = interval * 3
	for i := int64(0); i < runCount; i++ {
		node := newFakeNode(3, height+i, currentNode)

		// Set stake versions.
		for x := uint16(0); x < params.TicketsPerBlock; x++ {
			node.votes = append(node.votes,
				VoteVersionTuple{Version: 2})
		}

		sv, err := bc.calcStakeVersionByNode(currentNode)
		if err != nil {
			t.Fatalf("calcStakeVersionByNode: unexpected error: %v", err)
		}
		node.header.StakeVersion = sv

		currentNode = node
		bc.bestNode = currentNode

	}
	height += runCount

	if !bc.isStakeMajorityVersion(0, currentNode) {
		t.Fatalf("invalid StakeVersion expected 0 -> true")
	}
	if !bc.isStakeMajorityVersion(2, currentNode) {
		t.Fatalf("invalid StakeVersion expected 2 -> true")
	}
	if !bc.isStakeMajorityVersion(4, currentNode) {
		t.Fatalf("invalid StakeVersion expected 4 -> true")
	}
	if bc.isStakeMajorityVersion(5, currentNode) {
		t.Fatalf("invalid StakeVersion expected 5 -> false")
	}

	// Generate 2 interval with v5 votes
	runCount = interval * 2
	for i := int64(0); i < runCount; i++ {
		node := newFakeNode(3, height+i, currentNode)

		// Set stake versions.
		for x := uint16(0); x < params.TicketsPerBlock; x++ {
			node.votes = append(node.votes,
				VoteVersionTuple{Version: 5})
		}

		sv, err := bc.calcStakeVersionByNode(currentNode)
		if err != nil {
			t.Fatalf("calcStakeVersionByNode: unexpected error: %v", err)
		}
		node.header.StakeVersion = sv

		currentNode = node
		bc.bestNode = currentNode

	}
	height += runCount

	if !bc.isStakeMajorityVersion(0, currentNode) {
		t.Fatalf("invalid StakeVersion expected 0 -> true")
	}
	if !bc.isStakeMajorityVersion(2, currentNode) {
		t.Fatalf("invalid StakeVersion expected 2 -> true")
	}
	if !bc.isStakeMajorityVersion(4, currentNode) {
		t.Fatalf("invalid StakeVersion expected 4 -> true")
	}
	if !bc.isStakeMajorityVersion(5, currentNode) {
		t.Fatalf("invalid StakeVersion expected 5 -> true")
	}
	if bc.isStakeMajorityVersion(6, currentNode) {
		t.Fatalf("invalid StakeVersion expected 6 -> false")
	}

	// Generate 1 interval with v4 votes, to test the edge condition
	runCount = interval
	for i := int64(0); i < runCount; i++ {
		node := newFakeNode(3, height+i, currentNode)

		// Set stake versions.
		for x := uint16(0); x < params.TicketsPerBlock; x++ {
			node.votes = append(node.votes,
				VoteVersionTuple{Version: 4})
		}

		sv, err := bc.calcStakeVersionByNode(currentNode)
		if err != nil {
			t.Fatalf("calcStakeVersionByNode: unexpected error: %v", err)
		}
		node.header.StakeVersion = sv

		currentNode = node
		bc.bestNode = currentNode

	}
	height += runCount

	if !bc.isStakeMajorityVersion(0, currentNode) {
		t.Fatalf("invalid StakeVersion expected 0 -> true")
	}
	if !bc.isStakeMajorityVersion(2, currentNode) {
		t.Fatalf("invalid StakeVersion expected 2 -> true")
	}
	if !bc.isStakeMajorityVersion(4, currentNode) {
		t.Fatalf("invalid StakeVersion expected 4 -> true")
	}
	if !bc.isStakeMajorityVersion(5, currentNode) {
		t.Fatalf("invalid StakeVersion expected 5 -> true")
	}
	if bc.isStakeMajorityVersion(6, currentNode) {
		t.Fatalf("invalid StakeVersion expected 6 -> false")
	}

	// Generate 1 interval with v4 votes.
	runCount = interval
	for i := int64(0); i < runCount; i++ {
		node := newFakeNode(3, height+i, currentNode)

		// Set stake versions.
		for x := uint16(0); x < params.TicketsPerBlock; x++ {
			node.votes = append(node.votes,
				VoteVersionTuple{Version: 4})
		}

		sv, err := bc.calcStakeVersionByNode(currentNode)
		if err != nil {
			t.Fatalf("calcStakeVersionByNode: unexpected error: %v", err)
		}
		node.header.StakeVersion = sv

		currentNode = node
		bc.bestNode = currentNode

	}

	if !bc.isStakeMajorityVersion(0, currentNode) {
		t.Fatalf("invalid StakeVersion expected 0 -> true")
	}
	if !bc.isStakeMajorityVersion(2, currentNode) {
		t.Fatalf("invalid StakeVersion expected 2 -> true")
	}
	if !bc.isStakeMajorityVersion(4, currentNode) {
		t.Fatalf("invalid StakeVersion expected 4 -> true")
	}
	if !bc.isStakeMajorityVersion(5, currentNode) {
		t.Fatalf("invalid StakeVersion expected 5 -> true")
	}
	if bc.isStakeMajorityVersion(6, currentNode) {
		t.Fatalf("invalid StakeVersion expected 6 -> false")
	}
}

// TestCalcStakeVersionByNode doesn't work yet
func DNWTestCalcStakeVersionByNode(t *testing.T) {
	params := &chaincfg.SimNetParams

	tests := []struct {
		name          string
		numNodes      int64
		expectVersion uint32
		set           func(*blockNode)
	}{
		{
			name:          "headerStake 2 votes 3",
			numNodes:      params.StakeValidationHeight + params.StakeVersionInterval*3,
			expectVersion: 3,
			set: func(b *blockNode) {
				if int64(b.header.Height) > params.StakeValidationHeight {
					// set voter versions
					for x := 0; x < int(params.TicketsPerBlock); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 3})
					}

					// set header stake version
					b.header.StakeVersion = 2
					// set enforcement version
					b.header.Version = 3
				}
			},
		},
		{
			name:          "headerStake 3 votes 2",
			numNodes:      params.StakeValidationHeight + params.StakeVersionInterval*3,
			expectVersion: 3,
			set: func(b *blockNode) {
				if int64(b.header.Height) > params.StakeValidationHeight {
					// set voter versions
					for x := 0; x < int(params.TicketsPerBlock); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 2})
					}

					// set header stake version
					b.header.StakeVersion = 3
					// set enforcement version
					b.header.Version = 3
				}
			},
		},
	}

	for _, test := range tests {
		bc := newFakeChain(params)
		currentNode := genesisBlockNode(params)

		t.Logf("running: \"%v\"\n", test.name)
		for i := int64(1); i <= test.numNodes; i++ {
			/*
				// Make up a header.
				header := &wire.BlockHeader{
					Version: 1,
					Height:  uint32(i),
					Nonce:   uint32(0),
				}
			*/

			// add by sammy at 2017-10-25
			node := newBlockNode(FakeBlock(1, i, 0), nil, nil, nil)

			//node := newBlockNode(header, nil, nil, nil)
			node.height = i
			node.parent = currentNode

			test.set(node)

			currentNode = node
			bc.bestNode = currentNode
		}

		version, err := bc.calcStakeVersionByNode(bc.bestNode)
		t.Logf("name \"%v\" version %v err %v", test.name, version, err)
		if err != nil {
			t.Fatalf("calcStakeVersionByNode: unexpected error: %v", err)
		}
		if version != test.expectVersion {
			t.Fatalf("version mismatch: got %v expected %v",
				version, test.expectVersion)
		}
	}
}

// TestIsStakeMajorityVersion doesn't work yet
func DNWTestIsStakeMajorityVersion(t *testing.T) {
	params := &chaincfg.MainNetParams

	// Calculate super majority for 5 and 3 ticket maxes.
	maxTickets5 := int32(params.StakeVersionInterval) * int32(params.TicketsPerBlock)
	sm5 := maxTickets5 * params.StakeMajorityMultiplier / params.StakeMajorityDivisor
	maxTickets3 := int32(params.StakeVersionInterval) * int32(params.TicketsPerBlock-2)
	sm3 := maxTickets3 * params.StakeMajorityMultiplier / params.StakeMajorityDivisor

	// Keep track of ticketcount in set.  Must be reset every test.
	ticketCount := int32(0)

	tests := []struct {
		name                 string
		numNodes             int64
		set                  func(*blockNode)
		blockVersion         int32
		startStakeVersion    uint32
		expectedStakeVersion uint32
		expectedCalcVersion  uint32
		result               bool
	}{
		{
			name:                 "too shallow",
			numNodes:             params.StakeValidationHeight + params.StakeVersionInterval - 1,
			startStakeVersion:    1,
			expectedStakeVersion: 1,
			expectedCalcVersion:  0,
			result:               true,
		},
		{
			name:                 "just enough",
			numNodes:             params.StakeValidationHeight + params.StakeVersionInterval,
			startStakeVersion:    1,
			expectedStakeVersion: 1,
			expectedCalcVersion:  0,
			result:               true,
		},
		{
			name:                 "odd",
			numNodes:             params.StakeValidationHeight + params.StakeVersionInterval + 1,
			startStakeVersion:    1,
			expectedStakeVersion: 1,
			expectedCalcVersion:  0,
			result:               true,
		},
		{
			name:     "100%",
			numNodes: params.StakeValidationHeight + params.StakeVersionInterval,
			set: func(b *blockNode) {
				if int64(b.header.Height) > params.StakeValidationHeight {
					for x := 0; x < int(params.TicketsPerBlock); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 2})
					}
				}
			},
			startStakeVersion:    1,
			expectedStakeVersion: 2,
			expectedCalcVersion:  0,
			result:               true,
		},
		{
			name:     "50%",
			numNodes: params.StakeValidationHeight + (params.StakeVersionInterval * 2),
			set: func(b *blockNode) {
				if int64(b.header.Height) <= params.StakeValidationHeight {
					return
				}

				if int64(b.header.Height) < params.StakeValidationHeight+params.StakeVersionInterval {
					for x := 0; x < int(params.TicketsPerBlock); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 1})
					}
					return
				}

				threshold := maxTickets5 / 2

				v := uint32(1)
				for x := 0; x < int(params.TicketsPerBlock); x++ {
					if ticketCount >= threshold {
						v = 2
					}
					b.votes = append(b.votes,
						VoteVersionTuple{Version: v})
					ticketCount++
				}
			},
			startStakeVersion:    1,
			expectedStakeVersion: 2,
			expectedCalcVersion:  0,
			result:               false,
		},
		{
			name:     "75%-1",
			numNodes: params.StakeValidationHeight + (params.StakeVersionInterval * 2),
			set: func(b *blockNode) {
				if int64(b.header.Height) < params.StakeValidationHeight {
					return
				}

				if int64(b.header.Height) < params.StakeValidationHeight+params.StakeVersionInterval {
					for x := 0; x < int(params.TicketsPerBlock); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 1})
					}
					return
				}

				threshold := maxTickets5 - sm5 + 1

				v := uint32(1)
				for x := 0; x < int(params.TicketsPerBlock); x++ {
					if ticketCount >= threshold {
						v = 2
					}
					b.votes = append(b.votes,
						VoteVersionTuple{Version: v})
					ticketCount++
				}
			},
			startStakeVersion:    1,
			expectedStakeVersion: 2,
			expectedCalcVersion:  0,
			result:               false,
		},
		{
			name:     "75%",
			numNodes: params.StakeValidationHeight + (params.StakeVersionInterval * 2),
			set: func(b *blockNode) {
				if int64(b.header.Height) <= params.StakeValidationHeight {
					return
				}

				if int64(b.header.Height) < params.StakeValidationHeight+params.StakeVersionInterval {
					for x := 0; x < int(params.TicketsPerBlock); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 1})
					}
					return
				}

				threshold := maxTickets5 - sm5

				v := uint32(1)
				for x := 0; x < int(params.TicketsPerBlock); x++ {
					if ticketCount >= threshold {
						v = 2
					}
					b.votes = append(b.votes,
						VoteVersionTuple{Version: v})
					ticketCount++
				}
			},
			startStakeVersion:    1,
			expectedStakeVersion: 2,
			expectedCalcVersion:  0,
			result:               true,
		},
		{
			name:     "100% after several non majority intervals",
			numNodes: params.StakeValidationHeight + (params.StakeVersionInterval * 222),
			set: func(b *blockNode) {
				if int64(b.header.Height) <= params.StakeValidationHeight {
					return
				}

				if int64(b.header.Height) < params.StakeValidationHeight+params.StakeVersionInterval {
					for x := 0; x < int(params.TicketsPerBlock); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 1})
					}
					return
				}

				for x := 0; x < int(params.TicketsPerBlock); x++ {
					b.votes = append(b.votes,
						VoteVersionTuple{Version: uint32(x) % 5})
				}
			},
			startStakeVersion:    1,
			expectedStakeVersion: 1,
			expectedCalcVersion:  0,
			result:               true,
		},
		{
			name:     "no majority ever",
			numNodes: params.StakeValidationHeight + (params.StakeVersionInterval * 8),
			set: func(b *blockNode) {
				if int64(b.header.Height) <= params.StakeValidationHeight {
					return
				}

				for x := 0; x < int(params.TicketsPerBlock); x++ {
					b.votes = append(b.votes,
						VoteVersionTuple{Version: uint32(x) % 5})
				}
			},
			startStakeVersion:    1,
			expectedStakeVersion: 1,
			expectedCalcVersion:  0,
			result:               true,
		},
		{
			name:     "75%-1 with 3 votes",
			numNodes: params.StakeValidationHeight + (params.StakeVersionInterval * 2),
			set: func(b *blockNode) {
				if int64(b.header.Height) < params.StakeValidationHeight {
					return
				}

				if int64(b.header.Height) < params.StakeValidationHeight+params.StakeVersionInterval {
					for x := 0; x < int(params.TicketsPerBlock-2); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 1})
					}
					return
				}

				threshold := maxTickets3 - sm3 + 1

				v := uint32(1)
				for x := 0; x < int(params.TicketsPerBlock-2); x++ {
					if ticketCount >= threshold {
						v = 2
					}
					b.votes = append(b.votes,
						VoteVersionTuple{Version: v})
					ticketCount++
				}
			},
			startStakeVersion:    1,
			expectedStakeVersion: 2,
			expectedCalcVersion:  0,
			result:               false,
		},
		{
			name:     "75% with 3 votes",
			numNodes: params.StakeValidationHeight + (params.StakeVersionInterval * 2),
			set: func(b *blockNode) {
				if int64(b.header.Height) <= params.StakeValidationHeight {
					return
				}

				if int64(b.header.Height) < params.StakeValidationHeight+params.StakeVersionInterval {
					for x := 0; x < int(params.TicketsPerBlock-2); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 1})
					}
					return
				}

				threshold := maxTickets3 - sm3

				v := uint32(1)
				for x := 0; x < int(params.TicketsPerBlock-2); x++ {
					if ticketCount >= threshold {
						v = 2
					}
					b.votes = append(b.votes,
						VoteVersionTuple{Version: v})
					ticketCount++
				}
			},
			startStakeVersion:    1,
			expectedStakeVersion: 2,
			expectedCalcVersion:  0,
			result:               true,
		},
		{
			name:     "75% with 3 votes blockversion 3",
			numNodes: params.StakeValidationHeight + (params.StakeVersionInterval * 2),
			set: func(b *blockNode) {
				if int64(b.header.Height) <= params.StakeValidationHeight {
					return
				}

				if int64(b.header.Height) < params.StakeValidationHeight+params.StakeVersionInterval {
					for x := 0; x < int(params.TicketsPerBlock-2); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 1})
					}
					return
				}

				threshold := maxTickets3 - sm3

				v := uint32(1)
				for x := 0; x < int(params.TicketsPerBlock-2); x++ {
					if ticketCount >= threshold {
						v = 2
					}
					b.votes = append(b.votes,
						VoteVersionTuple{Version: v})
					ticketCount++
				}
			},
			blockVersion:         3,
			startStakeVersion:    1,
			expectedStakeVersion: 2,
			expectedCalcVersion:  2,
			result:               true,
		},
		{
			name:     "75%-1 with 3 votes blockversion 3",
			numNodes: params.StakeValidationHeight + (params.StakeVersionInterval * 2),
			set: func(b *blockNode) {
				if int64(b.header.Height) < params.StakeValidationHeight {
					return
				}

				if int64(b.header.Height) < params.StakeValidationHeight+params.StakeVersionInterval {
					for x := 0; x < int(params.TicketsPerBlock-2); x++ {
						b.votes = append(b.votes,
							VoteVersionTuple{Version: 1})
					}
					return
				}

				threshold := maxTickets3 - sm3 + 1

				v := uint32(1)
				for x := 0; x < int(params.TicketsPerBlock-2); x++ {
					if ticketCount >= threshold {
						v = 2
					}
					b.votes = append(b.votes,
						VoteVersionTuple{Version: v})
					ticketCount++
				}
			},
			blockVersion:         3,
			startStakeVersion:    1,
			expectedStakeVersion: 2,
			expectedCalcVersion:  1,
			result:               false,
		},
	}

	for _, test := range tests {
		// Create new BlockChain in order to blow away cache.
		bc := newFakeChain(params)

		ticketCount = 0

		genesisNode := genesisBlockNode(params)
		genesisNode.header.StakeVersion = test.startStakeVersion

		t.Logf("running: %v\n", test.name)
		var currentNode *blockNode
		currentNode = genesisNode
		for i := int64(1); i <= test.numNodes; i++ {
			/*
				// Make up a header.
				header := &wire.BlockHeader{
					Version:      test.blockVersion,
					Height:       uint32(i),
					Nonce:        uint32(0),
					StakeVersion: test.startStakeVersion,
				}
				node := newBlockNode(header, nil, nil, nil)
			*/

			node := newBlockNode(FakeBlock(test.blockVersion, i, test.startStakeVersion), nil, nil, nil)
			node.height = i
			node.parent = currentNode

			// Override version.
			if test.set != nil {
				test.set(node)
			} else {
				for x := 0; x < int(params.TicketsPerBlock); x++ {
					node.votes = append(node.votes,
						VoteVersionTuple{Version: test.startStakeVersion})
				}
			}

			currentNode = node
			bc.bestNode = currentNode
		}

		res := bc.isVoterMajorityVersion(test.expectedStakeVersion, currentNode)
		if res != test.result {
			t.Fatalf("%v isVoterMajorityVersion", test.name)
		}

		// validate calcStakeVersion
		version, err := bc.calcStakeVersionByNode(currentNode)
		if err != nil {
			t.Fatalf("calcStakeVersionByNode: unexpected error: %v", err)
		}
		if version != test.expectedCalcVersion {
			t.Fatalf("%v calcStakeVersionByNode got %v expected %v",
				test.name, version, test.expectedCalcVersion)
		}
	}
}

// TestLarge doestn't work yet
func DNWTestLarge(t *testing.T) {
	params := &chaincfg.MainNetParams

	numRuns := 5
	numBlocks := params.StakeVersionInterval * 100
	numBlocksShallow := params.StakeVersionInterval * 10
	tests := []struct {
		name                 string
		numNodes             int64
		set                  func(*blockNode)
		blockVersion         int32
		startStakeVersion    uint32
		expectedStakeVersion uint32
		expectedCalcVersion  uint32
		result               bool
	}{
		{
			name:                 "shallow cache",
			numNodes:             numBlocksShallow,
			startStakeVersion:    1,
			expectedStakeVersion: 1,
			expectedCalcVersion:  0,
			result:               true,
		},
		{
			name:                 "deep cache",
			numNodes:             numBlocks,
			startStakeVersion:    1,
			expectedStakeVersion: 1,
			expectedCalcVersion:  0,
			result:               true,
		},
	}

	for _, test := range tests {
		// Create new BlockChain in order to blow away cache.
		bc := newFakeChain(params)

		genesisNode := genesisBlockNode(params)
		genesisNode.header.StakeVersion = test.startStakeVersion

		t.Logf("running: %v with %v nodes\n", test.name, test.numNodes)
		var currentNode *blockNode
		currentNode = genesisNode
		for i := int64(1); i <= test.numNodes; i++ {
			/*
				// Make up a header.
				header := &wire.BlockHeader{
					Version:      test.blockVersion,
					Height:       uint32(i),
					Nonce:        uint32(0),
					StakeVersion: test.startStakeVersion,
				}
				node := newBlockNode(header, nil, nil, nil)
			*/

			node := newBlockNode(FakeBlock(test.blockVersion, i, test.startStakeVersion), nil, nil, nil)
			node.height = i
			node.parent = currentNode

			// Override version.
			for x := 0; x < int(params.TicketsPerBlock); x++ {
				node.votes = append(node.votes,
					VoteVersionTuple{Version: test.startStakeVersion})
			}

			currentNode = node
			bc.bestNode = currentNode
		}

		for x := 0; x < numRuns; x++ {
			start := time.Now()
			res := bc.isVoterMajorityVersion(test.expectedStakeVersion, currentNode)
			if res != test.result {
				t.Fatalf("%v isVoterMajorityVersion got %v expected %v", test.name, res, test.result)
			}

			// validate calcStakeVersion
			version, err := bc.calcStakeVersionByNode(currentNode)
			if err != nil {
				t.Fatalf("calcStakeVersionByNode: unexpected error: %v", err)
			}
			if version != test.expectedCalcVersion {
				t.Fatalf("%v calcStakeVersionByNode got %v expected %v",
					test.name, version, test.expectedCalcVersion)
			}
			end := time.Now()

			setup := "setup 0"
			if x != 0 {
				setup = fmt.Sprintf("run %v", x)
			}

			vkey := stakeMajorityCacheKeySize + 8 // bool on x86_64
			key := chainhash.HashSize + 4         // size of uint32

			cost := len(bc.isVoterMajorityVersionCache) * vkey
			cost += len(bc.isStakeMajorityVersionCache) * vkey
			cost += len(bc.calcPriorStakeVersionCache) * key
			cost += len(bc.calcVoterVersionIntervalCache) * key
			cost += len(bc.calcStakeVersionCache) * key
			memoryCost := fmt.Sprintf("memory cost: %v", cost)

			t.Logf("run time (%v) %v %v", setup, end.Sub(start),
				memoryCost)
		}
	}
}
