// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2016 The Decred developers
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

// TestChainSvrCmds tests all of the chain server commands marshal and unmarshal
// into valid results include handling of optional fields being omitted in the
// marshalled command, while optional fields with defaults have the default
// assigned on unmarshalled commands.
func TestChainSvrCmds(t *testing.T) {
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
			name: "addnode",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("addnode", "127.0.0.1", hcashjson.ANRemove)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewAddNodeCmd("127.0.0.1", hcashjson.ANRemove)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"addnode","params":["127.0.0.1","remove"],"id":1}`,
			unmarshalled: &hcashjson.AddNodeCmd{Addr: "127.0.0.1", SubCmd: hcashjson.ANRemove},
		},
		{
			name: "createrawtransaction",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				txInputs := []hcashjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return hcashjson.NewCreateRawTransactionCmd(txInputs, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1,"tree":0}],{"456":0.0123}],"id":1}`,
			unmarshalled: &hcashjson.CreateRawTransactionCmd{
				Inputs:  []hcashjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1,"tree":0}]`,
					`{"456":0.0123}`, int64(12312333333))
			},
			staticCmd: func() interface{} {
				txInputs := []hcashjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return hcashjson.NewCreateRawTransactionCmd(txInputs, amounts, hcashjson.Int64(12312333333))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1,"tree":0}],{"456":0.0123},12312333333],"id":1}`,
			unmarshalled: &hcashjson.CreateRawTransactionCmd{
				Inputs:   []hcashjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts:  map[string]float64{"456": .0123},
				LockTime: hcashjson.Int64(12312333333),
			},
		},
		{
			name: "decoderawtransaction",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("decoderawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewDecodeRawTransactionCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decoderawtransaction","params":["123"],"id":1}`,
			unmarshalled: &hcashjson.DecodeRawTransactionCmd{HexTx: "123"},
		},
		{
			name: "decodescript",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("decodescript", "00")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewDecodeScriptCmd("00")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decodescript","params":["00"],"id":1}`,
			unmarshalled: &hcashjson.DecodeScriptCmd{HexScript: "00"},
		},
		{
			name: "getaddednodeinfo",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getaddednodeinfo", true)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetAddedNodeInfoCmd(true, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true],"id":1}`,
			unmarshalled: &hcashjson.GetAddedNodeInfoCmd{DNS: true, Node: nil},
		},
		{
			name: "getaddednodeinfo optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getaddednodeinfo", true, "127.0.0.1")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetAddedNodeInfoCmd(true, hcashjson.String("127.0.0.1"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true,"127.0.0.1"],"id":1}`,
			unmarshalled: &hcashjson.GetAddedNodeInfoCmd{
				DNS:  true,
				Node: hcashjson.String("127.0.0.1"),
			},
		},
		{
			name: "getbestblockhash",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getbestblockhash")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBestBlockHashCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblockhash","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetBestBlockHashCmd{},
		},
		{
			name: "getblock",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblock", "123")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBlockCmd("123", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123"],"id":1}`,
			unmarshalled: &hcashjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   hcashjson.Bool(true),
				VerboseTx: hcashjson.Bool(false),
			},
		},
		{
			name: "getblock required optional1",
			newCmd: func() (interface{}, error) {
				// Intentionally use a source param that is
				// more pointers than the destination to
				// exercise that path.
				verbosePtr := hcashjson.Bool(true)
				return hcashjson.NewCmd("getblock", "123", &verbosePtr)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBlockCmd("123", hcashjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true],"id":1}`,
			unmarshalled: &hcashjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   hcashjson.Bool(true),
				VerboseTx: hcashjson.Bool(false),
			},
		},
		{
			name: "getblock required optional2",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblock", "123", true, true)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBlockCmd("123", hcashjson.Bool(true), hcashjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true,true],"id":1}`,
			unmarshalled: &hcashjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   hcashjson.Bool(true),
				VerboseTx: hcashjson.Bool(true),
			},
		},
		{
			name: "getblockchaininfo",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblockchaininfo")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBlockChainInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockchaininfo","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetBlockChainInfoCmd{},
		},
		{
			name: "getblockcount",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblockcount")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBlockCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockcount","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetBlockCountCmd{},
		},
		{
			name: "getblockhash",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblockhash", 123)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBlockHashCmd(123)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockhash","params":[123],"id":1}`,
			unmarshalled: &hcashjson.GetBlockHashCmd{Index: 123},
		},
		{
			name: "getblockheader",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblockheader", "123")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBlockHeaderCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockheader","params":["123"],"id":1}`,
			unmarshalled: &hcashjson.GetBlockHeaderCmd{
				Hash:    "123",
				Verbose: hcashjson.Bool(true),
			},
		},
		{
			name: "getblocksubsidy",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblocksubsidy", 123, 256)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBlockSubsidyCmd(123, 256)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocksubsidy","params":[123,256],"id":1}`,
			unmarshalled: &hcashjson.GetBlockSubsidyCmd{
				Height: 123,
				Voters: 256,
			},
		},
		{
			name: "getblocktemplate",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblocktemplate")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetBlockTemplateCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblocktemplate","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetBlockTemplateCmd{Request: nil},
		},
		{
			name: "getblocktemplate optional - template request",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"]}`)
			},
			staticCmd: func() interface{} {
				template := hcashjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				}
				return hcashjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"]}],"id":1}`,
			unmarshalled: &hcashjson.GetBlockTemplateCmd{
				Request: &hcashjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := hcashjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   500,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return hcashjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &hcashjson.GetBlockTemplateCmd{
				Request: &hcashjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   int64(500),
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks 2",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := hcashjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return hcashjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &hcashjson.GetBlockTemplateCmd{
				Request: &hcashjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getchaintips",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getchaintips")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetChainTipsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintips","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetChainTipsCmd{},
		},
		{
			name: "getconnectioncount",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getconnectioncount")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetConnectionCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getconnectioncount","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetConnectionCountCmd{},
		},
		{
			name: "getdifficulty",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getdifficulty")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetDifficultyCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdifficulty","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetDifficultyCmd{},
		},
		{
			name: "getgenerate",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getgenerate")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetGenerateCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getgenerate","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetGenerateCmd{},
		},
		{
			name: "gethashespersec",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("gethashespersec")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetHashesPerSecCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gethashespersec","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetHashesPerSecCmd{},
		},
		{
			name: "getinfo",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getinfo")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetInfoCmd{},
		},
		{
			name: "getmempoolinfo",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getmempoolinfo")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetMempoolInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmempoolinfo","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetMempoolInfoCmd{},
		},
		{
			name: "getmininginfo",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getmininginfo")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetMiningInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmininginfo","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetMiningInfoCmd{},
		},
		{
			name: "getnetworkinfo",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getnetworkinfo")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetNetworkInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnetworkinfo","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetNetworkInfoCmd{},
		},
		{
			name: "getnettotals",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getnettotals")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetNetTotalsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnettotals","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetNetTotalsCmd{},
		},
		{
			name: "getnetworkhashps",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getnetworkhashps")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetNetworkHashPSCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetNetworkHashPSCmd{
				Blocks: hcashjson.Int(120),
				Height: hcashjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional1",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getnetworkhashps", 200)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetNetworkHashPSCmd(hcashjson.Int(200), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200],"id":1}`,
			unmarshalled: &hcashjson.GetNetworkHashPSCmd{
				Blocks: hcashjson.Int(200),
				Height: hcashjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional2",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getnetworkhashps", 200, 123)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetNetworkHashPSCmd(hcashjson.Int(200), hcashjson.Int(123))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200,123],"id":1}`,
			unmarshalled: &hcashjson.GetNetworkHashPSCmd{
				Blocks: hcashjson.Int(200),
				Height: hcashjson.Int(123),
			},
		},
		{
			name: "getpeerinfo",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getpeerinfo")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetPeerInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getpeerinfo","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetPeerInfoCmd{},
		},
		{
			name: "getrawmempool",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getrawmempool")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetRawMempoolCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetRawMempoolCmd{
				Verbose: hcashjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getrawmempool", false)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetRawMempoolCmd(hcashjson.Bool(false), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false],"id":1}`,
			unmarshalled: &hcashjson.GetRawMempoolCmd{
				Verbose: hcashjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional 2",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getrawmempool", false, "all")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetRawMempoolCmd(hcashjson.Bool(false), hcashjson.String("all"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false,"all"],"id":1}`,
			unmarshalled: &hcashjson.GetRawMempoolCmd{
				Verbose: hcashjson.Bool(false),
				TxType:  hcashjson.String("all"),
			},
		},
		{
			name: "getrawtransaction",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getrawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetRawTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123"],"id":1}`,
			unmarshalled: &hcashjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: hcashjson.Int(0),
			},
		},
		{
			name: "getrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getrawtransaction", "123", 1)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetRawTransactionCmd("123", hcashjson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123",1],"id":1}`,
			unmarshalled: &hcashjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: hcashjson.Int(1),
			},
		},
		{
			name: "gettxout",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("gettxout", "123", 1)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetTxOutCmd("123", 1, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1],"id":1}`,
			unmarshalled: &hcashjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: hcashjson.Bool(true),
			},
		},
		{
			name: "gettxout optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("gettxout", "123", 1, true)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetTxOutCmd("123", 1, hcashjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1,true],"id":1}`,
			unmarshalled: &hcashjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: hcashjson.Bool(true),
			},
		},
		{
			name: "gettxoutsetinfo",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("gettxoutsetinfo")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetTxOutSetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gettxoutsetinfo","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetTxOutSetInfoCmd{},
		},
		{
			name: "getwork",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getwork")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetWorkCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":[],"id":1}`,
			unmarshalled: &hcashjson.GetWorkCmd{
				Data: nil,
			},
		},
		{
			name: "getwork optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("getwork", "00112233")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewGetWorkCmd(hcashjson.String("00112233"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":["00112233"],"id":1}`,
			unmarshalled: &hcashjson.GetWorkCmd{
				Data: hcashjson.String("00112233"),
			},
		},
		{
			name: "help",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("help")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewHelpCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":[],"id":1}`,
			unmarshalled: &hcashjson.HelpCmd{
				Command: nil,
			},
		},
		{
			name: "help optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("help", "getblock")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewHelpCmd(hcashjson.String("getblock"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":["getblock"],"id":1}`,
			unmarshalled: &hcashjson.HelpCmd{
				Command: hcashjson.String("getblock"),
			},
		},
		{
			name: "ping",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("ping")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewPingCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"ping","params":[],"id":1}`,
			unmarshalled: &hcashjson.PingCmd{},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("searchrawtransactions", "1Address")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSearchRawTransactionsCmd("1Address", nil, nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address"],"id":1}`,
			unmarshalled: &hcashjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcashjson.Int(1),
				Skip:        hcashjson.Int(0),
				Count:       hcashjson.Int(100),
				VinExtra:    hcashjson.Int(0),
				Reverse:     hcashjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("searchrawtransactions", "1Address", 0)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSearchRawTransactionsCmd("1Address",
					hcashjson.Int(0), nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0],"id":1}`,
			unmarshalled: &hcashjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcashjson.Int(0),
				Skip:        hcashjson.Int(0),
				Count:       hcashjson.Int(100),
				VinExtra:    hcashjson.Int(0),
				Reverse:     hcashjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("searchrawtransactions", "1Address", 0, 5)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSearchRawTransactionsCmd("1Address",
					hcashjson.Int(0), hcashjson.Int(5), nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5],"id":1}`,
			unmarshalled: &hcashjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcashjson.Int(0),
				Skip:        hcashjson.Int(5),
				Count:       hcashjson.Int(100),
				VinExtra:    hcashjson.Int(0),
				Reverse:     hcashjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSearchRawTransactionsCmd("1Address",
					hcashjson.Int(0), hcashjson.Int(5), hcashjson.Int(10), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10],"id":1}`,
			unmarshalled: &hcashjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcashjson.Int(0),
				Skip:        hcashjson.Int(5),
				Count:       hcashjson.Int(10),
				VinExtra:    hcashjson.Int(0),
				Reverse:     hcashjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSearchRawTransactionsCmd("1Address",
					hcashjson.Int(0), hcashjson.Int(5), hcashjson.Int(10), hcashjson.Int(1), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1],"id":1}`,
			unmarshalled: &hcashjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcashjson.Int(0),
				Skip:        hcashjson.Int(5),
				Count:       hcashjson.Int(10),
				VinExtra:    hcashjson.Int(1),
				Reverse:     hcashjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSearchRawTransactionsCmd("1Address",
					hcashjson.Int(0), hcashjson.Int(5), hcashjson.Int(10),
					hcashjson.Int(1), hcashjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true],"id":1}`,
			unmarshalled: &hcashjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcashjson.Int(0),
				Skip:        hcashjson.Int(5),
				Count:       hcashjson.Int(10),
				VinExtra:    hcashjson.Int(1),
				Reverse:     hcashjson.Bool(true),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSearchRawTransactionsCmd("1Address",
					hcashjson.Int(0), hcashjson.Int(5), hcashjson.Int(10),
					hcashjson.Int(1), hcashjson.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true,["1Address"]],"id":1}`,
			unmarshalled: &hcashjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcashjson.Int(0),
				Skip:        hcashjson.Int(5),
				Count:       hcashjson.Int(10),
				VinExtra:    hcashjson.Int(1),
				Reverse:     hcashjson.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "sendrawtransaction",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("sendrawtransaction", "1122")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSendRawTransactionCmd("1122", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
			unmarshalled: &hcashjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: hcashjson.Bool(false),
			},
		},
		{
			name: "sendrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("sendrawtransaction", "1122", false)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSendRawTransactionCmd("1122", hcashjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &hcashjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: hcashjson.Bool(false),
			},
		},
		{
			name: "setgenerate",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("setgenerate", true)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSetGenerateCmd(true, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true],"id":1}`,
			unmarshalled: &hcashjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: hcashjson.Int(-1),
			},
		},
		{
			name: "setgenerate optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("setgenerate", true, 6)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSetGenerateCmd(true, hcashjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true,6],"id":1}`,
			unmarshalled: &hcashjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: hcashjson.Int(6),
			},
		},
		{
			name: "stop",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("stop")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewStopCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stop","params":[],"id":1}`,
			unmarshalled: &hcashjson.StopCmd{},
		},
		{
			name: "submitblock",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("submitblock", "112233")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewSubmitBlockCmd("112233", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233"],"id":1}`,
			unmarshalled: &hcashjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options:  nil,
			},
		},
		{
			name: "submitblock optional",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("submitblock", "112233", `{"workid":"12345"}`)
			},
			staticCmd: func() interface{} {
				options := hcashjson.SubmitBlockOptions{
					WorkID: "12345",
				}
				return hcashjson.NewSubmitBlockCmd("112233", &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233",{"workid":"12345"}],"id":1}`,
			unmarshalled: &hcashjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options: &hcashjson.SubmitBlockOptions{
					WorkID: "12345",
				},
			},
		},
		{
			name: "validateaddress",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("validateaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewValidateAddressCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"validateaddress","params":["1Address"],"id":1}`,
			unmarshalled: &hcashjson.ValidateAddressCmd{
				Address: "1Address",
			},
		},
		{
			name: "verifychain",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("verifychain")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewVerifyChainCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[],"id":1}`,
			unmarshalled: &hcashjson.VerifyChainCmd{
				CheckLevel: hcashjson.Int64(3),
				CheckDepth: hcashjson.Int64(288),
			},
		},
		{
			name: "verifychain optional1",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("verifychain", 2)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewVerifyChainCmd(hcashjson.Int64(2), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2],"id":1}`,
			unmarshalled: &hcashjson.VerifyChainCmd{
				CheckLevel: hcashjson.Int64(2),
				CheckDepth: hcashjson.Int64(288),
			},
		},
		{
			name: "verifychain optional2",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("verifychain", 2, 500)
			},
			staticCmd: func() interface{} {
				return hcashjson.NewVerifyChainCmd(hcashjson.Int64(2), hcashjson.Int64(500))
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2,500],"id":1}`,
			unmarshalled: &hcashjson.VerifyChainCmd{
				CheckLevel: hcashjson.Int64(2),
				CheckDepth: hcashjson.Int64(500),
			},
		},
		{
			name: "verifymessage",
			newCmd: func() (interface{}, error) {
				return hcashjson.NewCmd("verifymessage", "1Address", "301234", "test")
			},
			staticCmd: func() interface{} {
				return hcashjson.NewVerifyMessageCmd("1Address", "301234", "test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifymessage","params":["1Address","301234","test"],"id":1}`,
			unmarshalled: &hcashjson.VerifyMessageCmd{
				Address:   "1Address",
				Signature: "301234",
				Message:   "test",
			},
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
			t.Errorf("\n%s\n%s", marshalled, test.marshalled)
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

