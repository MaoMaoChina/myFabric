/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blockcutter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

<<<<<<< HEAD
	cb "github.com/hyperledger/fabric-protos-go/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/orderer/common/blockcutter"
	"github.com/hyperledger/fabric/orderer/common/blockcutter/mock"
)

var _ = Describe("Blockcutter", func() {
	var (
		bc                blockcutter.Receiver
		fakeConfig        *mock.OrdererConfig
		fakeConfigFetcher *mock.OrdererConfigFetcher

		metrics               *blockcutter.Metrics
		fakeBlockFillDuration *mock.MetricsHistogram
	)

	BeforeEach(func() {
		fakeConfig = &mock.OrdererConfig{}
		fakeConfigFetcher = &mock.OrdererConfigFetcher{}
		fakeConfigFetcher.OrdererConfigReturns(fakeConfig, true)

		fakeBlockFillDuration = &mock.MetricsHistogram{}
		fakeBlockFillDuration.WithReturns(fakeBlockFillDuration)
		metrics = &blockcutter.Metrics{
			BlockFillDuration: fakeBlockFillDuration,
		}

		bc = blockcutter.NewReceiverImpl("mychannel", fakeConfigFetcher, metrics)
	})

	Describe("Ordered", func() {
		var (
			message *cb.Envelope
		)

		BeforeEach(func() {
			fakeConfig.BatchSizeReturns(&ab.BatchSize{
				MaxMessageCount:   2,
				PreferredMaxBytes: 100,
			})

			message = &cb.Envelope{Payload: []byte("Twenty Bytes of Data"), Signature: []byte("Twenty Bytes of Data")}
		})

		It("adds the message to the pending batches", func() {
			batches, pending := bc.Ordered(message)
			Expect(batches).To(BeEmpty())
			Expect(pending).To(BeTrue())
			Expect(fakeBlockFillDuration.ObserveCallCount()).To(Equal(0))
		})

		Context("when enough batches to fill the max message count are enqueued", func() {
			It("cuts the batch", func() {
				batches, pending := bc.Ordered(message)
				Expect(batches).To(BeEmpty())
				Expect(pending).To(BeTrue())
				batches, pending = bc.Ordered(message)
				Expect(len(batches)).To(Equal(1))
				Expect(len(batches[0])).To(Equal(2))
				Expect(pending).To(BeFalse())

				Expect(fakeBlockFillDuration.ObserveCallCount()).To(Equal(1))
				Expect(fakeBlockFillDuration.ObserveArgsForCall(0)).To(BeNumerically(">", 0))
				Expect(fakeBlockFillDuration.ObserveArgsForCall(0)).To(BeNumerically("<", 1))
				Expect(fakeBlockFillDuration.WithCallCount()).To(Equal(1))
				Expect(fakeBlockFillDuration.WithArgsForCall(0)).To(Equal([]string{"channel", "mychannel"}))
			})
		})

		Context("when the message does not exceed max message count or preferred size", func() {
			BeforeEach(func() {
				fakeConfig.BatchSizeReturns(&ab.BatchSize{
					MaxMessageCount:   3,
					PreferredMaxBytes: 100,
				})
			})

			It("adds the message to the pending batches", func() {
				batches, pending := bc.Ordered(message)
				Expect(batches).To(BeEmpty())
				Expect(pending).To(BeTrue())
				batches, pending = bc.Ordered(message)
				Expect(batches).To(BeEmpty())
				Expect(pending).To(BeTrue())
				Expect(fakeBlockFillDuration.ObserveCallCount()).To(Equal(0))
			})
		})

		Context("when the message is larger than the preferred max bytes", func() {
			BeforeEach(func() {
				fakeConfig.BatchSizeReturns(&ab.BatchSize{
					MaxMessageCount:   3,
					PreferredMaxBytes: 30,
				})
			})

			It("cuts the batch immediately", func() {
				batches, pending := bc.Ordered(message)
				Expect(len(batches)).To(Equal(1))
				Expect(pending).To(BeFalse())
				Expect(fakeBlockFillDuration.ObserveCallCount()).To(Equal(1))
				Expect(fakeBlockFillDuration.ObserveArgsForCall(0)).To(Equal(float64(0)))
				Expect(fakeBlockFillDuration.WithCallCount()).To(Equal(1))
				Expect(fakeBlockFillDuration.WithArgsForCall(0)).To(Equal([]string{"channel", "mychannel"}))
			})
		})

		Context("when the message causes the batch to exceed the preferred max bytes", func() {
			BeforeEach(func() {
				fakeConfig.BatchSizeReturns(&ab.BatchSize{
					MaxMessageCount:   3,
					PreferredMaxBytes: 50,
				})
			})

			It("cuts the previous batch immediately, enqueueing the second", func() {
				batches, pending := bc.Ordered(message)
				Expect(batches).To(BeEmpty())
				Expect(pending).To(BeTrue())

				batches, pending = bc.Ordered(message)
				Expect(len(batches)).To(Equal(1))
				Expect(len(batches[0])).To(Equal(1))
				Expect(pending).To(BeTrue())

				Expect(fakeBlockFillDuration.ObserveCallCount()).To(Equal(1))
				Expect(fakeBlockFillDuration.ObserveArgsForCall(0)).To(BeNumerically(">", 0))
				Expect(fakeBlockFillDuration.ObserveArgsForCall(0)).To(BeNumerically("<", 1))
				Expect(fakeBlockFillDuration.WithCallCount()).To(Equal(1))
				Expect(fakeBlockFillDuration.WithArgsForCall(0)).To(Equal([]string{"channel", "mychannel"}))
			})

			Context("when the new message is larger than the preferred max bytes", func() {
				var (
					bigMessage *cb.Envelope
				)

				BeforeEach(func() {
					bigMessage = &cb.Envelope{Payload: make([]byte, 1000)}
				})

				It("cuts both the previous batch and the next batch immediately", func() {
					batches, pending := bc.Ordered(message)
					Expect(batches).To(BeEmpty())
					Expect(pending).To(BeTrue())

					batches, pending = bc.Ordered(bigMessage)
					Expect(len(batches)).To(Equal(2))
					Expect(len(batches[0])).To(Equal(1))
					Expect(len(batches[1])).To(Equal(1))
					Expect(pending).To(BeFalse())

					Expect(fakeBlockFillDuration.ObserveCallCount()).To(Equal(2))
					Expect(fakeBlockFillDuration.ObserveArgsForCall(0)).To(BeNumerically(">", 0))
					Expect(fakeBlockFillDuration.ObserveArgsForCall(0)).To(BeNumerically("<", 1))
					Expect(fakeBlockFillDuration.ObserveArgsForCall(1)).To(Equal(float64(0)))
					Expect(fakeBlockFillDuration.WithCallCount()).To(Equal(2))
					Expect(fakeBlockFillDuration.WithArgsForCall(0)).To(Equal([]string{"channel", "mychannel"}))
					Expect(fakeBlockFillDuration.WithArgsForCall(1)).To(Equal([]string{"channel", "mychannel"}))
				})
			})
		})

		Context("when the orderer config cannot be retrieved", func() {
			BeforeEach(func() {
				fakeConfigFetcher.OrdererConfigReturns(nil, false)
			})

			It("panics", func() {
				Expect(func() { bc.Ordered(message) }).To(Panic())
			})
		})
	})

	Describe("Cut", func() {
		It("cuts an empty batch", func() {
			batch := bc.Cut()
			Expect(batch).To(BeNil())
			Expect(fakeBlockFillDuration.ObserveCallCount()).To(Equal(0))
		})
	})
})
=======
	mockconfig "github.com/hyperledger/fabric/common/mocks/config"
	"github.com/hyperledger/fabric/orderer/common/filter"
	cb "github.com/hyperledger/fabric/protos/common"
	ab "github.com/hyperledger/fabric/protos/orderer"
	logging "github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.SetLevel(logging.DEBUG, "")
}

