/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blockcutter

import (
	"time"

	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/hyperledger/fabric/common/flogging"
)

var logger = flogging.MustGetLogger("orderer.common.blockcutter")

type OrdererConfigFetcher interface {
	OrdererConfig() (channelconfig.Orderer, bool)
}

// Receiver defines a sink for the ordered broadcast messages
type Receiver interface {
	// Ordered should be invoked sequentially as messages are ordered
<<<<<<< HEAD
	// Each batch in `messageBatches` will be wrapped into a block.
	// `pending` indicates if there are still messages pending in the receiver.
	Ordered(msg *cb.Envelope) (messageBatches [][]*cb.Envelope, pending bool)
=======
	// If the current message valid, and no batches need to be cut:
	//   - Ordered will return nil, nil, and true (indicating valid Tx).
	// If the current message valid, and batches need to be cut:
	//   - Ordered will return 1 or 2 batches of messages, 1 or 2 batches of committers, and true (indicating valid Tx).
	// If the current message is invalid:
	//   - Ordered will return nil, nil, and false (to indicate invalid Tx).
	//
	// Given a valid message, if the current message needs to be isolated (as determined during filtering).
	//   - Ordered will return:
	//     * The pending batch of (if not empty), and a second batch containing only the isolated message.
	//     * The corresponding batches of committers.
	//     * true (indicating ok).
	// Otherwise, given a valid message, the pending batch, if not empty, will be cut and returned if:
	//   - The current message needs to be isolated (as determined during filtering).
	//   - The current message will cause the pending batch size in bytes to exceed BatchSize.PreferredMaxBytes.
	//   - After adding the current message to the pending batch, the message count has reached BatchSize.MaxMessageCount.
	//
	// In any case, `pending` is set to true if there are still messages pending in the receiver after cutting the block.
	Ordered(msg *cb.Envelope) (messageBatches [][]*cb.Envelope, committers [][]filter.Committer, validTx bool, pending bool)
>>>>>>> release-1.0

	// Cut returns the current batch and starts a new one
	Cut() []*cb.Envelope
}

type receiver struct {
	sharedConfigFetcher   OrdererConfigFetcher
	pendingBatch          []*cb.Envelope
	pendingBatchSizeBytes uint32

	PendingBatchStartTime time.Time
	ChannelID             string
	Metrics               *Metrics
}

// NewReceiverImpl creates a Receiver implementation based on the given configtxorderer manager
func NewReceiverImpl(channelID string, sharedConfigFetcher OrdererConfigFetcher, metrics *Metrics) Receiver {
	return &receiver{
		sharedConfigFetcher: sharedConfigFetcher,
		Metrics:             metrics,
		ChannelID:           channelID,
	}
}