// TestChainSvrCmdErrors ensures any errors that occur in the command during
// custom mashal and unmarshal are as expected.
func TestChainSvrCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     interface{}
		marshalled string
		err        error
	}{
		{
			name:       "template request with invalid type",
			result:     &hcashjson.TemplateRequest{},
			marshalled: `{"mode":1}`,
			err:        &json.UnmarshalTypeError{},
		},
		{
			name:       "invalid template request sigoplimit field",
			result:     &hcashjson.TemplateRequest{},
			marshalled: `{"sigoplimit":"invalid"}`,
			err:        hcashjson.Error{Code: hcashjson.ErrInvalidType},
		},
		{
			name:       "invalid template request sizelimit field",
			result:     &hcashjson.TemplateRequest{},
			marshalled: `{"sizelimit":"invalid"}`,
			err:        hcashjson.Error{Code: hcashjson.ErrInvalidType},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.Unmarshal([]byte(test.marshalled), &test.result)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error type - got `%T` (%v), got `%T`",
				i, test.name, err, err, test.err)
			continue
		}

		if terr, ok := test.err.(hcashjson.Error); ok {
			gotErrorCode := err.(hcashjson.Error).Code
			if gotErrorCode != terr.Code {
				t.Errorf("Test #%d (%s) mismatched error code "+
					"- got %v (%v), want %v", i, test.name,
					gotErrorCode, terr, terr.Code)
				continue
			}
		}
	}
}
