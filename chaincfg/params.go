// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2017 The Hcash developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/HcashOrg/hcashd/chaincfg/chainhash"
	"github.com/HcashOrg/hcashd/wire"
)

// These variables are the chain proof-of-work limit parameters for each default
// network.
var (
	// bigOne is 1 represented as a big.Int.  It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)

	// mainPowLimit is the highest proof of work value a Hypercash block can
	// have for the main network.  It is the value 2^224 - 1.
	mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 228), bigOne)

	// testNetPowLimit is the highest proof of work value a Hypercash block
	// can have for the test network.  It is the value 2^232 - 1.
	testNetPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 232), bigOne)

	// simNetPowLimit is the highest proof of work value a Hypercash block
	// can have for the simulation test network.  It is the value 2^255 - 1.
	simNetPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)

	VoteBitsNotFound = fmt.Errorf("vote bits not found")
)

// SigHashOptimization is an optimization for verification of transactions that
// do CHECKSIG operations with hashType SIGHASH_ALL. Although there should be no
// consequences to daemons that are simply running a node, it may be the case
// that you could cause database corruption if you turn this code on, create and
// manipulate your own MsgTx, then include them in blocks. For safety, if you're
// using the daemon with wallet or mining with the daemon this should be disabled.
// If you believe that any MsgTxs in your daemon will be used mutably, do NOT
// turn on this feature. It is disabled by default.
// This feature is considered EXPERIMENTAL, enable at your own risk!
var SigHashOptimization = false

// CheckForDuplicateHashes checks for duplicate hashes when validating blocks.
// Because of the rule inserting the height into the second (nonce) txOut, there
// should never be a duplicate transaction hash that overwrites another. However,
// because there is a 2^128 chance of a collision, the paranoid user may wish to
// turn this feature on.
var CheckForDuplicateHashes = false

// CPUMinerThreads is the default number of threads to utilize with the
// CPUMiner when mining.
var CPUMinerThreads = 1

// Checkpoint identifies a known good point in the block chain.  Using
// checkpoints allows a few optimizations for old blocks during initial download
// and also prevents forks from old blocks.
//
// Each checkpoint is selected based upon several factors.  See the
// documentation for chain.IsCheckpointCandidate for details on the selection
// criteria.
type Checkpoint struct {
	Height int64
	Hash   *chainhash.Hash
}

// Vote describes a voting instance.  It is self-describing so that the UI can
// be directly implemented using the fields.  Mask determines which bits can be
// used.  Bits are enumerated and must be consecutive.  Each vote requires one
// and only one abstain (bits = 0) and reject vote (IsNo = true).
//
// For example, change block height from int64 to uint64.
// Vote {
//	Id:          "blockheight",
//	Description: "Change block height from int64 to uint64"
//	Mask:        0x0006,
//	Choices:     []Choice{
//		{
//			Id:          "abstain",
//			Description: "abstain voting for change",
//			Bits:        0x0000,
//			IsAbstain:   true,
//			IsNo:        false,
//		},
//		{
//			Id:          "no",
//			Description: "reject changing block height to uint64",
//			Bits:        0x0002,
//			IsAbstain:   false,
//			IsNo:        false,
//		},
//		{
//			Id:          "yes",
//			Description: "accept changing block height to uint64",
//			Bits:        0x0004,
//			IsAbstain:   false,
//			IsNo:        true,
//		},
//	},
// }
//
type Vote struct {
	// Single unique word identifying the vote.
	Id string

	// Longer description of what the vote is about.
	Description string

	// Usable bits for this vote.
	Mask uint16

	Choices []Choice
}

// Choice defines one of the possible Choices that make up a vote. The 0 value
// in Bits indicates the default choice.  Care should be taken not to bias a
// vote with the default choice.
type Choice struct {
	// Single unique word identifying vote (e.g. yes)
	Id string

	// Longer description of the vote.
	Description string

	// Bits used for this vote.
	Bits uint16

	// This is the abstain choice.  By convention this must be the 0 vote
	// (abstain) and exist only once in the Vote.Choices array.
	IsAbstain bool

	// This coince indicates a hard No Vote.  By convention this must exist
	// only once in the Vote.Choices array.
	IsNo bool
}

// VoteIndex compares vote to Choice.Bits and returns the index into the
// Choices array.  If the vote is invalid it returns -1.
func (v *Vote) VoteIndex(vote uint16) int {
	vote &= v.Mask
	for k := range v.Choices {
		if vote == v.Choices[k].Bits {
			return k
		}
	}

	return -1
}

const (
	// VoteIDMaxBlockSize is the vote ID for the the maximum block size
	// increase agenda used for the hard fork demo.
	VoteIDMaxBlockSize = "maxblocksize"

	// VoteIDSDiffAlgorithm is the vote ID for the new stake difficulty
	// algorithm (aka ticket price) agenda defined by DCP0001.
	VoteIDSDiffAlgorithm = "sdiffalgorithm"

	// VoteIDLNSupport is the vote ID for determining if the developers
	// should work on integrating Lightning Network support.
	VoteIDLNSupport = "lnsupport"
)

