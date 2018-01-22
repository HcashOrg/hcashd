// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hcashjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/HcashOrg/hcashd/hcashjson"
)

// TestChainSvrWsNtfns tests all of the chain server websocket-specific
// notifications marshal and unmarshal into valid results include handling of
// optional fields being omitted in the marshalled command, while optional
// fields with defaults have the default assigned on unmarshalled commands.
func TestChainSvrWsNtfns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		newNtfn      func() (interface{}, error)
		staticNtfn   func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "blockconnected",
			newNtfn: func() (interface{}, error) {
				return hcashjson.NewCmd("blockconnected", "header", []string{"tx0", "tx1"})
			},
			staticNtfn: func() interface{} {
				return hcashjson.NewBlockConnectedNtfn("header", []string{"tx0", "tx1"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"blockconnected","params":["header",["tx0","tx1"]],"id":null}`,
			unmarshalled: &hcashjson.BlockConnectedNtfn{
				Header:        "header",
				SubscribedTxs: []string{"tx0", "tx1"},
			},
		},
		{
			name: "blockdisconnected",
			newNtfn: func() (interface{}, error) {
				return hcashjson.NewCmd("blockdisconnected", "header")
			},
			staticNtfn: func() interface{} {
				return hcashjson.NewBlockDisconnectedNtfn("header")
			},
			marshalled: `{"jsonrpc":"1.0","method":"blockdisconnected","params":["header"],"id":null}`,
			unmarshalled: &hcashjson.BlockDisconnectedNtfn{
				Header: "header",
			},
		},
		{
			name: "relevanttxaccepted",
			newNtfn: func() (interface{}, error) {
				return hcashjson.NewCmd("relevanttxaccepted", "001122")
			},
			staticNtfn: func() interface{} {
				return hcashjson.NewRelevantTxAcceptedNtfn("001122")
			},
			marshalled: `{"jsonrpc":"1.0","method":"relevanttxaccepted","params":["001122"],"id":null}`,
			unmarshalled: &hcashjson.RelevantTxAcceptedNtfn{
				Transaction: "001122",
			},
		},
		{
			name: "txaccepted",
			newNtfn: func() (interface{}, error) {
				return hcashjson.NewCmd("txaccepted", "123", 1.5)
			},
			staticNtfn: func() interface{} {
				return hcashjson.NewTxAcceptedNtfn("123", 1.5)
			},
			marshalled: `{"jsonrpc":"1.0","method":"txaccepted","params":["123",1.5],"id":null}`,
			unmarshalled: &hcashjson.TxAcceptedNtfn{
				TxID:   "123",
				Amount: 1.5,
			},
		},
		{
			name: "txacceptedverbose",
			newNtfn: func() (interface{}, error) {
				return hcashjson.NewCmd("txacceptedverbose", `{"hex":"001122","txid":"123","version":1,"locktime":4294967295,"vin":null,"vout":null,"confirmations":0,"txtype":"regular","size":256}`)
			},
			staticNtfn: func() interface{} {
				txResult := hcashjson.TxRawResult{
					Hex:           "001122",
					Txid:          "123",
					Version:       1,
					LockTime:      4294967295,
					Vin:           nil,
					Vout:          nil,
					TxType:        "regular",
					Size:          256,
					Confirmations: 0,
				}
				return hcashjson.NewTxAcceptedVerboseNtfn(txResult)
			},
			marshalled: `{"jsonrpc":"1.0","method":"txacceptedverbose","params":[{"hex":"001122","txid":"123","version":1,"locktime":4294967295,"expiry":0,"vin":null,"vout":null,"blockheight":0,"blockkeyheight":0,"txtype":"regular","size":256}],"id":null}`,
			unmarshalled: &hcashjson.TxAcceptedVerboseNtfn{
				RawTx: hcashjson.TxRawResult{
					Hex:           "001122",
					Txid:          "123",
					Version:       1,
					LockTime:      4294967295,
					Vin:           nil,
					Vout:          nil,
					TxType:        "regular",
					Size:          256,
					Confirmations: 0,
				},
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the notification as created by the new static
		// creation function.  The ID is nil for notifications.
		marshalled, err := hcashjson.MarshalCmd(nil, test.staticNtfn())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		// Ensure the notification is created without error via the
		// generic new notification creation function.
		cmd, err := test.newNtfn()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the notification as created by the generic new
		// notification creation function.    The ID is nil for
		// notifications.
		marshalled, err = hcashjson.MarshalCmd(nil, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request hcashjson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = hcashjson.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}
