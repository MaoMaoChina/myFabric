/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package config

import (
	cb "github.com/hyperledger/fabric-protos-go/common"
)

// Config encapsulates config (channel or resource) tree
type Config interface {
	// ConfigProto returns the current config
	ConfigProto() *cb.Config

	// ProposeConfigUpdate attempts to validate a new configtx against the current config state
	ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error)
}

<<<<<<< HEAD
// Manager provides access to the resource config
type Manager interface {
	// GetChannelConfig defines methods that are related to channel configuration
	GetChannelConfig(channel string) Config
=======
// Application stores the common shared application config
type Application interface {
	// Organizations returns a map of org ID to ApplicationOrg
	Organizations() map[string]ApplicationOrg
}

// Channel gives read only access to the channel configuration
type Channel interface {
	// HashingAlgorithm returns the default algorithm to be used when hashing
	// such as computing block hashes, and CreationPolicy digests
	HashingAlgorithm() func(input []byte) []byte

	// BlockDataHashingStructureWidth returns the width to use when constructing the
	// Merkle tree to compute the BlockData hash
	BlockDataHashingStructureWidth() uint32

	// OrdererAddresses returns the list of valid orderer addresses to connect to to invoke Broadcast/Deliver
	OrdererAddresses() []string
}

// Consortiums represents the set of consortiums serviced by an ordering service
type Consortiums interface {
	// Consortiums returns the set of consortiums
	Consortiums() map[string]Consortium
}

// Consortium represents a group of orgs which may create channels together
type Consortium interface {
	// ChannelCreationPolicy returns the policy to check when instantiating a channel for this consortium
	ChannelCreationPolicy() *cb.Policy
}

// Orderer stores the common shared orderer config
type Orderer interface {
	// ConsensusType returns the configured consensus type
	// 共识协议 Solo/Kafka
	ConsensusType() string

	// BatchSize returns the maximum number of messages to include in a block
	BatchSize() *ab.BatchSize

	// BatchTimeout returns the amount of time to wait before creating a batch
	BatchTimeout() time.Duration

	// MaxChannelsCount returns the maximum count of channels to allow for an ordering network
	MaxChannelsCount() uint64

	// KafkaBrokers returns the addresses (IP:port notation) of a set of "bootstrap"
	// Kafka brokers, i.e. this is not necessarily the entire set of Kafka brokers
	// used for ordering
	KafkaBrokers() []string

	// Organizations returns the organizations for the ordering service
	Organizations() map[string]Org
}

type ValueProposer interface {
	// BeginValueProposals called when a config proposal is begun
	BeginValueProposals(tx interface{}, groups []string) (ValueDeserializer, []ValueProposer, error)

	// RollbackProposals called when a config proposal is abandoned
	RollbackProposals(tx interface{})

	// PreCommit is invoked before committing the config to catch
	// any errors which cannot be caught on a per proposal basis
	// TODO, rename other methods to remove Value/Proposal references
	PreCommit(tx interface{}) error

	// CommitProposals called when a config proposal is committed
	CommitProposals(tx interface{})
>>>>>>> release-1.0
}