type isolatedCommitter struct{}

func (ic isolatedCommitter) Isolated() bool { return true }

func (ic isolatedCommitter) Commit() {}

type mockIsolatedFilter struct{}

func (mif *mockIsolatedFilter) Apply(msg *cb.Envelope) (filter.Action, filter.Committer) {
	if bytes.Equal(msg.Payload, isolatedTx.Payload) {
		return filter.Accept, isolatedCommitter{}
	}
	return filter.Forward, nil
}

type mockRejectFilter struct{}

func (mrf mockRejectFilter) Apply(message *cb.Envelope) (filter.Action, filter.Committer) {
	if bytes.Equal(message.Payload, badTx.Payload) {
		return filter.Reject, nil
	}
	return filter.Forward, nil
}

type mockAcceptFilter struct{}

func (mrf mockAcceptFilter) Apply(message *cb.Envelope) (filter.Action, filter.Committer) {
	if bytes.Equal(message.Payload, goodTx.Payload) {
		return filter.Accept, filter.NoopCommitter
	}
	return filter.Forward, nil
}

func getFilters() *filter.RuleSet {
	return filter.NewRuleSet([]filter.Rule{
		&mockIsolatedFilter{},
		&mockRejectFilter{},
		&mockAcceptFilter{},
	})
}