// ConsensusDeployment defines details related to a specific consensus rule
// change that is voted in.  This is part of BIP0009.
type ConsensusDeployment struct {
	// Vote describes the what is being voted on and what the choices are.
	// This is sitting in a struct in order to make merging between btcd
	// easier.
	Vote Vote

	// StartTime is the median block time after which voting on the
	// deployment starts.
	StartTime uint64

	// ExpireTime is the median block time after which the attempted
	// deployment expires.
	ExpireTime uint64
}

// TokenPayout is a payout for block 1 which specifies an address and an amount
// to pay to that address in a transaction output.
type TokenPayout struct {
	Address string
	Amount  int64
}

// Params defines a Hypercash network by its parameters.  These parameters may be
// used by Hypercash applications to differentiate networks as well as addresses
// and keys for one network from those intended for use on another network.
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name string

	// Net defines the magic bytes used to identify the network.
	Net wire.CurrencyNet

	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string

	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds []string

	// GenesisBlock defines the first block of the chain.
	GenesisBlock *wire.MsgBlock

	// GenesisHash is the starting block hash.
	GenesisHash *chainhash.Hash

	// PowLimit defines the highest allowed proof of work value for a block
	// as a uint256.
	PowLimit *big.Int

	// DifficultyRate is used for micro block
	DifficultyRate uint32

	// MaxMicroPerKeys defines the max number of microblocks between two keyblocks
	MaxMicroPerKey uint32

	// PowLimitBits defines the highest allowed proof of work value for a
	// block in compact form.
	PowLimitBits uint32

	// ReduceMinDifficulty defines whether the network should reduce the
	// minimum required difficulty after a long enough period of time has
	// passed without finding a block.  This is really only useful for test
	// networks and should not be set on a main network.
	ReduceMinDifficulty bool

	// MinDiffReductionTime is the amount of time after which the minimum
	// required difficulty should be reduced when a block hasn't been found.
	//
	// NOTE: This only applies if ReduceMinDifficulty is true.
	MinDiffReductionTime time.Duration

	// GenerateSupported specifies whether or not CPU mining is allowed.
	GenerateSupported bool

	// MaximumBlockSizes are the maximum sizes of a block that can be
	// generated on the network.  It is an array because the max block size
	// can be different values depending on the results of a voting agenda.
	// The first entry is the initial block size for the network, while the
	// other entries are potential block size changes which take effect when
	// the vote for the associated agenda succeeds.
	MaximumBlockSizes []int

	// MaxTxSize is the maximum number of bytes a serialized transaction can
	// be in order to be considered valid by consensus.
	MaxTxSize int

	// TargetTimePerBlock is the desired amount of time to generate each
	// block.
	TargetTimePerBlock time.Duration

	// WorkDiffAlpha is the stake difficulty EMA calculation alpha (smoothing)
	// value. It is different from a normal EMA alpha. Closer to 1 --> smoother.
	WorkDiffAlpha int64

	// WorkDiffWindowSize is the number of windows (intervals) used for calculation
	// of the exponentially weighted average.
	WorkDiffWindowSize int64

	// WorkDiffWindows is the number of windows (intervals) used for calculation
	// of the exponentially weighted average.
	WorkDiffWindows int64

	// TargetTimespan is the desired amount of time that should elapse
	// before the block difficulty requirement is examined to determine how
	// it should be changed in order to maintain the desired block
	// generation rate.  This value should correspond to the product of
	// WorkDiffWindowSize and TimePerBlock above.
	TargetTimespan time.Duration

	// RetargetAdjustmentFactor is the adjustment factor used to limit
	// the minimum and maximum amount of adjustment that can occur between
	// difficulty retargets.
	RetargetAdjustmentFactor int64

	// Subsidy parameters.
	//
	// Subsidy calculation for exponential reductions:
	// 0 for i in range (0, height / SubsidyReductionInterval):
	// 1     subsidy *= MulSubsidy
	// 2     subsidy /= DivSubsidy
	//
	// Caveat: Don't overflow the int64 register!!

	// BaseSubsidy is the starting subsidy amount for mined blocks.
	BaseSubsidy int64

	// Subsidy reduction multiplier.
	MulSubsidy int64

	// Subsidy reduction divisor.
	DivSubsidy int64

	// SubsidyReductionInterval is the reduction interval in blocks.
	SubsidyReductionInterval int64

	// WorkRewardProportion is the comparative amount of the subsidy given for
	// creating a block.
	WorkRewardProportion uint16

	// StakeRewardProportion is the comparative amount of the subsidy given for
	// casting stake votes (collectively, per block).
	StakeRewardProportion uint16

	// BlockTaxProportion is the inverse of the percentage of funds for each
	// block to allocate to the developer organization.
	// e.g. 10% --> 10 (or 1 / (1/10))
	// Special case: disable taxes with a value of 0
	BlockTaxProportion uint16

	// Checkpoints ordered from oldest to newest.
	Checkpoints []Checkpoint

	// These fields are related to voting on consensus rule changes as
	// defined by BIP0009.
	//
	// RuleChangeActivationQurom is the number of votes required for a vote
	// to take effect.
	//
	// RuleChangeActivationInterval is the number of blocks in each threshold
	// state retarget window.
	//
	// Deployments define the specific consensus rule changes to be voted
	// on for the stake version (the map key).
	RuleChangeActivationQuorum     uint32
	RuleChangeActivationMultiplier uint32
	RuleChangeActivationDivisor    uint32
	RuleChangeActivationInterval   uint32
	Deployments                    map[uint32][]ConsensusDeployment

	// Enforce current block version once network has upgraded.
	BlockEnforceNumRequired uint64

	// Reject previous block versions once network has upgraded.
	BlockRejectNumRequired uint64

	// The number of nodes to check.
	BlockUpgradeNumToCheck uint64

	MicroBlockValidationHeight int64

	// Mempool parameters
	RelayNonStdTxs bool

	// NetworkAddressPrefix is the first letter of the network
	// for any given address encoded as a string.
	NetworkAddressPrefix string

	// Address encoding magics
	PubKeyAddrID     [2]byte // First 2 bytes of a P2PK address
	PubKeyBlissAddrID [2]byte
	PubKeyLmsAddrID  [2]byte
	PubKeyHashAddrID [2]byte // First 2 bytes of a P2PKH address
	PKHEdwardsAddrID [2]byte // First 2 bytes of an Edwards P2PKH address
	PKHSchnorrAddrID [2]byte // First 2 bytes of a secp256k1 Schnorr P2PKH address
	PKHBlissAddrID   [2]byte // First 2 bytes of a Bliss address
	PKHLmsAddrID   	 [2]byte // First 2 bytes of a LMS address
	ScriptHashAddrID [2]byte // First 2 bytes of a P2SH address
	PrivateKeyID     [2]byte // First 2 bytes of a WIF private key

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID [4]byte
	HDPublicKeyID  [4]byte

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType uint32

	// MinimumStakeDiff if the minimum amount of Atoms required to purchase a
	// stake ticket.
	MinimumStakeDiff int64

	// Ticket pool sizes for Hypercash PoS. This denotes the number of possible
	// buckets/number of different ticket numbers. It is also the number of
	// possible winner numbers there are.
	TicketPoolSize uint16

	// Average number of tickets per block for Hypercash PoS.
	TicketsPerBlock uint16

	// Number of blocks for tickets to mature (spendable at TicketMaturity+1).
	TicketMaturity uint16

	// Number of blocks for tickets to expire after they have matured. This MUST
	// be >= (StakeEnabledHeight + StakeValidationHeight).
	TicketExpiry uint32

	// CoinbaseMaturity is the number of blocks required before newly mined
	// coins (coinbase transactions) can be spent.
	CoinbaseMaturity uint16

	// Maturity for spending SStx change outputs.
	SStxChangeMaturity uint16

	// TicketPoolSizeWeight is the multiplicative weight applied to the
	// ticket pool size difference between a window period and its target
	// when determining the stake system.
	TicketPoolSizeWeight uint16

	// StakeDiffAlpha is the stake difficulty EMA calculation alpha (smoothing)
	// value. It is different from a normal EMA alpha. Closer to 1 --> smoother.
	StakeDiffAlpha int64

	// StakeDiffWindowSize is the number of blocks used for each interval in
	// exponentially weighted average.
	StakeDiffWindowSize int64

	// StakeDiffWindows is the number of windows (intervals) used for calculation
	// of the exponentially weighted average.
	StakeDiffWindows int64

	// StakeVersionInterval determines the interval where the stake version
	// is calculated.
	StakeVersionInterval int64

	// MaxFreshStakePerBlock is the maximum number of new tickets that may be
	// submitted per block.
	MaxFreshStakePerBlock uint8

	// StakeEnabledHeight is the height in which the first ticket could possibly
	// mature.
	StakeEnabledHeight int64

	// StakeValidationHeight is the height at which votes (SSGen) are required
	// to add a new block to the top of the blockchain. This height is the
	// first block that will be voted on, but will include in itself no votes.
	StakeValidationHeight int64

	// StakeBaseSigScript is the consensus stakebase signature script for all
	// votes on the network. This isn't signed in any way, so without forcing
	// it to be this value miners/daemons could freely change it.
	StakeBaseSigScript []byte

	// StakeMajorityMultiplier and StakeMajorityDivisor are used
	// to calculate the super majority of stake votes using integer math as
	// such: X*StakeMajorityMultiplier/StakeMajorityDivisor
	StakeMajorityMultiplier int32
	StakeMajorityDivisor    int32

	// OrganizationPkScript is the output script for block taxes to be
	// distributed to in every block's coinbase. It should ideally be a P2SH
	// multisignature address.  OrganizationPkScriptVersion is the version
	// of the output script.  Until PoS hardforking is implemented, this
	// version must always match for a block to validate.
	OrganizationPkScript        []byte
	OrganizationPkScriptVersion uint16

	// BlockOneLedger specifies the list of payouts in the coinbase of
	// block height 1. If there are no payouts to be given, set this
	// to an empty slice.
	BlockOneLedger []*TokenPayout
}

