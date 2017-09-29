// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

// BlockOneLedgerMainNet is the block one output ledger for the main
// network.
var BlockOneLedgerMainNet = []*TokenPayout{
	{"HsTRMWkvC2GGkg6k1kTjrV1tyYPGN6tTAGS", 50000 * 1e8},  //wancc
	{"HsGEnUynvBzaMcgydSiRmEKKb1nZTHdZVoC", 50000 * 1e8}, //fanlq
	{"HsDnJ7hMNXjYrfJfUNbwwuDBhWWNfuYoeuR", 50000 * 1e8}, //shanyl
	{"HsU7GADmBspFq5bdKT3Lkz9rg1x7VB6TQPt", 50000 * 1e8}, //panc
}

// BlockOneLedgerTestNet is the block one output ledger for the test
// network.
var BlockOneLedgerTestNet = []*TokenPayout{
}

// BlockOneLedgerTestNet2 is the block one output ledger for the 2nd test
// network.
var BlockOneLedgerTestNet2 = []*TokenPayout{
}

// BlockOneLedgerSimNet is the block one output ledger for the simulation
// network. See under "Hypercash organization related parameters" in params.go
// for information on how to spend these outputs.
var BlockOneLedgerSimNet = []*TokenPayout{
}
