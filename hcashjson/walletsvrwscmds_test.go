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

// TestWalletSvrWsCmds tests all of the wallet server websocket-specific
// commands marshal and unmarshal into valid results include handling of
// optional fields being omitted in the marshalled command, while optional
// fields with defaults have the default assigned on unmarshalled commands.
func DNWTestWalletSvrWsCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "createencryptedwallet",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("createencryptedwallet", "pass")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewCreateEncryptedWalletCmd("pass")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"createencryptedwallet","params":["pass"],"id":1}`,
			unmarshalled: &hcashjson.CreateEncryptedWalletCmd{Passphrase: "pass"},
		},
		{
			name: "exportwatchingwallet",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("exportwatchingwallet")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewExportWatchingWalletCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"exportwatchingwallet","params":[],"id":1}`,
			unmarshalled: &hcashjson.ExportWatchingWalletCmd{
				Account:  nil,
				Download: hcashjson.Bool(false),
			},
		},
		{
			name: "exportwatchingwallet optional1",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("exportwatchingwallet", "acct")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewExportWatchingWalletCmd(hcashjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"exportwatchingwallet","params":["acct"],"id":1}`,
			unmarshalled: &hcashjson.ExportWatchingWalletCmd{
				Account:  hcashjson.String("acct"),
				Download: hcashjson.Bool(false),
			},
		},
		{
			name: "exportwatchingwallet optional2",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("exportwatchingwallet", "acct", true)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewExportWatchingWalletCmd(hcashjson.String("acct"),
					hcashjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"exportwatchingwallet","params":["acct",true],"id":1}`,
			unmarshalled: &hcashjson.ExportWatchingWalletCmd{
				Account:  hcashjson.String("acct"),
				Download: hcashjson.Bool(true),
			},
		},
		{
			name: "getunconfirmedbalance",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getunconfirmedbalance")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetUnconfirmedBalanceCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getunconfirmedbalance","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetUnconfirmedBalanceCmd{
				Account: nil,
			},
		},
		{
			name: "getunconfirmedbalance optional1",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getunconfirmedbalance", "acct")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetUnconfirmedBalanceCmd(hcashjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getunconfirmedbalance","params":["acct"],"id":1}`,
			unmarshalled: &hcashjson.GetUnconfirmedBalanceCmd{
				Account: hcashjson.String("acct"),
			},
		},
		{
			name: "listaddresstransactions",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("listaddresstransactions", `["1Address"]`)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewListAddressTransactionsCmd([]string{"1Address"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaddresstransactions","params":[["1Address"]],"id":1}`,
			unmarshalled: &hcashjson.ListAddressTransactionsCmd{
				Addresses: []string{"1Address"},
				Account:   nil,
			},
		},
		{
			name: "listaddresstransactions optional1",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("listaddresstransactions", `["1Address"]`, "acct")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewListAddressTransactionsCmd([]string{"1Address"},
					hcashjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaddresstransactions","params":[["1Address"],"acct"],"id":1}`,
			unmarshalled: &hcashjson.ListAddressTransactionsCmd{
				Addresses: []string{"1Address"},
				Account:   hcashjson.String("acct"),
			},
		},
		{
			name: "listalltransactions",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("listalltransactions")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewListAllTransactionsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listalltransactions","params":[],"id":1}`,
			unmarshalled: &hcashjson.ListAllTransactionsCmd{
				Account: nil,
			},
		},
		{
			name: "listalltransactions optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("listalltransactions", "acct")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewListAllTransactionsCmd(hcashjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listalltransactions","params":["acct"],"id":1}`,
			unmarshalled: &hcashjson.ListAllTransactionsCmd{
				Account: hcashjson.String("acct"),
			},
		},
		{
			name: "recoveraddresses",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("recoveraddresses", "acct", 10)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewRecoverAddressesCmd("acct", 10)
			},
			marshalled: `{"jsonrpc":"1.0","method":"recoveraddresses","params":["acct",10],"id":1}`,
			unmarshalled: &hcashjson.RecoverAddressesCmd{
				Account: "acct",
				N:       10,
			},
		},
		{
			name: "walletislocked",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("walletislocked")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewWalletIsLockedCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletislocked","params":[],"id":1}`,
			unmarshalled: &hcashjson.WalletIsLockedCmd{},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := hcashjson.MarshalCmd(testID, test.staticCmd())
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

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = hcashjson.MarshalCmd(testID, cmd)
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