// MainNetParams defines the network parameters for the main Hypercash network.
var MainNetParams = Params{
	Name:        "testdata2",
	Net:         wire.MainNet,
	DefaultPort: "14008",
	DNSSeeds: []string{
		"testnet-seeds.hcashtech.org",
	},

	// Chain parameters
	GenesisBlock:             &genesisBlock,
	GenesisHash:              &genesisHash,
	PowLimit:                 mainPowLimit,
	DifficultyRate:           16,
	MaxMicroPerKey:           31,
	PowLimitBits:             0x1d0fffff,
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0, // Does not apply since ReduceMinDifficulty false
	GenerateSupported:        false,
	MaximumBlockSizes:        []int{2048000},
	MaxTxSize:                2048000,
	TargetTimePerBlock:       time.Minute * 5,
	WorkDiffAlpha:            1,
	WorkDiffWindowSize:       144,
	WorkDiffWindows:          20,
	TargetTimespan:           time.Minute * 5 * 144, // TimePerBlock * WindowSize
	RetargetAdjustmentFactor: 4,

	// Subsidy parameters.
	BaseSubsidy:              5000000000, //
	MulSubsidy:               1000,
	DivSubsidy:               1005,
	SubsidyReductionInterval: 1543,
	WorkRewardProportion:     45,
	StakeRewardProportion:    45,
	BlockTaxProportion:       10,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
	//	{30,newHashFromStr("00000df9e4054bd941145c7ea9dbefc29e47ed564cc2fdb254720ab07a016938")},
	//	{200,newHashFromStr("00000019ed43fba03c72b03cbd7a706c50b56819379b478da57184363fd90a68")},
		{28952,newHashFromStr("00000005f6263ebb38395fcdc810bc462f34f063454056a4965cca9c3d280d70")},
	},

	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationQuorum:     4032, // 10 % of RuleChangeActivationInterval * TicketsPerBlock
	RuleChangeActivationMultiplier: 3,    // 75%
	RuleChangeActivationDivisor:    4,
	RuleChangeActivationInterval:   2016 * 4, // 4 weeks
	Deployments: map[uint32][]ConsensusDeployment{
		4: {{
			Vote: Vote{
				Id:          VoteIDSDiffAlgorithm,
				Description: "Change stake difficulty algorithm as defined in DCP0001",
				Mask:        0x0006, // Bits 1 and 2
				Choices: []Choice{{
					Id:          "abstain",
					Description: "abstain voting for change",
					Bits:        0x0000,
					IsAbstain:   true,
					IsNo:        false,
				}, {
					Id:          "no",
					Description: "keep the existing algorithm",
					Bits:        0x0002, // Bit 1
					IsAbstain:   false,
					IsNo:        true,
				}, {
					Id:          "yes",
					Description: "change to the new algorithm",
					Bits:        0x0004, // Bit 2
					IsAbstain:   false,
					IsNo:        false,
				}},
			},
			StartTime:  1493164800, // Apr 26th, 2017
			ExpireTime: 1524700800, // Apr 26th, 2018
		}, {
			Vote: Vote{
				Id:          VoteIDLNSupport,
				Description: "Request developers begin work on Lightning Network (LN) integration",
				Mask:        0x0018, // Bits 3 and 4
				Choices: []Choice{{
					Id:          "abstain",
					Description: "abstain from voting",
					Bits:        0x0000,
					IsAbstain:   true,
					IsNo:        false,
				}, {
					Id:          "no",
					Description: "no, do not work on integrating LN support",
					Bits:        0x0008, // Bit 3
					IsAbstain:   false,
					IsNo:        true,
				}, {
					Id:          "yes",
					Description: "yes, begin work on integrating LN support",
					Bits:        0x0010, // Bit 4
					IsAbstain:   false,
					IsNo:        false,
				}},
			},
			StartTime:  1493164800, // Apr 26th, 2017
			ExpireTime: 1508976000, // Oct 26th, 2017
		}},
	},

	// Enforce current block version once majority of the network has
	// upgraded.
	// 75% (750 / 1000)
	// Reject previous block versions once a majority of the network has
	// upgraded.
	// 95% (950 / 1000)
	BlockEnforceNumRequired: 750,
	BlockRejectNumRequired:  950,
	BlockUpgradeNumToCheck:  1000,

	MicroBlockValidationHeight: 64,

	// Mempool parameters
	RelayNonStdTxs: false,

	// Address encoding magics
	NetworkAddressPrefix: "H",
	PubKeyAddrID:         [2]byte{0x19, 0xa4}, // starts with Hk
	PubKeyBlissAddrID:    [2]byte{0x07, 0xc3}, // starts with Hk
	PubKeyLmsAddrID:      [2]byte{0x07, 0x77}, // starts with Hp
	PubKeyHashAddrID:     [2]byte{0x09, 0x7f}, // starts with Hs
	PKHEdwardsAddrID:     [2]byte{0x09, 0x60}, // starts with He
	PKHSchnorrAddrID:     [2]byte{0x09, 0x41}, // starts with HS
	PKHBlissAddrID:       [2]byte{0x09, 0x57}, // starts with Hb
	PKHLmsAddrID:		  [2]byte{0x09, 0x34}, // starts with HM
	ScriptHashAddrID:     [2]byte{0x09, 0x5a}, // starts with Hc
	PrivateKeyID:         [2]byte{0x19, 0xab}, // starts with Hm

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x02, 0xfd, 0xa4, 0xe8}, // starts with dprv
	HDPublicKeyID:  [4]byte{0x02, 0xfd, 0xa9, 0x26}, // starts with dpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 20,

	// Hypercash PoS parameters
	MinimumStakeDiff:        2 * 1e8, // 2 Coin
	TicketPoolSize:          8192,
	TicketsPerBlock:         5,
	TicketMaturity:          128/*256*/,
	TicketExpiry:            40960, // 5*TicketPoolSize
	CoinbaseMaturity:        128/*256*/,
	SStxChangeMaturity:      1,
	TicketPoolSizeWeight:    4,
	StakeDiffAlpha:          1, // Minimal
	StakeDiffWindowSize:     144,
	StakeDiffWindows:        20,
	StakeVersionInterval:    144 * 2 * 7, // ~1 week
	MaxFreshStakePerBlock:   20,          // 4*TicketsPerBlock
	StakeEnabledHeight:      128 + 128/*256 + 256*/,   // CoinbaseMaturity + TicketMaturity
	StakeValidationHeight:   512,        // ~14 days
	StakeBaseSigScript:      []byte{0x00, 0x00},
	StakeMajorityMultiplier: 3,
	StakeMajorityDivisor:    4,

	// Hypercash organization related parameters
	// Organization address is HcR2g1eGf6mpQz7QNa1u5AFeWq8yKnfRUeM
	OrganizationPkScript:        hexDecode("a914cc53f47615e01dcae979f5813278095627f095cc87"),  //HccmQdPB1hkFQcB7jE2drqRYbavEhoQuAdx
	OrganizationPkScriptVersion: 0,
	BlockOneLedger:              BlockOneLedgerMainNet,
}