var badTx = &cb.Envelope{Payload: []byte("BAD")}
var goodTx = &cb.Envelope{Payload: []byte("GOOD")}
var goodTxLarge = &cb.Envelope{Payload: []byte("GOOD"), Signature: make([]byte, 1000)}
var isolatedTx = &cb.Envelope{Payload: []byte("ISOLATED")}
var unmatchedTx = &cb.Envelope{Payload: []byte("UNMATCHED")}

func TestNormalBatch(t *testing.T) {
	filters := getFilters()
	maxMessageCount := uint32(2)
	absoluteMaxBytes := uint32(1000)
	preferredMaxBytes := uint32(100)
	r := NewReceiverImpl(&mockconfig.Orderer{BatchSizeVal: &ab.BatchSize{MaxMessageCount: maxMessageCount, AbsoluteMaxBytes: absoluteMaxBytes, PreferredMaxBytes: preferredMaxBytes}}, filters)

	batches, committers, ok, pending := r.Ordered(goodTx)

	assert.Nil(t, batches, "Should not have created batch")
	assert.Nil(t, committers, "Should not have created batch")
	assert.True(t, ok, "Should have enqueued message into batch")
	assert.True(t, pending, "Should have pending messages")

	batches, committers, ok, pending = r.Ordered(goodTx)

	assert.Len(t, batches, 1, "Should have created 1 message batch, got %d", len(batches))
	assert.Len(t, committers, 1, "Should have created 1 committer batch, got %d", len(committers))
	assert.True(t, ok, "Should have enqueued message into batch")
	assert.False(t, pending, "Should not have pending messages")
}

func TestBadMessageInBatch(t *testing.T) {
	filters := getFilters()
	maxMessageCount := uint32(2)
	absoluteMaxBytes := uint32(1000)
	preferredMaxBytes := uint32(100)
	r := NewReceiverImpl(&mockconfig.Orderer{BatchSizeVal: &ab.BatchSize{MaxMessageCount: maxMessageCount, AbsoluteMaxBytes: absoluteMaxBytes, PreferredMaxBytes: preferredMaxBytes}}, filters)

	batches, committers, ok, _ := r.Ordered(badTx)

	assert.Nil(t, batches, "Should not have created batch")
	assert.Nil(t, committers, "Should not have created batch")
	assert.False(t, ok, "Should not have enqueued bad message into batch")

	batches, committers, ok, pending := r.Ordered(goodTx)

	assert.Nil(t, batches, "Should not have created batch")
	assert.Nil(t, committers, "Should not have created batch")
	assert.True(t, ok, "Should have enqueued good message into batch")
	assert.True(t, pending, "Should have pending messages")

	batches, committers, ok, _ = r.Ordered(badTx)

	assert.Nil(t, batches, "Should not have created batch")
	assert.Nil(t, committers, "Should not have created batch")
	assert.False(t, ok, "Should not have enqueued second bad message into batch")
}

