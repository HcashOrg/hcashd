// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2017 The Hcash developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"bytes"
	"fmt"
	"io"

	"github.com/HcashOrg/hcashd/chaincfg/chainhash"
)

// blockHeaderLen is a constant that represents the number of bytes for a block
// header.
const lightBlockHeaderLen = 180 + 32 + 4

// MsgBlock implements the Message interface and represents a hypercash
// block message.  It is used to deliver block and transaction information in
// response to a getdata message (MsgGetData) for a given block hash.
type MsgLightBlock struct {
	Header        BlockHeader
	TxIds  []*chainhash.Hash
	STxIds []*chainhash.Hash
}

func (msg *MsgLightBlock) PrintMsgLightBlock(start string) {
	fmt.Printf("[test]%v\n", start)
	fmt.Printf("[test]block Hash:%v \n", msg.Header.BlockHash())

	for _, txid := range msg.TxIds{
		fmt.Printf("[test]txid:%v \n", txid)
	}

	for _, stxid := range msg.STxIds{
		fmt.Printf("[test]stxid:%v \n", stxid)
	}

	fmt.Printf("[test]End Block\n")
}

// AddTransaction adds a transaction to the message.
func (msg *MsgLightBlock) AddTransactionID(txid chainhash.Hash) error {
	msg.TxIds = append(msg.TxIds, &txid)
	return nil

}

// AddSTransaction adds a stake transaction to the message.
func (msg *MsgLightBlock) AddSTransactionID(txid chainhash.Hash) error {
	msg.STxIds = append(msg.STxIds, &txid)
	return nil
}

// ClearTransactions removes all transactions from the message.
func (msg *MsgLightBlock) ClearTransactionIDs() {
	msg.TxIds = make([]*chainhash.Hash, 0, defaultTransactionAlloc)
}

// ClearSTransactions removes all stake transactions from the message.
func (msg *MsgLightBlock) ClearSTransactionIDs() {
	msg.STxIds = make([]*chainhash.Hash, 0, defaultTransactionAlloc)
}

// BtcDecode decodes r using the hypercash protocol encoding into the receiver.
// This is part of the Message interface implementation.
// See Deserialize for decoding blocks stored to disk, such as in a database, as
// opposed to decoding blocks from the wire.
func (msg *MsgLightBlock) BtcDecode(r io.Reader, pver uint32) error {
	err := readBlockHeader(r, pver, &msg.Header)
	if err != nil {
		return err
	}

	txCount, err := ReadVarInt(r, pver)
	if err != nil {
		return err
	}

	// Prevent more transactions than could possibly fit into the regular
	// tx tree.
	// It would be possible to cause memory exhaustion and panics without
	// a sane upper bound on this count.
	maxTxPerTree := MaxTxPerTxTree(pver)
	if txCount > maxTxPerTree {
		str := fmt.Sprintf("too many transactions to fit into a block "+
			"[count %d, max %d]", txCount, maxTxPerTree)
		return messageError("MsgBlock.BtcDecode", str)
	}

	msg.TxIds = make([]*chainhash.Hash, 0, txCount)
	for i := uint64(0); i < txCount; i++ {
		var txId chainhash.Hash
		readElement(r, &txId)
		if err != nil {
			return err
		}
		msg.TxIds = append(msg.TxIds, &txId)
	}

	// Prevent more transactions than could possibly fit into the stake
	// tx tree.
	// It would be possible to cause memory exhaustion and panics without
	// a sane upper bound on this count.
	stakeTxCount, err := ReadVarInt(r, pver)
	if err != nil {
		return err
	}
	if stakeTxCount > maxTxPerTree {
		str := fmt.Sprintf("too many stransactions to fit into a block "+
			"[count %d, max %d]", stakeTxCount, maxTxPerTree)
		return messageError("MsgBlock.BtcDecode", str)
	}

	msg.STxIds = make([]*chainhash.Hash, 0, stakeTxCount)
	for i := uint64(0); i < stakeTxCount; i++ {
		var stxId chainhash.Hash
		readElement(r, &stxId)
		if err != nil {
			return err
		}
		msg.TxIds = append(msg.STxIds, &stxId)
	}
	msg.PrintMsgLightBlock("BtcDecode LightBlock")
	return nil
}

// Deserialize decodes a block from r into the receiver using a format that is
// suitable for long-term storage such as a database while respecting the
// Version field in the block.  This function differs from BtcDecode in that
// BtcDecode decodes from the hypercash wire protocol as it was sent across the
// network.  The wire encoding can technically differ depending on the protocol
// version and doesn't even really need to match the format of a stored block at
// all.  As of the time this comment was written, the encoded block is the same
// in both instances, but there is a distinct difference and separating the two
// allows the API to be flexible enough to deal with changes.
func (msg *MsgLightBlock) Deserialize(r io.Reader) error {
	// At the current time, there is no difference between the wire encoding
	// at protocol version 0 and the stable long-term storage format.  As
	// a result, make use of BtcDecode.
	return msg.BtcDecode(r, 0)
}