// TestNet2Params defines the network parameters for the test currency network.
// This network is sometimes simply called "testnet".
// This is the second public iteration of testnet.
var TestNet2Params = Params{
	Name:        "testnet2",
	Net:         wire.TestNet2,
	DefaultPort: "12008",
	DNSSeeds: []string{
	},

	// Chain parameters
	GenesisBlock:             &testNet2GenesisBlock,
	GenesisHash:              &testNet2GenesisHash,
	PowLimit:                 testNetPowLimit,
	DifficultyRate:           16,
	MaxMicroPerKey:           31,
	PowLimitBits:             0x1e00ffff,
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0, // Does not apply since ReduceMinDifficulty false
	GenerateSupported:        true,
	MaximumBlockSizes:        []int{2048000},
	MaxTxSize:                2048000,
	TargetTimePerBlock:       time.Minute * 2,
	WorkDiffAlpha:            1,
	WorkDiffWindowSize:       144,
	WorkDiffWindows:          20,
	TargetTimespan:           time.Minute * 2 * 144, // TimePerBlock * WindowSize
	RetargetAdjustmentFactor: 4,

	// Subsidy parameters.
	BaseSubsidy:              5000000000, // 25 Coin
	MulSubsidy:               1000,
	DivSubsidy:               1005,
	SubsidyReductionInterval: 1543,
	WorkRewardProportion:     45,
	StakeRewardProportion:    45,
	BlockTaxProportion:       10,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationQuorum:     2520, // 10 % of RuleChangeActivationInterval * TicketsPerBlock
	RuleChangeActivationMultiplier: 3,    // 75%
	RuleChangeActivationDivisor:    4,
	RuleChangeActivationInterval:   5040, // 1 week
	Deployments: map[uint32][]ConsensusDeployment{
		5: {{
			Vote: Vote{
				Id:          VoteIDSDiffAlgorithm,
				Description: "Change stake difficulty algorithm as defined in DCP0001",
				Mask:        0x0006, // Bits 1 and 2
				Choices: []Choice{{
					Id:          "abstain",
					Description: "abstain voting for change",
					Bits:        0x0000,
					IsAbstain:   true,
					IsNo:        false,
				}, {
					Id:          "no",
					Description: "keep the existing algorithm",
					Bits:        0x0002, // Bit 1
					IsAbstain:   false,
					IsNo:        true,
				}, {
					Id:          "yes",
					Description: "change to the new algorithm",
					Bits:        0x0004, // Bit 2
					IsAbstain:   false,
					IsNo:        false,
				}},
			},
			StartTime:  1493164800, // Apr 26th, 2017
			ExpireTime: 1524700800, // Apr 26th, 2018
		}},
	},

	// Enforce current block version once majority of the network has
	// upgraded.
	// 51% (51 / 100)
	// Reject previous block versions once a majority of the network has
	// upgraded.
	// 75% (75 / 100)
	BlockEnforceNumRequired: 51,
	BlockRejectNumRequired:  75,
	BlockUpgradeNumToCheck:  100,
	MicroBlockValidationHeight: 256,

	// Mempool parameters
	RelayNonStdTxs: true,

	// Address encoding magics
	NetworkAddressPrefix: "T",
	PubKeyAddrID:         [2]byte{0x28, 0xf7}, // starts with Tk
	PubKeyBlissAddrID:    [2]byte{0x0e, 0xd1}, // starts with Tk
	PubKeyLmsAddrID:      [2]byte{0x0f, 0x19}, // starts with Tp
	PubKeyHashAddrID:     [2]byte{0x0f, 0x21}, // starts with Ts
	PKHEdwardsAddrID:     [2]byte{0x0f, 0x01}, // starts with Te
	PKHSchnorrAddrID:     [2]byte{0x0e, 0xe3}, // starts with TS
	PKHBlissAddrID:       [2]byte{0x0e, 0xf8}, // starts with Tb
	PKHLmsAddrID:		  [2]byte{0x0e, 0xd6}, // starts with TM
	ScriptHashAddrID:     [2]byte{0x0e, 0xfc}, // starts with Tc
	PrivateKeyID:         [2]byte{0x23, 0x0e}, // starts with Pt

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x97}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xd1}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 11,

	// Hypercash PoS parameters
	MinimumStakeDiff:        2 * 1e8, // 0.2 Coin
	TicketPoolSize:          8192,
	TicketsPerBlock:         5,
	TicketMaturity:          256,
	TicketExpiry:            40960, // 6*TicketPoolSize
	CoinbaseMaturity:        256,
	SStxChangeMaturity:      1,
	TicketPoolSizeWeight:    4,
	StakeDiffAlpha:          1,
	StakeDiffWindowSize:     144,
	StakeDiffWindows:        20,
	StakeVersionInterval:    144 * 2 * 7, // ~1 week
	MaxFreshStakePerBlock:   20,          // 4*TicketsPerBlock
	StakeEnabledHeight:      256 + 256,     // CoinbaseMaturity + TicketMaturity
	StakeValidationHeight:   1024,         // Arbitrary
	StakeBaseSigScript:      []byte{0x00, 0x00},
	StakeMajorityMultiplier: 3,
	StakeMajorityDivisor:    4,

	// Hypercash organization related parameters.
	// Organization address is TccTkqj8wFqrUemmHMRSx8SYEueQYLmuuFk
	OrganizationPkScript:        hexDecode("4fa6cbd0dbe5ec407fe4c8ad374e667771fa0d44"),
	OrganizationPkScriptVersion: 0,
	BlockOneLedger:              BlockOneLedgerTestNet2,
}

