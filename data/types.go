package data

import (
  "math/big"

  "github.com/harmony-one/harmony/numeric"
  "github.com/harmony-one/harmony/staking/types"
)

type NodeMetadataReply struct {
  BLSKeys      []string `json:"blskey"`
	Version      string   `json:"version"`
	NetworkType  string   `json:"network"`
	ShardID      float64  `json:"shard-id"`
	NodeRole     string   `json:"role"`
	ArchivalNode bool     `json:"is-archival"`
}

type LatestHeaderReply struct {
	BlockHash        string  `json:"blockHash"`
	BlockNumber      float64 `json:"blockNumber"`
	ShardID          uint32  `json:"shardID"`
	Leader           string  `json:"leader"`
	ViewID           uint64  `json:"viewID"`
	Epoch            int     `json:"epoch"`
	Timestamp        string  `json:"timestamp"`
	UnixTime         int64   `json:"unixtime"`
	LastCommitSig    string  `json:"lastCommitSig"`
	LastCommitBitmap string  `json:"lastCommitBitmap"`
}

type BlockByNumberReply struct {
  BlockHash              string   `json:"hash"`
  BlockNumber            string   `json:"number"`
  BlockSize              string   `json:"size"`
  BlockSizeInt           int
  ParentHash             string   `json:"parentHash"`
  StateRoot              string   `json:"stateRoot"`
  Transactions           []string `json:"transactions"`
  NumTransactions        int
  StakingTransactions    []string `json:"stakingTransactions"`
  NumStakingTransactions int
}

// HACK: To get UnmarshalJSON to unmarshal one addresses
type ValidatorInformationReply struct {
  Wrapper              ValidatorWrapper               `json:"validator"`
	Performance          *types.CurrentEpochPerformance `json:"current-epoch-performance"`
	TotalDelegated       *big.Int                       `json:"total-delegation"`
	CurrentlyInCommittee bool                           `json:"currently-in-committee"`
	EPoSStatus           string                         `json:"epos-status"`
	EPoSWinningStake     *numeric.Dec                   `json:"epos-winning-stake"`
	BootedStatus         *string                        `json:"booted-status"`
	Lifetime             *types.AccumulatedOverLifetime `json:"lifetime"`
}

type ValidatorWrapper struct {
	Validator Validator
	Delegations types.Delegations
}

type Validator struct {
  Address              string `json:"address"`
  SlotPubKeys          []string `json:"bls-public-keys"`
  LastEpochInCommittee *big.Int `json:"last-epoch-in-committee"`
  MinSelfDelegation    *big.Int `json:"min-self-delegation"`
  MaxTotalDelegation   *big.Int `json:"max-total-delegation"`
  types.Commission
  types.Description
  CreationHeight       *big.Int `json:"creation-height"`
}

type StructureReply struct {
  ShardID uint32 `json:"shardID"`
  HTTP    string `json:"http"`
}