// FromBytes deserializes a transaction byte slice.
func (msg *MsgLightBlock) FromBytes(b []byte) error {
	r := bytes.NewReader(b)
	return msg.Deserialize(r)
}

// BtcEncode encodes the receiver to w using the hypercash protocol encoding.
// This is part of the Message interface implementation.
// See Serialize for encoding blocks to be stored to disk, such as in a
// database, as opposed to encoding blocks for the wire.
func (msg *MsgLightBlock) BtcEncode(w io.Writer, pver uint32) error {
	msg.PrintMsgLightBlock("BtcEncode LightBlock")
	err := writeBlockHeader(w, pver, &msg.Header)
	if err != nil {
		return err
	}

	err = WriteVarInt(w, pver, uint64(len(msg.TxIds)))
	if err != nil {
		return err
	}

	for _, txid := range msg.TxIds {
		err := writeElement(w, txid)
		if err != nil {
			return err
		}
	}

	err = WriteVarInt(w, pver, uint64(len(msg.STxIds)))
	if err != nil {
		return err
	}

	for _, stxid := range msg.STxIds {
		err := writeElement(w, stxid)
		if err != nil {
			return err
		}
	}

	return nil
}

// Serialize encodes the block to w using a format that suitable for long-term
// storage such as a database while respecting the Version field in the block.
// This function differs from BtcEncode in that BtcEncode encodes the block to
// the hypercash wire protocol in order to be sent across the network.  The wire
// encoding can technically differ depending on the protocol version and doesn't
// even really need to match the format of a stored block at all.  As of the
// time this comment was written, the encoded block is the same in both
// instances, but there is a distinct difference and separating the two allows
// the API to be flexible enough to deal with changes.
func (msg *MsgLightBlock) Serialize(w io.Writer) error {
	// At the current time, there is no difference between the wire encoding
	// at protocol version 0 and the stable long-term storage format.  As
	// a result, make use of BtcEncode.
	return msg.BtcEncode(w, 0)
}

// Bytes returns the serialized form of the block in bytes.
func (msg *MsgLightBlock) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, msg.SerializeSize()))
	err := msg.Serialize(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// SerializeSize returns the number of bytes it would take to serialize the
// the block.
func (msg *MsgLightBlock) SerializeSize() int {
	// Check to make sure that all transactions have the correct
	// type and version to be included in a block.

	// Block header bytes + Serialized varint size for the number of
	// transactions + Serialized varint size for the number of
	// stake transactions

	n := lightBlockHeaderLen + uint64(len(msg.TxIds)) * chainhash.HashSize +
		uint64(len(msg.STxIds)) * chainhash.HashSize


	return int(n)
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgLightBlock) Command() string {
	return CmdLightBlock
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgLightBlock) MaxPayloadLength(pver uint32) uint32 {
	// Protocol version 3 and lower have a different max block payload.
	if pver <= 3 {
		return MaxBlockPayloadV3
	}

	// Block header at 80 bytes + transaction count + max transactions
	// which can vary up to the MaxBlockPayload (including the block header
	// and transaction count).
	return MaxBlockPayload
}

// BlockHash computes the block identifier hash for this block.
func (msg *MsgLightBlock) BlockHash() chainhash.Hash {
	return msg.Header.BlockHash()
}

// NewMsgBlock returns a new hypercash block message that conforms to the
// Message interface.  See MsgBlock for details.
func NewMsgLightBlock(blockHeader *BlockHeader) *MsgLightBlock {
	return &MsgLightBlock{
		Header:        *blockHeader,
		TxIds:  make([]*chainhash.Hash, 0, defaultTransactionAlloc),
		STxIds: make([]*chainhash.Hash, 0, defaultTransactionAlloc),
	}
}

// NewMsgBlock returns a new hypercash block message that conforms to the
// Message interface.  See MsgBlock for details.
func NewMsgLightBlockFromMsgBlock(msgBlock *MsgBlock) *MsgLightBlock {
	msgLightBlock := &MsgLightBlock{
		Header:        msgBlock.Header,
		TxIds:  make([]*chainhash.Hash, 0, defaultTransactionAlloc),
		STxIds: make([]*chainhash.Hash, 0, defaultTransactionAlloc),
	}
	for _, tx := range msgBlock.Transactions {
		msgLightBlock.AddTransactionID(tx.TxHash())
	}
	for _, stx := range msgBlock.STransactions {
		msgLightBlock.AddSTransactionID(stx.TxHash())
	}

	return msgLightBlock
}