// SimNetParams defines the network parameters for the simulation test Hypercash
// network.  This network is similar to the normal test network except it is
// intended for private use within a group of individuals doing simulation
// testing.  The functionality is intended to differ in that the only nodes
// which are specifically specified are used to create the network rather than
// following normal discovery rules.  This is important as otherwise it would
// just turn into another public testnet.
var SimNetParams = Params{
	Name:        "simnet",
	Net:         wire.SimNet,
	DefaultPort: "13008",
	DNSSeeds:    []string{}, // NOTE: There must NOT be any seeds.

	// Chain parameters
	GenesisBlock:             &simNetGenesisBlock,
	GenesisHash:              &simNetGenesisHash,
	PowLimit:                 simNetPowLimit,
	DifficultyRate:           16,
	MaxMicroPerKey:           31,
	PowLimitBits:             0x207fffff,
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0, // Does not apply since ReduceMinDifficulty false
	GenerateSupported:        true,
	MaximumBlockSizes:        []int{2048000},
	MaxTxSize:                2048000,
	TargetTimePerBlock:       time.Second,
	WorkDiffAlpha:            1,
	WorkDiffWindowSize:       8,
	WorkDiffWindows:          4,
	TargetTimespan:           time.Second * 8, // TimePerBlock * WindowSize
	RetargetAdjustmentFactor: 4,

	// Subsidy parameters.
	BaseSubsidy:              50000000000,
	MulSubsidy:               100,
	DivSubsidy:               101,
	SubsidyReductionInterval: 128,
	WorkRewardProportion:     45,
	StakeRewardProportion:    45,
	BlockTaxProportion:       10,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationQuorum:     160, // 10 % of RuleChangeActivationInterval * TicketsPerBlock
	RuleChangeActivationMultiplier: 3,   // 75%
	RuleChangeActivationDivisor:    4,
	RuleChangeActivationInterval:   320, // 320 seconds
	Deployments: map[uint32][]ConsensusDeployment{
		4: {{
			Vote: Vote{
				Id:          VoteIDMaxBlockSize,
				Description: "Change maximum allowed block size from 1MiB to 1.25MB",
				Mask:        0x0006, // Bits 1 and 2
				Choices: []Choice{{
					Id:          "abstain",
					Description: "abstain voting for change",
					Bits:        0x0000,
					IsAbstain:   true,
					IsNo:        false,
				}, {
					Id:          "no",
					Description: "reject changing max allowed block size",
					Bits:        0x0002, // Bit 1
					IsAbstain:   false,
					IsNo:        true,
				}, {
					Id:          "yes",
					Description: "accept changing max allowed block size",
					Bits:        0x0004, // Bit 2
					IsAbstain:   false,
					IsNo:        false,
				}},
			},
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		}},
		5: {{
			Vote: Vote{
				Id:          VoteIDSDiffAlgorithm,
				Description: "Change stake difficulty algorithm as defined in DCP0001",
				Mask:        0x0006, // Bits 1 and 2
				Choices: []Choice{{
					Id:          "abstain",
					Description: "abstain voting for change",
					Bits:        0x0000,
					IsAbstain:   true,
					IsNo:        false,
				}, {
					Id:          "no",
					Description: "keep the existing algorithm",
					Bits:        0x0002, // Bit 1
					IsAbstain:   false,
					IsNo:        true,
				}, {
					Id:          "yes",
					Description: "change to the new algorithm",
					Bits:        0x0004, // Bit 2
					IsAbstain:   false,
					IsNo:        false,
				}},
			},
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		}},
	},

	// Enforce current block version once majority of the network has
	// upgraded.
	// 51% (51 / 100)
	// Reject previous block versions once a majority of the network has
	// upgraded.
	// 75% (75 / 100)
	BlockEnforceNumRequired: 51,
	BlockRejectNumRequired:  75,
	BlockUpgradeNumToCheck:  100,
	MicroBlockValidationHeight:256,

	// Mempool parameters
	RelayNonStdTxs: true,

	// Address encoding magics
	NetworkAddressPrefix: "S",
	PubKeyAddrID:         [2]byte{0x27, 0x6f}, // starts with Sk
	PubKeyBlissAddrID:    [2]byte{0x0e, 0x40}, // starts with SK
	PubKeyLmsAddrID:      [2]byte{0x0e, 0x89}, // starts with Sp
	PubKeyHashAddrID:     [2]byte{0x0e, 0x91}, // starts with Ss
	PKHEdwardsAddrID:     [2]byte{0x0e, 0x71}, // starts with Se
	PKHSchnorrAddrID:     [2]byte{0x0e, 0x53}, // starts with SS
	PKHBlissAddrID:       [2]byte{0x0e, 0x68}, // starts with Sb
	PKHLmsAddrID:		  [2]byte{0x0e, 0x45}, // starts with SM
	ScriptHashAddrID:     [2]byte{0x0e, 0x6c}, // starts with Sc
	PrivateKeyID:         [2]byte{0x23, 0x07}, // starts with Ps

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x20, 0xb9, 0x03}, // starts with sprv
	HDPublicKeyID:  [4]byte{0x04, 0x20, 0xbd, 0x3d}, // starts with spub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 115, // ASCII for s

	// Hypercash PoS parameters
	MinimumStakeDiff:        20000,
	TicketPoolSize:          64,
	TicketsPerBlock:         5,
	TicketMaturity:          16,
	TicketExpiry:            384, // 6*TicketPoolSize
	CoinbaseMaturity:        16,
	SStxChangeMaturity:      1,
	TicketPoolSizeWeight:    4,
	StakeDiffAlpha:          1,
	StakeDiffWindowSize:     8,
	StakeDiffWindows:        8,
	StakeVersionInterval:    8 * 2 * 7,
	MaxFreshStakePerBlock:   20,            // 4*TicketsPerBlock
	StakeEnabledHeight:      16 + 16,       // CoinbaseMaturity + TicketMaturity
	StakeValidationHeight:   16 + (64 * 2), // CoinbaseMaturity + TicketPoolSize*2
	StakeBaseSigScript:      []byte{0xDE, 0xAD, 0xBE, 0xEF},
	StakeMajorityMultiplier: 3,
	StakeMajorityDivisor:    4,

	// Hypercash organization related parameters
	//
	// "Dev org" address is a 3-of-3 P2SH going to wallet:
	// aardvark adroitness aardvark adroitness
	// aardvark adroitness aardvark adroitness
	// aardvark adroitness aardvark adroitness
	// aardvark adroitness aardvark adroitness
	// aardvark adroitness aardvark adroitness
	// aardvark adroitness aardvark adroitness
	// aardvark adroitness aardvark adroitness
	// aardvark adroitness aardvark adroitness
	// briefcase
	// (seed 0x00000000000000000000000000000000000000000000000000000000000000)
	//
	// This same wallet owns the three ledger outputs for simnet.
	//
	// P2SH details for simnet dev org is below.
	//
	// address: Scc4ZC844nzuZCXsCFXUBXTLks2mD6psWom
	// redeemScript: 532103e8c60c7336744c8dcc7b85c27789950fc52aa4e48f895ebbfb
	// ac383ab893fc4c2103ff9afc246e0921e37d12e17d8296ca06a8f92a07fbe7857ed1d4
	// f0f5d94e988f21033ed09c7fa8b83ed53e6f2c57c5fa99ed2230c0d38edf53c0340d0f
	// c2e79c725a53ae
	//   (3-of-3 multisig)
	// Pubkeys used:
	//   SkQmxbeuEFDByPoTj41TtXat8tWySVuYUQpd4fuNNyUx51tF1csSs
	//   SkQn8ervNvAUEX5Ua3Lwjc6BAuTXRznDoDzsyxgjYqX58znY7w9e4
	//   SkQkfkHZeBbMW8129tZ3KspEh1XBFC1btbkgzs6cjSyPbrgxzsKqk
	//
	// Organization address is ScuQxvveKGfpG1ypt6u27F99Anf7EW3cqhq
	OrganizationPkScript:        hexDecode("a914cbb08d6ca783b533b2c7d24a51fbca92d937bf9987"),
	OrganizationPkScriptVersion: 0,
	BlockOneLedger:              BlockOneLedgerSimNet,
}