// Ordered should be invoked sequentially as messages are ordered
<<<<<<< HEAD
//
// messageBatches length: 0, pending: false
//   - impossible, as we have just received a message
// messageBatches length: 0, pending: true
//   - no batch is cut and there are messages pending
// messageBatches length: 1, pending: false
//   - the message count reaches BatchSize.MaxMessageCount
// messageBatches length: 1, pending: true
//   - the current message will cause the pending batch size in bytes to exceed BatchSize.PreferredMaxBytes.
// messageBatches length: 2, pending: false
//   - the current message size in bytes exceeds BatchSize.PreferredMaxBytes, therefore isolated in its own batch.
// messageBatches length: 2, pending: true
//   - impossible
//
// Note that messageBatches can not be greater than 2.
func (r *receiver) Ordered(msg *cb.Envelope) (messageBatches [][]*cb.Envelope, pending bool) {
	if len(r.pendingBatch) == 0 {
		// We are beginning a new batch, mark the time
		r.PendingBatchStartTime = time.Now()
	}

	ordererConfig, ok := r.sharedConfigFetcher.OrdererConfig()
	if !ok {
		logger.Panicf("Could not retrieve orderer config to query batch parameters, block cutting is not possible")
	}
=======
// If the current message valid, and no batches need to be cut:
//   - Ordered will return nil, nil, true (indicating valid tx) and true (indicating there are pending messages).
// If the current message valid, and batches need to be cut:
//   - Ordered will return 1 or 2 batches of messages, 1 or 2 batches of committers, and true (indicating valid tx).
// If the current message is invalid:
//   - Ordered will return nil, nil, and false (to indicate invalid tx).
//
// Given a valid message, if the current message needs to be isolated (as determined during filtering).
//   - Ordered will return:
//     * The pending batch of (if not empty), and a second batch containing only the isolated message.
//     * The corresponding batches of committers.
//     * true (indicating valid tx).
// Otherwise, given a valid message, the pending batch, if not empty, will be cut and returned if:
//   - The current message needs to be isolated (as determined during filtering).
//   - The current message will cause the pending batch size in bytes to exceed BatchSize.PreferredMaxBytes.
//   - After adding the current message to the pending batch, the message count has reached BatchSize.MaxMessageCount.
//
// In any case, `pending` is set to true if there are still messages pending in the receiver after cutting the block.
func (r *receiver) Ordered(msg *cb.Envelope) (messageBatches [][]*cb.Envelope, committerBatches [][]filter.Committer, validTx bool, pending bool) {
	// The messages must be filtered a second time in case configuration has changed since the message was received
	committer, err := r.filters.Apply(msg)
	if err != nil {
		logger.Debugf("Rejecting message: %s", err)
		return // We don't bother to determine `pending` here as it's not processed in error case
	}

	// message is valid
	validTx = true

	messageSizeBytes := messageSizeBytes(msg)
>>>>>>> release-1.0

	batchSize := ordererConfig.BatchSize()

<<<<<<< HEAD
	messageSizeBytes := messageSizeBytes(msg)
	if messageSizeBytes > batchSize.PreferredMaxBytes {
		logger.Debugf("The current message, with %v bytes, is larger than the preferred batch size of %v bytes and will be isolated.", messageSizeBytes, batchSize.PreferredMaxBytes)
=======
		if committer.Isolated() {
			logger.Debugf("Found message which requested to be isolated, cutting into its own batch")
		} else {
			logger.Debugf("The current message, with %v bytes, is larger than the preferred batch size of %v bytes and will be isolated.", messageSizeBytes, r.sharedConfigManager.BatchSize().PreferredMaxBytes)
		}
>>>>>>> release-1.0

		// cut pending batch, if it has any messages
		if len(r.pendingBatch) > 0 {
			messageBatch := r.Cut()
			messageBatches = append(messageBatches, messageBatch)
		}

		// create new batch with single message
		messageBatches = append(messageBatches, []*cb.Envelope{msg})

<<<<<<< HEAD
		// Record that this batch took no time to fill
		r.Metrics.BlockFillDuration.With("channel", r.ChannelID).Observe(0)

		return
	}

	messageWillOverflowBatchSizeBytes := r.pendingBatchSizeBytes+messageSizeBytes > batchSize.PreferredMaxBytes
=======
		return
	}

	messageWillOverflowBatchSizeBytes := r.pendingBatchSizeBytes+messageSizeBytes > r.sharedConfigManager.BatchSize().PreferredMaxBytes
>>>>>>> release-1.0

	if messageWillOverflowBatchSizeBytes {
		logger.Debugf("The current message, with %v bytes, will overflow the pending batch of %v bytes.", messageSizeBytes, r.pendingBatchSizeBytes)
		logger.Debugf("Pending batch would overflow if current message is added, cutting batch now.")
		messageBatch := r.Cut()
		r.PendingBatchStartTime = time.Now()
		messageBatches = append(messageBatches, messageBatch)
	}

	logger.Debugf("Enqueuing message into batch")
	r.pendingBatch = append(r.pendingBatch, msg)
	r.pendingBatchSizeBytes += messageSizeBytes
<<<<<<< HEAD
=======
	r.pendingCommitters = append(r.pendingCommitters, committer)
>>>>>>> release-1.0
	pending = true

	if uint32(len(r.pendingBatch)) >= batchSize.MaxMessageCount {
		logger.Debugf("Batch size met, cutting batch")
		messageBatch := r.Cut()
		messageBatches = append(messageBatches, messageBatch)
<<<<<<< HEAD
=======
		committerBatches = append(committerBatches, committerBatch)
>>>>>>> release-1.0
		pending = false
	}

	return
}

// Cut returns the current batch and starts a new one
func (r *receiver) Cut() []*cb.Envelope {
	if r.pendingBatch != nil {
		r.Metrics.BlockFillDuration.With("channel", r.ChannelID).Observe(time.Since(r.PendingBatchStartTime).Seconds())
	}
	r.PendingBatchStartTime = time.Time{}
	batch := r.pendingBatch
	r.pendingBatch = nil
	r.pendingBatchSizeBytes = 0
	return batch
}

func messageSizeBytes(message *cb.Envelope) uint32 {
	return uint32(len(message.Payload) + len(message.Signature))
}