func TestUnmatchedMessageInBatch(t *testing.T) {
	filters := getFilters()
	maxMessageCount := uint32(2)
	absoluteMaxBytes := uint32(1000)
	preferredMaxBytes := uint32(100)
	r := NewReceiverImpl(&mockconfig.Orderer{BatchSizeVal: &ab.BatchSize{MaxMessageCount: maxMessageCount, AbsoluteMaxBytes: absoluteMaxBytes, PreferredMaxBytes: preferredMaxBytes}}, filters)

	batches, committers, ok, _ := r.Ordered(unmatchedTx)

	assert.Nil(t, batches, "Should not have created batch")
	assert.Nil(t, committers, "Should not have created batch")
	assert.False(t, ok, "Should not have enqueued unmatched message into batch")

	batches, committers, ok, pending := r.Ordered(goodTx)

	assert.Nil(t, batches, "Should not have created batch")
	assert.Nil(t, committers, "Should not have created batch")
	assert.True(t, ok, "Should have enqueued good message into batch")
	assert.True(t, pending, "Should have pending messages")

	batches, committers, ok, _ = r.Ordered(unmatchedTx)

	assert.Nil(t, batches, "Should not have created batch from unmatched message")
	assert.Nil(t, committers, "Should not have created batch")
	assert.False(t, ok, "Should not have enqueued second unmatched message into batch")
}

func TestIsolatedEmptyBatch(t *testing.T) {
	filters := getFilters()
	maxMessageCount := uint32(2)
	absoluteMaxBytes := uint32(1000)
	preferredMaxBytes := uint32(100)
	r := NewReceiverImpl(&mockconfig.Orderer{BatchSizeVal: &ab.BatchSize{MaxMessageCount: maxMessageCount, AbsoluteMaxBytes: absoluteMaxBytes, PreferredMaxBytes: preferredMaxBytes}}, filters)

	batches, committers, ok, pending := r.Ordered(isolatedTx)

	assert.Len(t, batches, 1, "Should created 1 new message batch, got %d", len(batches))
	assert.Len(t, batches[0], 1, "Should have had one isolatedTx in the message batch, got %d", len(batches[0]))
	assert.Len(t, committers, 1, "Should created 1 new committer batch, got %d", len(committers))
	assert.Len(t, committers[0], 1, "Should have had one isolatedTx in the committer batch, got %d", len(committers[0]))
	assert.True(t, ok, "Should have enqueued isolated message into batch")
	assert.False(t, pending, "Should not have pending messages")
	assert.Equal(t, isolatedTx.Payload, batches[0][0].Payload, "Should have had the isolated tx in the first batch")
}

func TestIsolatedPartialBatch(t *testing.T) {
	filters := getFilters()
	maxMessageCount := uint32(2)
	absoluteMaxBytes := uint32(1000)
	preferredMaxBytes := uint32(100)
	r := NewReceiverImpl(&mockconfig.Orderer{BatchSizeVal: &ab.BatchSize{MaxMessageCount: maxMessageCount, AbsoluteMaxBytes: absoluteMaxBytes, PreferredMaxBytes: preferredMaxBytes}}, filters)

	batches, committers, ok, pending := r.Ordered(goodTx)

	assert.Nil(t, batches, "Should not have created batch")
	assert.Nil(t, committers, "Should not have created batch")
	assert.True(t, ok, "Should have enqueued message into batch")
	assert.True(t, pending, "Should have pending messages")

	batches, committers, ok, pending = r.Ordered(isolatedTx)

	assert.Len(t, batches, 2, "Should created 2 new message batch, got %d", len(batches))
	assert.Len(t, batches[0], 1, "Should have had one goodTx in the first message batch, got %d", len(batches[0]))
	assert.Len(t, batches[1], 1, "Should have had one isolatedTx in the second message batch, got %d", len(batches[1]))
	assert.Len(t, committers, 2, "Should created 2 new committer batch, got %d", len(committers))
	assert.Len(t, committers[0], 1, "Should have had 1 committer in the first committer batch, got %d", len(committers[0]))
	assert.Len(t, committers[1], 1, "Should have had 1 committer in the second committer batch, got %d", len(committers[1]))
	assert.True(t, ok, "Should have enqueued isolated message into batch")
	assert.False(t, pending, "Should not have pending messages")
	assert.Equal(t, goodTx.Payload, batches[0][0].Payload, "Should have had the good tx in the first batch")
	assert.Equal(t, isolatedTx.Payload, batches[1][0].Payload, "Should have had the isolated tx in the second batch")
}