var (
	// ErrDuplicateNet describes an error where the parameters for a Hypercash
	// network could not be set due to the network already being a standard
	// network or previously-registered into this package.
	ErrDuplicateNet = errors.New("duplicate Hypercash network")

	// ErrUnknownHDKeyID describes an error where the provided id which
	// is intended to identify the network for a hierarchical deterministic
	// private extended key is not registered.
	ErrUnknownHDKeyID = errors.New("unknown hd private extended key bytes")
)

var (
	registeredNets    = make(map[wire.CurrencyNet]struct{})
	pubKeyAddrIDs     = make(map[[2]byte]struct{})
	pubKeyHashAddrIDs = make(map[[2]byte]struct{})
	pkhEdwardsAddrIDs = make(map[[2]byte]struct{})
	pkhSchnorrAddrIDs = make(map[[2]byte]struct{})
	scriptHashAddrIDs = make(map[[2]byte]struct{})
	hdPrivToPubKeyIDs = make(map[[4]byte][]byte)
)

// Register registers the network parameters for a Hypercash network.  This may
// error with ErrDuplicateNet if the network is already registered (either
// due to a previous Register call, or the network being one of the default
// networks).
//
// Network parameters should be registered into this package by a main package
// as early as possible.  Then, library packages may lookup networks or network
// parameters based on inputs and work regardless of the network being standard
// or not.
func Register(params *Params) error {
	if _, ok := registeredNets[params.Net]; ok {
		return ErrDuplicateNet
	}
	registeredNets[params.Net] = struct{}{}
	pubKeyAddrIDs[params.PubKeyAddrID] = struct{}{}
	pubKeyHashAddrIDs[params.PubKeyHashAddrID] = struct{}{}
	scriptHashAddrIDs[params.ScriptHashAddrID] = struct{}{}
	hdPrivToPubKeyIDs[params.HDPrivateKeyID] = params.HDPublicKeyID[:]
	return nil
}

