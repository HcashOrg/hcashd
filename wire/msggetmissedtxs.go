// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"
	"github.com/HcashOrg/hcashd/chaincfg/chainhash"
)

// defaultInvListAlloc is the default size used for the backing array for an
// inventory list.  The array will dynamically grow as needed, but this
// figure is intended to provide enough space for the max number of inventory
// vectors in a *typical* inventory message without needing to grow the backing
// array multiple times.  Technically, the list can grow to MaxInvPerMsg, but
// rather than using that large figure, this figure more accurately reflects the
// typical case.
const defaultTxInvListAlloc = 1000

// MsgInv implements the Message interface and represents a hypercash inv message.
// It is used to advertise a peer's known data such as blocks and transactions
// through inventory vectors.  It may be sent unsolicited to inform other peers
// of the data or in response to a getblocks message (MsgGetBlocks).  Each
// message is limited to a maximum number of inventory vectors, which is
// currently 50,000.
//
// Use the AddInvVect function to build up the list of inventory vectors when
// sending an inv message to another peer.
type MsgGetMissedTxs struct {
	BlockInv *InvVect
	TxInvList []*InvVect
}

func (msg *MsgGetMissedTxs) PrintMsgGetMissedTxs(start string) {
	fmt.Printf("[test]%v\n", start)
	fmt.Printf("[test]Block inv Type:%v \n", msg.BlockInv.Type)
	fmt.Printf("[test]Block inv Hash:%v \n", msg.BlockInv.Hash)

	for _, tx := range msg.TxInvList{
		fmt.Printf("[test]tx type:%v \n", tx.Type)
		fmt.Printf("[test]tx id:%v \n", tx.Hash)
	}

	fmt.Printf("[test]End TxInvVect\n")
}

// AddInvVect adds an inventory vector to the message.
func (msg *MsgGetMissedTxs) AddTxInvVect(iv *InvVect) error {
	if len(msg.TxInvList)+1 > MaxInvPerMsg {
		str := fmt.Sprintf("too many invvect in message [max %v]",
			MaxInvPerMsg)
		return messageError("MsgInv.AddInvVect", str)
	}

	msg.TxInvList = append(msg.TxInvList, iv)
	return nil
}

// AddInvVect adds an inventory vector to the message.
func (msg *MsgGetMissedTxs) SetBlockInv(blockHash *chainhash.Hash) error {
	iv := NewInvVect(InvTypeLightBlock, blockHash)
	msg.BlockInv = iv
	return nil
}

// BtcDecode decodes r using the hypercash protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgGetMissedTxs) BtcDecode(r io.Reader, pver uint32) error {
	count, err := ReadVarInt(r, pver)
	if err != nil {
		return err
	}

	// Limit to max inventory vectors per message.
	if count > MaxInvPerMsg {
		str := fmt.Sprintf("too many invvect in message [%v]", count)
		return messageError("MsgInv.BtcDecode", str)
	}
	
	blockInvList := make([]InvVect, 1)
	err = readInvVect(r, pver, &blockInvList[0])
	if err != nil {
		return err
	}
	msg.BlockInv = &blockInvList[0]

	// Create a contiguous slice of inventory vectors to deserialize into in
	// order to reduce the number of allocations.
	invList := make([]InvVect, count)
	msg.TxInvList = make([]*InvVect, 0, count)
	for i := uint64(0); i < count; i++ {
		iv := &invList[i]
		err := readInvVect(r, pver, iv)
		if err != nil {
			return err
		}
		msg.AddTxInvVect(iv)
	}
	msg.PrintMsgGetMissedTxs("BtcDecode MsgGetMissedTxs")
	return nil
}

// BtcEncode encodes the receiver to w using the hypercash protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgGetMissedTxs) BtcEncode(w io.Writer, pver uint32) error {
	// Limit to max inventory vectors per message.
	msg.PrintMsgGetMissedTxs("BtcEncode MsgGetMissedTxs")
	count := len(msg.TxInvList)
	if count > MaxInvPerMsg {
		str := fmt.Sprintf("too many invvect in message [%v]", count)
		return messageError("MsgInv.BtcEncode", str)
	}

	err := WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	//Write Block info
	err = writeInvVect(w, pver, msg.BlockInv)
	if err != nil {
		return err
	}

	for _, iv := range msg.TxInvList {
		err := writeInvVect(w, pver, iv)
		if err != nil {
			return err
		}
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgGetMissedTxs) Command() string {
	return CmdGetMissedTxs
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgGetMissedTxs) MaxPayloadLength(pver uint32) uint32 {
	// Num inventory vectors (varInt) + max allowed inventory vectors.
	return MaxVarIntPayload + (MaxInvPerMsg * maxInvVectPayload)
}

// NewMsgInv returns a new hypercash inv message that conforms to the Message
// interface.  See MsgInv for details.
func NewMsgGetMissedTxs() *MsgGetMissedTxs {
	return &MsgGetMissedTxs{
		TxInvList: make([]*InvVect, 0, defaultTxInvListAlloc),
	}
}

// NewMsgInvSizeHint returns a new hypercash inv message that conforms to the
// Message interface.  See MsgInv for details.  This function differs from
// NewMsgInv in that it allows a default allocation size for the backing array
// which houses the inventory vector list.  This allows callers who know in
// advance how large the inventory list will grow to avoid the overhead of
// growing the internal backing array several times when appending large amounts
// of inventory vectors with AddInvVect.  Note that the specified hint is just
// that - a hint that is used for the default allocation size.  Adding more
// (or less) inventory vectors will still work properly.  The size hint is
// limited to MaxInvPerMsg.
func NewMsgGetMissedTxsSizeHint(sizeHint uint) *MsgGetMissedTxs {
	// Limit the specified hint to the maximum allow per message.
	if sizeHint > MaxInvPerMsg {
		sizeHint = MaxInvPerMsg
	}

	return &MsgGetMissedTxs{
		TxInvList: make([]*InvVect, 0, sizeHint),
	}
}
