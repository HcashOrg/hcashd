// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

// BlockOneLedgerMainNet is the block one output ledger for the main
// network.
var BlockOneLedgerMainNet = []*TokenPayout{
	{"HsMNycPD277U4Zw2qNHNqY8r3MSsnhhiiGW", 50000 * 1e8}, //wancc
	{"HsN4DLc5n7kKyfUMKm5SW56J8LJGTK6u91h", 50000 * 1e8}, //fanlq
	{"HsF4tYLz9JpUFk9aPLC7U2AN8Deq6LkyoWc", 50000 * 1e8}, //shanyl
	{"HsZKbCUvcpjfHAJpDfWikD7E2oXUGR4ge6q", 50000 * 1e8}, //panc
	{"HsNu7JN9SeNb3cH7BJWMiSSpPqY6rz8BSXW", 50000 * 1e8}, //guxy
	{"HsZKbCUvcpjfHAJpDfWikD7E2oXUGR4ge6q", 50000 * 1e8}, //dengcg
	{"HsZKbCUvcpjfHAJpDfWikD7E2oXUGR4ge6q", 50000 * 1e8}, //lixm
	{"HsS6Hqt7yB5Fr2HDBz5gRhK75q7ciuxa7au", 50000 * 1e8}, //yaoyq
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