// mustRegister performs the same function as Register except it panics if there
// is an error.  This should only be called from package init functions.
func mustRegister(params *Params) {
	if err := Register(params); err != nil {
		panic("failed to register network: " + err.Error())
	}
}

// IsPubKeyAddrID returns whether the id is an identifier known to prefix a
// pay-to-pubkey address on any default or registered network.
func IsPubKeyAddrID(id [2]byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}

// IsPubKeyHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-pubkey-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsScriptHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsPubKeyHashAddrID(id [2]byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}

// IsPKHEdwardsAddrID returns whether the id is an identifier know to prefix a
// pay-to-pubkey-hash Edwards address.
func IsPKHEdwardsAddrID(id [2]byte) bool {
	_, ok := pkhEdwardsAddrIDs[id]
	return ok
}

// IsPKHSchnorrAddrID returns whether the id is an identifier know to prefix a
// pay-to-pubkey-hash secp256k1 Schnorr address.
func IsPKHSchnorrAddrID(id [2]byte) bool {
	_, ok := pkhSchnorrAddrIDs[id]
	return ok
}

// IsScriptHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-script-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsPubKeyHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsScriptHashAddrID(id [2]byte) bool {
	_, ok := scriptHashAddrIDs[id]
	return ok
}

// HDPrivateKeyToPublicKeyID accepts a private hierarchical deterministic
// extended key id and returns the associated public key id.  When the provided
// id is not registered, the ErrUnknownHDKeyID error will be returned.
func HDPrivateKeyToPublicKeyID(id []byte) ([]byte, error) {
	if len(id) != 4 {
		return nil, ErrUnknownHDKeyID
	}

	var key [4]byte
	copy(key[:], id)
	pubBytes, ok := hdPrivToPubKeyIDs[key]
	if !ok {
		return nil, ErrUnknownHDKeyID
	}

	return pubBytes, nil
}