func TestBatchSizePreferredMaxBytesOverflow(t *testing.T) {
	filters := getFilters()

	goodTxBytes := messageSizeBytes(goodTx)

	// set preferred max bytes such that 10 goodTx will not fit
	preferredMaxBytes := goodTxBytes*10 - 1

	// set message count > 9
	maxMessageCount := uint32(20)

	r := NewReceiverImpl(&mockconfig.Orderer{BatchSizeVal: &ab.BatchSize{MaxMessageCount: maxMessageCount, AbsoluteMaxBytes: preferredMaxBytes * 2, PreferredMaxBytes: preferredMaxBytes}}, filters)

	// enqueue 9 messages
	for i := 0; i < 9; i++ {
		batches, committers, ok, pending := r.Ordered(goodTx)

		assert.Nil(t, batches, "Should not have created batch")
		assert.Nil(t, committers, "Should not have created batch")
		assert.True(t, ok, "Should have enqueued message into batch")
		assert.True(t, pending, "Should have pending messages")
	}

	// next message should create batch
	batches, committers, ok, pending := r.Ordered(goodTx)

	assert.Len(t, batches, 1, "Should have created 1 message batch, got %d", len(batches))
	assert.Len(t, batches[0], 9, "Should have had nine normal tx in the message batch, got %d", len(batches[0]))
	assert.Len(t, committers, 1, "Should have created 1 committer batch, got %d", len(committers))
	assert.Len(t, committers[0], 9, "Should have had nine committers in the committer batch, got %d", len(committers[0]))
	assert.True(t, ok, "Should have enqueued message into batch")
	assert.True(t, pending, "Should still have pending messages")

	// force a batch cut
	messageBatch, committerBatch := r.Cut()

	assert.NotNil(t, messageBatch, "Should have created message batch")
	assert.Len(t, messageBatch, 1, "Should have had 1 tx in the batch, got %d", len(messageBatch))
	assert.NotNil(t, committerBatch, "Should have created committer batch")
	assert.Len(t, committerBatch, 1, "Should have had 1 committer in the committer batch, got %d", len(committerBatch))
}

func TestBatchSizePreferredMaxBytesOverflowNoPending(t *testing.T) {
	filters := getFilters()

	goodTxLargeBytes := messageSizeBytes(goodTxLarge)

	// set preferred max bytes such that 1 goodTxLarge will not fit
	preferredMaxBytes := goodTxLargeBytes - 1

	// set message count > 1
	maxMessageCount := uint32(20)

	r := NewReceiverImpl(&mockconfig.Orderer{BatchSizeVal: &ab.BatchSize{MaxMessageCount: maxMessageCount, AbsoluteMaxBytes: preferredMaxBytes * 3, PreferredMaxBytes: preferredMaxBytes}}, filters)

	// submit large message
	batches, committers, ok, pending := r.Ordered(goodTxLarge)

	assert.Len(t, batches, 1, "Should have created 1 message batch, got %d", len(batches))
	assert.Len(t, batches[0], 1, "Should have had 1 normal tx in the message batch, got %d", len(batches[0]))
	assert.Len(t, committers, 1, "Should have created 1 committer batch, got %d", len(committers))
	assert.Len(t, committers[0], 1, "Should have had 1 committer in the committer batch, got %d", len(committers[0]))
	assert.True(t, ok, "Should have enqueued message into batch")
	assert.False(t, pending, "Should not have pending messages")
}
>>>>>>> release-1.0
