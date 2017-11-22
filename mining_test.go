package main

import (
	"container/heap"
	"math/rand"
	"testing"

	"github.com/HcashOrg/hcashd/blockchain/stake"
)

// fakePrioritizedTxes prepares some fake prioritized txes for test
func fakePrioritizedTxes() []*txPrioItem {
	const numRandTx = 1000
	// some items under the edge condition
	prioritizedTxes := []*txPrioItem{
		{feePerKB: 5678, txType: stake.TxTypeRegular, priority: 3},
		{feePerKB: 5678, txType: stake.TxTypeRegular, priority: 1},
		{feePerKB: 5678, txType: stake.TxTypeRegular, priority: 1}, // Duplicate fee and prio
		{feePerKB: 5678, txType: stake.TxTypeRegular, priority: 5},
		{feePerKB: 5678, txType: stake.TxTypeRegular, priority: 2},
		{feePerKB: 1234, txType: stake.TxTypeRegular, priority: 3},
		{feePerKB: 1234, txType: stake.TxTypeRegular, priority: 1},
		{feePerKB: 1234, txType: stake.TxTypeRegular, priority: 5},
		{feePerKB: 1234, txType: stake.TxTypeRegular, priority: 5}, // Duplicate fee and prio
		{feePerKB: 1234, txType: stake.TxTypeRegular, priority: 2},
		{feePerKB: 10000, txType: stake.TxTypeRegular, priority: 0}, // Higher fee, smaller prio
		{feePerKB: 0, txType: stake.TxTypeRegular, priority: 10000}, // Higher prio, lower fee
	}

	for i := 0; i < numRandTx; i++ {
		// fake a random tx
		randTxType := stake.TxType(rand.Intn(4))
		randPriority := rand.Float64() * 100
		randFeePerKB := rand.Float64() * 10

		prioritizedTxes = append(prioritizedTxes, &txPrioItem{
			tx:       nil,
			txType:   randTxType,
			feePerKB: randFeePerKB,
			priority: randPriority,
		})
	}

	return prioritizedTxes
}

// TestTxPQOnStakePriorityAndFeeAndTxPriority tests the priority
// queue on stake priority, fee per kb and tx priority
func TestTxPQOnStakePriorityAndFeeAndTxPriority(t *testing.T) {
	// get fake txes as samples
	prioritizedTxes := fakePrioritizedTxes()
	pq := newTxPriorityQueue(len(prioritizedTxes), txPQByStakeAndFee)

	// build the priority queue tx by tx
	for _, tx := range prioritizedTxes {
		heap.Push(pq, tx)
	}

	// ensure the size of queue is correct
	if pq.Len() != len(prioritizedTxes) {
		t.Errorf("size of priority queue built is %v, want %v\n", pq.Len(), len(prioritizedTxes))
	}

	// last item popped out
	prev := &txPrioItem{
		tx:       nil,
		txType:   stake.TxTypeSSGen,
		priority: 10000.0, // since highest priority of the fake tx is 10000
		feePerKB: 10000.0, // since highest feePerKB of the fake tx is 10000
	}

	for pq.Len() > 0 {
		item := heap.Pop(pq)
		if tx, ok := item.(*txPrioItem); ok {
			// check correctness of order
			// higher stake priority, fee per kb and tx priority comes first
			if (compareStakePriority(tx, prev) > 0) ||
				((0 == compareStakePriority(tx, prev)) && (tx.feePerKB > prev.feePerKB)) {
				t.Errorf("bad pop: %v fee per KB was more than previous of %v "+
					"while the txtype was %v but previous was %v",
					tx.feePerKB, prev.feePerKB, tx.txType, prev.txType)
			}

			prev = tx
		}
	}
}

// TestTxPQOnStakePriorityAndFeeAndTxPriorityConsideringTxType tests the priority
// queue on stake priority, fee per kb
// if both tx are of type regular or revocation, plus the tx priority
func TestTxPQOnStakePriorityAndFeeAndConditionalTxPriority(t *testing.T) {
	// get fake txes as samples
	prioritizedTxes := fakePrioritizedTxes()
	pq := newTxPriorityQueue(len(prioritizedTxes), txPQByStakeAndFeeAndThenPriority)

	// build the priority queue tx by tx
	for _, tx := range prioritizedTxes {
		heap.Push(pq, tx)
	}

	// ensure the size of queue is correct
	if pq.Len() != len(prioritizedTxes) {
		t.Errorf("size of priority queue built is %v, want %v\n", pq.Len(), len(prioritizedTxes))
	}

	// last item popped out
	prev := &txPrioItem{
		tx:       nil,
		txType:   stake.TxTypeSSGen,
		priority: 10000.0, // since highest priority of the fake tx is 10000
		feePerKB: 10000.0, // since highest feePerKB of the fake tx is 10000
	}

	for pq.Len() > 0 {
		item := heap.Pop(pq)
		if tx, ok := item.(*txPrioItem); ok {
			// check correctness of order
			// higher stake priority, fee per kb
			// if both tx types are either regular or revocation, plus priority
			// comes first

			stakePriorityDelta := compareStakePriority(tx, prev)
			if (txStakePriority(tx.txType) == regOrRevocPriority) &&
				(txStakePriority(prev.txType) == regOrRevocPriority) {
				// both are of low stake priority

				if (stakePriorityDelta > 0) ||
					((0 == stakePriorityDelta) && (tx.priority > prev.priority)) {
					t.Errorf("bad pop: %v priority was more than previous of %v "+
						"while the tx type was %v but previous was %v",
						tx.priority, prev.priority, tx.txType, prev.txType)
				}
			} else {
				// neither are of low stake priority

				if (stakePriorityDelta > 0) ||
					((0 == stakePriorityDelta) && (tx.feePerKB > prev.feePerKB)) {
					t.Errorf("bad pop: %v fee per KB was more than previous of %v "+
						"while the tx type was %v but previous was %v",
						tx.feePerKB, prev.feePerKB, tx.txType, prev.txType)
				}
			}

			prev = tx
		}
	}
}