// newHashFromStr converts the passed big-endian hex string into a
// chainhash.Hash.  It only differs from the one available in chainhash in that
// it panics on an error since it will only (and must only) be called with
// hard-coded, and therefore known good, hashes.
func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		// Ordinarily I don't like panics in library code since it
		// can take applications down without them having a chance to
		// recover which is extremely annoying, however an exception is
		// being made in this case because the only way this can panic
		// is if there is an error in the hard-coded hashes.  Thus it
		// will only ever potentially panic on init and therefore is
		// 100% predictable.
		panic(err)
	}
	return hash
}

func hexDecode(hexStr string) []byte {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return b
}

// BlockOneSubsidy returns the total subsidy of block height 1 for the
// network.
func (p *Params) BlockOneSubsidy() int64 {
	if len(p.BlockOneLedger) == 0 {
		return 0
	}

	sum := int64(0)
	for _, output := range p.BlockOneLedger {
		sum += output.Amount
	}

	return sum
}

// TotalSubsidyProportions is the sum of WorkReward, StakeReward, and BlockTax
// proportions.
func (p *Params) TotalSubsidyProportions() uint16 {
	return p.WorkRewardProportion + p.StakeRewardProportion + p.BlockTaxProportion
}

// LatestCheckpointHeight is the height of the latest checkpoint block in the
// parameters.
func (p *Params) LatestCheckpointHeight() int64 {
	if len(p.Checkpoints) == 0 {
		return 0
	}
	return p.Checkpoints[len(p.Checkpoints)-1].Height
}

func init() {
	// Register all default networks when the package is initialized.
	mustRegister(&MainNetParams)
	mustRegister(&TestNet2Params)
	mustRegister(&SimNetParams)
}
