/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser_test

import (
<<<<<<< HEAD
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	mspproto "github.com/hyperledger/fabric-protos-go/msp"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/metrics/metricsfakes"
	"github.com/hyperledger/fabric/core/chaincode/lifecycle"
	"github.com/hyperledger/fabric/core/endorser"
	"github.com/hyperledger/fabric/core/endorser/fake"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/protoutil"

	"github.com/golang/protobuf/proto"
=======
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	mockpolicies "github.com/hyperledger/fabric/common/mocks/policies"
	"github.com/hyperledger/fabric/common/policies"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode"
	"github.com/hyperledger/fabric/core/common/ccprovider"
	"github.com/hyperledger/fabric/core/config"
	"github.com/hyperledger/fabric/core/container"
	"github.com/hyperledger/fabric/core/peer"
	syscc "github.com/hyperledger/fabric/core/scc"
	"github.com/hyperledger/fabric/core/testutil"
	"github.com/hyperledger/fabric/msp"
	mspmgmt "github.com/hyperledger/fabric/msp/mgmt"
	"github.com/hyperledger/fabric/msp/mgmt/testtools"
	"github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	pbutils "github.com/hyperledger/fabric/protos/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
>>>>>>> release-1.0
)

var _ = Describe("Endorser", func() {
	var (
		fakeProposalDuration         *metricsfakes.Histogram
		fakeProposalsReceived        *metricsfakes.Counter
		fakeSuccessfulProposals      *metricsfakes.Counter
		fakeProposalValidationFailed *metricsfakes.Counter
		fakeProposalACLCheckFailed   *metricsfakes.Counter
		fakeInitFailed               *metricsfakes.Counter
		fakeEndorsementsFailed       *metricsfakes.Counter
		fakeDuplicateTxsFailure      *metricsfakes.Counter

		fakeLocalIdentity                *fake.Identity
		fakeLocalMSPIdentityDeserializer *fake.IdentityDeserializer

		fakeChannelIdentity                *fake.Identity
		fakeChannelMSPIdentityDeserializer *fake.IdentityDeserializer

		fakeChannelFetcher *fake.ChannelFetcher

		fakePrivateDataDistributor *fake.PrivateDataDistributor

		fakeSupport              *fake.Support
		fakeTxSimulator          *fake.TxSimulator
		fakeHistoryQueryExecutor *fake.HistoryQueryExecutor

		signedProposal *pb.SignedProposal
		channelID      string
		chaincodeName  string

		chaincodeResponse *pb.Response
		chaincodeEvent    *pb.ChaincodeEvent
		chaincodeInput    *pb.ChaincodeInput

		e *endorser.Endorser
	)

	BeforeEach(func() {
		fakeProposalDuration = &metricsfakes.Histogram{}
		fakeProposalDuration.WithReturns(fakeProposalDuration)

		fakeProposalACLCheckFailed = &metricsfakes.Counter{}
		fakeProposalACLCheckFailed.WithReturns(fakeProposalACLCheckFailed)

		fakeInitFailed = &metricsfakes.Counter{}
		fakeInitFailed.WithReturns(fakeInitFailed)

		fakeEndorsementsFailed = &metricsfakes.Counter{}
		fakeEndorsementsFailed.WithReturns(fakeEndorsementsFailed)

		fakeDuplicateTxsFailure = &metricsfakes.Counter{}
		fakeDuplicateTxsFailure.WithReturns(fakeDuplicateTxsFailure)

		fakeProposalsReceived = &metricsfakes.Counter{}
		fakeSuccessfulProposals = &metricsfakes.Counter{}
		fakeProposalValidationFailed = &metricsfakes.Counter{}

		fakeLocalIdentity = &fake.Identity{}
		fakeLocalMSPIdentityDeserializer = &fake.IdentityDeserializer{}
		fakeLocalMSPIdentityDeserializer.DeserializeIdentityReturns(fakeLocalIdentity, nil)

		fakeChannelIdentity = &fake.Identity{}
		fakeChannelMSPIdentityDeserializer = &fake.IdentityDeserializer{}
		fakeChannelMSPIdentityDeserializer.DeserializeIdentityReturns(fakeChannelIdentity, nil)

		fakeChannelFetcher = &fake.ChannelFetcher{}
		fakeChannelFetcher.ChannelReturns(&endorser.Channel{
			IdentityDeserializer: fakeChannelMSPIdentityDeserializer,
		})

		fakePrivateDataDistributor = &fake.PrivateDataDistributor{}

		channelID = "channel-id"
		chaincodeName = "chaincode-name"
		chaincodeInput = &pb.ChaincodeInput{
			Args: [][]byte{[]byte("arg1"), []byte("arg2"), []byte("arg3")},
		}

		chaincodeResponse = &pb.Response{
			Status:  200,
			Payload: []byte("response-payload"),
		}
		chaincodeEvent = &pb.ChaincodeEvent{
			ChaincodeId: "chaincode-id",
			TxId:        "event-txid",
			EventName:   "event-name",
			Payload:     []byte("event-payload"),
		}

		fakeSupport = &fake.Support{}
		fakeSupport.ExecuteReturns(
			chaincodeResponse,
			chaincodeEvent,
			nil,
		)

		fakeSupport.ChaincodeEndorsementInfoReturns(&lifecycle.ChaincodeEndorsementInfo{
			Version:           "chaincode-definition-version",
			EndorsementPlugin: "plugin-name",
		}, nil)

		fakeSupport.GetLedgerHeightReturns(7, nil)

		fakeSupport.EndorseWithPluginReturns(
			&pb.Endorsement{
				Endorser:  []byte("endorser-identity"),
				Signature: []byte("endorser-signature"),
			},
			[]byte("endorser-modified-payload"),
			nil,
		)

		fakeTxSimulator = &fake.TxSimulator{}
		fakeTxSimulator.GetTxSimulationResultsReturns(
			&ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
				PvtSimulationResults: &rwset.TxPvtReadWriteSet{},
			},
			nil,
		)

		fakeHistoryQueryExecutor = &fake.HistoryQueryExecutor{}
		fakeSupport.GetHistoryQueryExecutorReturns(fakeHistoryQueryExecutor, nil)

		fakeSupport.GetTxSimulatorReturns(fakeTxSimulator, nil)

		fakeSupport.GetTransactionByIDReturns(nil, fmt.Errorf("txid-error"))

		e = &endorser.Endorser{
			LocalMSP:               fakeLocalMSPIdentityDeserializer,
			PrivateDataDistributor: fakePrivateDataDistributor,
			Metrics: &endorser.Metrics{
				ProposalDuration:         fakeProposalDuration,
				ProposalsReceived:        fakeProposalsReceived,
				SuccessfulProposals:      fakeSuccessfulProposals,
				ProposalValidationFailed: fakeProposalValidationFailed,
				ProposalACLCheckFailed:   fakeProposalACLCheckFailed,
				InitFailed:               fakeInitFailed,
				EndorsementsFailed:       fakeEndorsementsFailed,
				DuplicateTxsFailure:      fakeDuplicateTxsFailure,
			},
			Support:        fakeSupport,
			ChannelFetcher: fakeChannelFetcher,
		}
	})

	JustBeforeEach(func() {
		signedProposal = &pb.SignedProposal{
			ProposalBytes: protoutil.MarshalOrPanic(&pb.Proposal{
				Header: protoutil.MarshalOrPanic(&cb.Header{
					ChannelHeader: protoutil.MarshalOrPanic(&cb.ChannelHeader{
						Type:      int32(cb.HeaderType_ENDORSER_TRANSACTION),
						ChannelId: channelID,
						Extension: protoutil.MarshalOrPanic(&pb.ChaincodeHeaderExtension{
							ChaincodeId: &pb.ChaincodeID{
								Name: chaincodeName,
							},
						}),
						TxId: "6f142589e4ef6a1e62c9c816e2074f70baa9f7cf67c2f0c287d4ef907d6d2015",
					}),
					SignatureHeader: protoutil.MarshalOrPanic(&cb.SignatureHeader{
						Creator: protoutil.MarshalOrPanic(&mspproto.SerializedIdentity{
							Mspid: "msp-id",
						}),
						Nonce: []byte("nonce"),
					}),
				}),
				Payload: protoutil.MarshalOrPanic(&pb.ChaincodeProposalPayload{
					Input: protoutil.MarshalOrPanic(&pb.ChaincodeInvocationSpec{
						ChaincodeSpec: &pb.ChaincodeSpec{
							Input: chaincodeInput,
						},
					}),
				}),
			}),
			Signature: []byte("signature"),
		}
<<<<<<< HEAD
	})

	It("successfully endorses the proposal", func() {
		proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(proposalResponse.Endorsement).To(Equal(&pb.Endorsement{
			Endorser:  []byte("endorser-identity"),
			Signature: []byte("endorser-signature"),
		}))
		Expect(proposalResponse.Timestamp).To(BeNil())
		Expect(proposalResponse.Version).To(Equal(int32(1)))
		Expect(proposalResponse.Payload).To(Equal([]byte("endorser-modified-payload")))
		Expect(proto.Equal(proposalResponse.Response, &pb.Response{
			Status:  200,
			Payload: []byte("response-payload"),
		})).To(BeTrue())

		Expect(fakeSupport.EndorseWithPluginCallCount()).To(Equal(1))
		pluginName, cid, propRespPayloadBytes, sp := fakeSupport.EndorseWithPluginArgsForCall(0)
		Expect(sp).To(Equal(signedProposal))
		Expect(pluginName).To(Equal("plugin-name"))
		Expect(cid).To(Equal("channel-id"))

		prp := &pb.ProposalResponsePayload{}
		err = proto.Unmarshal(propRespPayloadBytes, prp)
		Expect(err).NotTo(HaveOccurred())
		Expect(fmt.Sprintf("%x", prp.ProposalHash)).To(Equal("6fa450b00ebef6c7de9f3479148f6d6ff2c645762e17fcaae989ff7b668be001"))

		ccAct := &pb.ChaincodeAction{}
		err = proto.Unmarshal(prp.Extension, ccAct)
		Expect(err).NotTo(HaveOccurred())
		Expect(ccAct.Events).To(Equal(protoutil.MarshalOrPanic(chaincodeEvent)))
		Expect(proto.Equal(ccAct.Response, &pb.Response{
			Status:  200,
			Payload: []byte("response-payload"),
		})).To(BeTrue())
		Expect(fakeSupport.GetHistoryQueryExecutorCallCount()).To(Equal(1))
		ledgerName := fakeSupport.GetHistoryQueryExecutorArgsForCall(0)
		Expect(ledgerName).To(Equal("channel-id"))
	})

	Context("when the chaincode endorsement fails", func() {
		BeforeEach(func() {
			fakeSupport.EndorseWithPluginReturns(nil, nil, fmt.Errorf("fake-endorserment-error"))
		})

		It("returns the error, but with no payload encoded", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Payload).To(BeNil())
			Expect(proposalResponse.Response).To(Equal(&pb.Response{
				Status:  500,
				Message: "endorsing with plugin failed: fake-endorserment-error",
			}))
		})
	})

	It("checks for duplicate transactions", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeSupport.GetTransactionByIDCallCount()).To(Equal(1))
		channelID, txid := fakeSupport.GetTransactionByIDArgsForCall(0)
		Expect(channelID).To(Equal("channel-id"))
		Expect(txid).To(Equal("6f142589e4ef6a1e62c9c816e2074f70baa9f7cf67c2f0c287d4ef907d6d2015"))
	})

	Context("when the txid is duplicated", func() {
		BeforeEach(func() {
			fakeSupport.GetTransactionByIDReturns(nil, nil)
		})

		It("wraps and returns an error and responds to the client", func() {
			proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
			Expect(err).To(MatchError("duplicate transaction found [6f142589e4ef6a1e62c9c816e2074f70baa9f7cf67c2f0c287d4ef907d6d2015]. Creator [0a066d73702d6964]"))
			Expect(proposalResponse).To(Equal(&pb.ProposalResponse{
				Response: &pb.Response{
					Status:  500,
					Message: "duplicate transaction found [6f142589e4ef6a1e62c9c816e2074f70baa9f7cf67c2f0c287d4ef907d6d2015]. Creator [0a066d73702d6964]",
				},
			}))
		})
	})

	It("gets a transaction simulator", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeSupport.GetTxSimulatorCallCount()).To(Equal(1))
		ledgerName, txid := fakeSupport.GetTxSimulatorArgsForCall(0)
		Expect(ledgerName).To(Equal("channel-id"))
		Expect(txid).To(Equal("6f142589e4ef6a1e62c9c816e2074f70baa9f7cf67c2f0c287d4ef907d6d2015"))
	})

	Context("when getting the tx simulator fails", func() {
		BeforeEach(func() {
			fakeSupport.GetTxSimulatorReturns(nil, fmt.Errorf("fake-simulator-error"))
		})

		It("returns a response with the error", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Payload).To(BeNil())
			Expect(proposalResponse.Response).To(Equal(&pb.Response{
				Status:  500,
				Message: "fake-simulator-error",
			}))
		})
	})

	It("gets a history query executor", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeSupport.GetHistoryQueryExecutorCallCount()).To(Equal(1))
		ledgerName := fakeSupport.GetHistoryQueryExecutorArgsForCall(0)
		Expect(ledgerName).To(Equal("channel-id"))
	})

	Context("when getting the history query executor fails", func() {
		BeforeEach(func() {
			fakeSupport.GetHistoryQueryExecutorReturns(nil, fmt.Errorf("fake-history-error"))
		})

		It("returns a response with the error", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Payload).To(BeNil())
			Expect(proposalResponse.Response).To(Equal(&pb.Response{
				Status:  500,
				Message: "fake-history-error",
			}))
		})
	})

	It("gets the channel context", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeChannelFetcher.ChannelCallCount()).To(Equal(1))
		channelID := fakeChannelFetcher.ChannelArgsForCall(0)
		Expect(channelID).To(Equal("channel-id"))
	})

	Context("when the channel context cannot be retrieved", func() {
		BeforeEach(func() {
			fakeChannelFetcher.ChannelReturns(nil)
		})

		It("returns a response with the error", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Payload).To(BeNil())
			Expect(proposalResponse.Response).To(Equal(&pb.Response{
				Status:  500,
				Message: "channel 'channel-id' not found",
			}))
		})
	})

	It("checks the submitter's identity", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeChannelMSPIdentityDeserializer.DeserializeIdentityCallCount()).To(Equal(1))
		identity := fakeChannelMSPIdentityDeserializer.DeserializeIdentityArgsForCall(0)
		Expect(identity).To(Equal(protoutil.MarshalOrPanic(&mspproto.SerializedIdentity{
			Mspid: "msp-id",
		})))

		Expect(fakeLocalMSPIdentityDeserializer.DeserializeIdentityCallCount()).To(Equal(0))
	})

	Context("when the proposal is not validly signed", func() {
		BeforeEach(func() {
			fakeChannelMSPIdentityDeserializer.DeserializeIdentityReturns(nil, fmt.Errorf("fake-deserialize-error"))
		})

		It("wraps and returns an error and responds to the client", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).To(MatchError("error validating proposal: access denied: channel [channel-id] creator org [msp-id]"))
			Expect(proposalResponse).To(Equal(&pb.ProposalResponse{
				Response: &pb.Response{
					Status:  500,
					Message: "error validating proposal: access denied: channel [channel-id] creator org [msp-id]",
				},
			}))
		})
	})

	It("checks the ACLs for the identity", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeProposalACLCheckFailed.WithCallCount()).To(Equal(0))
		Expect(fakeInitFailed.WithCallCount()).To(Equal(0))
		Expect(fakeEndorsementsFailed.WithCallCount()).To(Equal(0))
		Expect(fakeDuplicateTxsFailure.WithCallCount()).To(Equal(0))
	})

	Context("when the acl check fails", func() {
		BeforeEach(func() {
			fakeSupport.CheckACLReturns(fmt.Errorf("fake-acl-error"))
		})

		It("wraps and returns an error and responds to the client", func() {
			proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
			Expect(err).To(MatchError("fake-acl-error"))
			Expect(proposalResponse).To(Equal(&pb.ProposalResponse{
				Response: &pb.Response{
					Status:  500,
					Message: "fake-acl-error",
				},
			}))
		})

		Context("when it's for a system chaincode", func() {
			BeforeEach(func() {
				fakeSupport.IsSysCCReturns(true)
			})

			It("skips the acl check", func() {
				proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
				Expect(err).NotTo(HaveOccurred())
				Expect(proposalResponse.Response.Status).To(Equal(int32(200)))
			})
		})
	})

	It("gets the chaincode definition", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeSupport.ChaincodeEndorsementInfoCallCount()).To(Equal(1))
		channelID, chaincodeName, txSim := fakeSupport.ChaincodeEndorsementInfoArgsForCall(0)
		Expect(channelID).To(Equal("channel-id"))
		Expect(chaincodeName).To(Equal("chaincode-name"))
		Expect(txSim).To(Equal(fakeTxSimulator))
	})

	Context("when the chaincode definition is not found", func() {
		BeforeEach(func() {
			fakeSupport.ChaincodeEndorsementInfoReturns(nil, fmt.Errorf("fake-definition-error"))
		})

		It("returns an error in the response", func() {
			proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Response).To(Equal(&pb.Response{
				Status:  500,
				Message: "make sure the chaincode chaincode-name has been successfully defined on channel channel-id and try again: fake-definition-error",
			}))
		})
	})

	It("calls the chaincode", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeSupport.ExecuteCallCount()).To(Equal(1))
		txParams, chaincodeName, input := fakeSupport.ExecuteArgsForCall(0)
		Expect(txParams.ChannelID).To(Equal("channel-id"))
		Expect(txParams.SignedProp).To(Equal(signedProposal))
		Expect(txParams.TXSimulator).To(Equal(fakeTxSimulator))
		Expect(txParams.HistoryQueryExecutor).To(Equal(fakeHistoryQueryExecutor))
		Expect(chaincodeName).To(Equal("chaincode-name"))
		Expect(proto.Equal(input, &pb.ChaincodeInput{
			Args: [][]byte{[]byte("arg1"), []byte("arg2"), []byte("arg3")},
		})).To(BeTrue())
	})

	Context("when calling the chaincode returns an error", func() {
		BeforeEach(func() {
			fakeSupport.ExecuteReturns(nil, nil, fmt.Errorf("fake-chaincode-execution-error"))
		})

		It("returns a response with the error and no payload", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Payload).To(BeNil())
			Expect(proposalResponse.Response).To(Equal(&pb.Response{
				Status:  500,
				Message: "error in simulation: fake-chaincode-execution-error",
			}))
		})
	})

	It("distributes private data", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakePrivateDataDistributor.DistributePrivateDataCallCount()).To(Equal(1))
		cid, txid, privateData, blkHt := fakePrivateDataDistributor.DistributePrivateDataArgsForCall(0)
		Expect(cid).To(Equal("channel-id"))
		Expect(txid).To(Equal("6f142589e4ef6a1e62c9c816e2074f70baa9f7cf67c2f0c287d4ef907d6d2015"))
		Expect(blkHt).To(Equal(uint64(7)))

		// TODO, this deserves a better test, but there was none before and this logic,
		// really seems far too jumbled to be in the endorser package.  There are seperate
		// tests of the private data assembly functions in their test file.
		Expect(privateData).NotTo(BeNil())
		Expect(privateData.EndorsedAt).To(Equal(uint64(7)))
	})

	Context("when the private data cannot be distributed", func() {
		BeforeEach(func() {
			fakePrivateDataDistributor.DistributePrivateDataReturns(fmt.Errorf("fake-private-data-error"))
		})

		It("returns a response with the error and no payload", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Payload).To(BeNil())
			Expect(proposalResponse.Response).To(Equal(&pb.Response{
				Status:  500,
				Message: "error in simulation: fake-private-data-error",
			}))
		})
	})

	It("checks the block height", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeSupport.GetLedgerHeightCallCount()).To(Equal(1))
	})

	Context("when the block height cannot be determined", func() {
		BeforeEach(func() {
			fakeSupport.GetLedgerHeightReturns(0, fmt.Errorf("fake-block-height-error"))
		})

		It("returns a response with the error and no payload", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Payload).To(BeNil())
			Expect(proposalResponse.Response).To(Equal(&pb.Response{
				Status:  500,
				Message: "error in simulation: failed to obtain ledger height for channel 'channel-id': fake-block-height-error",
			}))
		})
	})

	It("records metrics about the proposal processing", func() {
		_, err := e.ProcessProposal(context.Background(), signedProposal)
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeProposalsReceived.AddCallCount()).To(Equal(1))
		Expect(fakeSuccessfulProposals.AddCallCount()).To(Equal(1))
		Expect(fakeProposalValidationFailed.AddCallCount()).To(Equal(0))

		Expect(fakeProposalDuration.WithCallCount()).To(Equal(1))
		Expect(fakeProposalDuration.WithArgsForCall(0)).To(Equal([]string{
			"channel", "channel-id",
			"chaincode", "chaincode-name",
			"success", "true",
		}))
	})

	Context("when the channel id is empty", func() {
		BeforeEach(func() {
			channelID = ""
		})

		It("returns a successful proposal response with no endorsement", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Endorsement).To(BeNil())
			Expect(proposalResponse.Timestamp).To(BeNil())
			Expect(proposalResponse.Version).To(Equal(int32(0)))
			Expect(proposalResponse.Payload).To(BeNil())
			Expect(proto.Equal(proposalResponse.Response, &pb.Response{
				Status:  200,
				Payload: []byte("response-payload"),
			})).To(BeTrue())
		})

		It("does not attempt to get a history query executor", func() {
			_, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeSupport.GetHistoryQueryExecutorCallCount()).To(Equal(0))
		})

		It("does not attempt to deduplicate the txid", func() {
			_, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeSupport.GetTransactionByIDCallCount()).To(Equal(0))
		})

		It("does not attempt to get a tx simulator", func() {
			_, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeSupport.GetTxSimulatorCallCount()).To(Equal(0))
		})

		It("uses the local MSP to authorize the creator", func() {
			_, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeChannelMSPIdentityDeserializer.DeserializeIdentityCallCount()).To(Equal(0))

			Expect(fakeLocalMSPIdentityDeserializer.DeserializeIdentityCallCount()).To(Equal(1))
			identity := fakeLocalMSPIdentityDeserializer.DeserializeIdentityArgsForCall(0)
			Expect(identity).To(Equal(protoutil.MarshalOrPanic(&mspproto.SerializedIdentity{
				Mspid: "msp-id",
			})))
		})

		Context("when the proposal is not validly signed", func() {
			BeforeEach(func() {
				fakeLocalMSPIdentityDeserializer.DeserializeIdentityReturns(nil, fmt.Errorf("fake-deserialize-error"))
			})

			It("wraps and returns an error and responds to the client", func() {
				proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
				Expect(err).To(MatchError("error validating proposal: access denied: channel [] creator org [msp-id]"))
				Expect(proposalResponse).To(Equal(&pb.ProposalResponse{
					Response: &pb.Response{
						Status:  500,
						Message: "error validating proposal: access denied: channel [] creator org [msp-id]",
					},
				}))
			})
		})

		It("records metrics but without a channel ID set", func() {
			_, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeProposalsReceived.AddCallCount()).To(Equal(1))
			Expect(fakeSuccessfulProposals.AddCallCount()).To(Equal(1))
			Expect(fakeProposalValidationFailed.AddCallCount()).To(Equal(0))

			Expect(fakeProposalDuration.WithCallCount()).To(Equal(1))
			Expect(fakeProposalDuration.WithArgsForCall(0)).To(Equal([]string{
				"channel", "",
				"chaincode", "chaincode-name",
				"success", "true",
			}))
			Expect(fakeProposalACLCheckFailed.WithCallCount()).To(Equal(0))
			Expect(fakeInitFailed.WithCallCount()).To(Equal(0))
			Expect(fakeEndorsementsFailed.WithCallCount()).To(Equal(0))
			Expect(fakeDuplicateTxsFailure.WithCallCount()).To(Equal(0))
		})

		Context("when the chaincode response is >= 500", func() {
			BeforeEach(func() {
				chaincodeResponse.Status = 500
			})

			It("returns the result, but with the proposal encoded, and no endorsements", func() {
				proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
				Expect(err).NotTo(HaveOccurred())
				Expect(proposalResponse.Endorsement).To(BeNil())
				Expect(proposalResponse.Timestamp).To(BeNil())
				Expect(proposalResponse.Version).To(Equal(int32(0)))
				Expect(proto.Equal(proposalResponse.Response, &pb.Response{
					Status:  500,
					Payload: []byte("response-payload"),
				})).To(BeTrue())

				// This is almost definitely a bug, but, adding a test in case someone is relying on this behavior.
				// When the response is >= 500, we return a payload, but not on success.  A payload is only meaningful
				// if it is endorsed, so it's unclear why we're returning it here.
				prp := &pb.ProposalResponsePayload{}
				err = proto.Unmarshal(proposalResponse.Payload, prp)
				Expect(err).NotTo(HaveOccurred())
				Expect(fmt.Sprintf("%x", prp.ProposalHash)).To(Equal("f2c27f04f897dc28fd1b2983e7b22ebc8fbbb3d0617c140d913b33e463886788"))

				ccAct := &pb.ChaincodeAction{}
				err = proto.Unmarshal(prp.Extension, ccAct)
				Expect(err).NotTo(HaveOccurred())
				Expect(proto.Equal(ccAct.Response, &pb.Response{
					Status:  500,
					Payload: []byte("response-payload"),
				})).To(BeTrue())

				// This is an especially weird bit of the behavior, the chaincode event is nil-ed before creating
				// the proposal response. (That probably shouldn't be created)
				Expect(ccAct.Events).To(BeNil())
			})
		})

		Context("when the 200 < chaincode response < 500", func() {
			BeforeEach(func() {
				chaincodeResponse.Status = 499
			})

			It("returns the result, but with the proposal encoded, and no endorsements", func() {
				proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
				Expect(err).NotTo(HaveOccurred())
				Expect(proposalResponse.Endorsement).To(BeNil())
				Expect(proposalResponse.Timestamp).To(BeNil())
				Expect(proposalResponse.Version).To(Equal(int32(0)))
				Expect(proto.Equal(proposalResponse.Response, &pb.Response{
					Status:  499,
					Payload: []byte("response-payload"),
				})).To(BeTrue())
				Expect(proposalResponse.Payload).To(BeNil())
			})
		})
	})

	Context("when the proposal is malformed", func() {
		JustBeforeEach(func() {
			signedProposal = &pb.SignedProposal{
				ProposalBytes: []byte("garbage"),
			}
		})

		It("wraps and returns an error and responds to the client", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).To(MatchError("error unmarshaling Proposal: proto: can't skip unknown wire type 7"))
			Expect(proposalResponse).To(Equal(&pb.ProposalResponse{
				Response: &pb.Response{
					Status:  500,
					Message: "error unmarshaling Proposal: proto: can't skip unknown wire type 7",
				},
			}))
		})
	})

	Context("when the chaincode response is >= 500", func() {
		BeforeEach(func() {
			chaincodeResponse.Status = 500
		})

		It("returns the result, but with the proposal encoded, and no endorsements", func() {
			proposalResponse, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Endorsement).To(BeNil())
			Expect(proposalResponse.Timestamp).To(BeNil())
			Expect(proposalResponse.Version).To(Equal(int32(0)))
			Expect(proto.Equal(proposalResponse.Response, &pb.Response{
				Status:  500,
				Payload: []byte("response-payload"),
			})).To(BeTrue())
			Expect(proposalResponse.Payload).NotTo(BeNil())
		})
	})

	Context("when the chaincode name is qscc", func() {
		BeforeEach(func() {
			chaincodeName = "qscc"
		})

		It("skips fetching the tx simulator and history query exucutor", func() {
			_, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeSupport.GetTxSimulatorCallCount()).To(Equal(0))
			Expect(fakeSupport.GetHistoryQueryExecutorCallCount()).To(Equal(0))
		})
	})

	Context("when the chaincode name is cscc", func() {
		BeforeEach(func() {
			chaincodeName = "cscc"
		})

		It("skips fetching the tx simulator and history query exucutor", func() {
			_, err := e.ProcessProposal(context.Background(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeSupport.GetTxSimulatorCallCount()).To(Equal(0))
			Expect(fakeSupport.GetHistoryQueryExecutorCallCount()).To(Equal(0))
		})
	})

	Context("when the chaincode response is >= 400 but < 500", func() {
		BeforeEach(func() {
			chaincodeResponse.Status = 400
		})

		It("returns the response with no payload", func() {
			proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proposalResponse.Payload).To(BeNil())
			Expect(proto.Equal(proposalResponse.Response, &pb.Response{
				Status:  400,
				Payload: []byte("response-payload"),
			})).To(BeTrue())
		})
	})

	Context("when we're in the degenerate legacy lifecycle case", func() {
		BeforeEach(func() {
			chaincodeName = "lscc"
			chaincodeInput.Args = [][]byte{
				[]byte("deploy"),
				nil,
				protoutil.MarshalOrPanic(&pb.ChaincodeDeploymentSpec{
					ChaincodeSpec: &pb.ChaincodeSpec{
						ChaincodeId: &pb.ChaincodeID{
							Name:    "deploy-name",
							Version: "deploy-version",
						},
						Input: &pb.ChaincodeInput{
							Args: [][]byte{[]byte("target-arg")},
						},
					},
				}),
			}

			fakeTxSimulator.GetTxSimulationResultsReturns(
				&ledger.TxSimulationResults{
					PubSimulationResults: &rwset.TxReadWriteSet{},
					// We don't return private data in this case because lscc forbids it
				},
				nil,
			)
		})

		It("triggers the legacy init, and returns the response from lscc", func() {
			proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
			Expect(err).NotTo(HaveOccurred())
			Expect(proto.Equal(proposalResponse.Response, &pb.Response{
				Status:  200,
				Payload: []byte("response-payload"),
			})).To(BeTrue())

			Expect(fakeSupport.ExecuteLegacyInitCallCount()).To(Equal(1))
			_, name, version, input := fakeSupport.ExecuteLegacyInitArgsForCall(0)
			Expect(name).To(Equal("deploy-name"))
			Expect(version).To(Equal("deploy-version"))
			Expect(input.Args).To(Equal([][]byte{[]byte("target-arg")}))
		})

		Context("when the chaincode spec contains a code package", func() {
			BeforeEach(func() {
				chaincodeInput.Args = [][]byte{
					[]byte("deploy"),
					nil,
					protoutil.MarshalOrPanic(&pb.ChaincodeDeploymentSpec{
						ChaincodeSpec: &pb.ChaincodeSpec{
							ChaincodeId: &pb.ChaincodeID{
								Name:    "deploy-name",
								Version: "deploy-version",
							},
							Input: &pb.ChaincodeInput{
								Args: [][]byte{[]byte("target-arg")},
							},
						},
						CodePackage: []byte("some-code"),
					}),
				}
			})

			It("returns an error to the client", func() {
				proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
				Expect(err).NotTo(HaveOccurred())
				Expect(proposalResponse.Response).To(Equal(&pb.Response{
					Status:  500,
					Message: "error in simulation: lscc upgrade/deploy should not include a code packages",
				}))
			})
		})

		Context("when the simulation uses private data", func() {
			BeforeEach(func() {
				fakeTxSimulator.GetTxSimulationResultsReturns(
					&ledger.TxSimulationResults{
						PubSimulationResults: &rwset.TxReadWriteSet{},
						PvtSimulationResults: &rwset.TxPvtReadWriteSet{},
					},
					nil,
				)
			})

			It("returns an error to the client", func() {
				proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
				Expect(err).NotTo(HaveOccurred())
				Expect(proposalResponse.Response).To(Equal(&pb.Response{
					Status:  500,
					Message: "error in simulation: Private data is forbidden to be used in instantiate",
				}))
			})
		})

		Context("when the init fails", func() {
			BeforeEach(func() {
				fakeSupport.ExecuteLegacyInitReturns(nil, nil, fmt.Errorf("fake-legacy-init-error"))
			})

			It("returns an error and increments the metric", func() {
				proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
				Expect(err).NotTo(HaveOccurred())
				Expect(proposalResponse.Response).To(Equal(&pb.Response{
					Status:  500,
					Message: "error in simulation: fake-legacy-init-error",
				}))

				Expect(fakeInitFailed.WithCallCount()).To(Equal(1))
				Expect(fakeInitFailed.WithArgsForCall(0)).To(Equal([]string{
					"channel", "channel-id",
					"chaincode", "deploy-name",
				}))
			})
		})

		Context("when the deploying chaincode is the name of a builtin system chaincode", func() {
			BeforeEach(func() {
				fakeSupport.IsSysCCStub = func(name string) bool {
					return name == "deploy-name"
				}
			})

			It("triggers the legacy init, and returns the response from lscc", func() {
				proposalResponse, err := e.ProcessProposal(context.TODO(), signedProposal)
				Expect(err).NotTo(HaveOccurred())
				Expect(proposalResponse.Response).To(Equal(&pb.Response{
					Status:  500,
					Message: "error in simulation: attempting to deploy a system chaincode deploy-name/channel-id",
				}))
			})
		})
	})
})
=======
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)

	tempDir := newTempDir()
	viper.Set("peer.fileSystemPath", filepath.Join(tempDir, "hyperledger", "production"))

	peerAddress, err := peer.GetLocalAddress()
	if err != nil {
		return nil, fmt.Errorf("Error obtaining peer address: %s", err)
	}
	lis, err := net.Listen("tcp", peerAddress)
	if err != nil {
		return nil, fmt.Errorf("Error starting peer listener %s", err)
	}

	//initialize ledger
	peer.MockInitialize()

	mspGetter := func(cid string) []string {
		return []string{"DEFAULT"}
	}

	peer.MockSetMSPIDGetter(mspGetter)

	getPeerEndpoint := func() (*pb.PeerEndpoint, error) {
		return &pb.PeerEndpoint{Id: &pb.PeerID{Name: "testpeer"}, Address: peerAddress}, nil
	}

	ccStartupTimeout := time.Duration(30000) * time.Millisecond
	pb.RegisterChaincodeSupportServer(grpcServer, chaincode.NewChaincodeSupport(getPeerEndpoint, false, ccStartupTimeout))

	syscc.RegisterSysCCs()

	if err = peer.MockCreateChain(chainID); err != nil {
		closeListenerAndSleep(lis)
		return nil, err
	}

	syscc.DeploySysCCs(chainID)

	go grpcServer.Serve(lis)

	return &testEnvironment{tempDir: tempDir, listener: lis}, nil
}

func finitPeer(tev *testEnvironment) {
	closeListenerAndSleep(tev.listener)
	os.RemoveAll(tev.tempDir)
}

func closeListenerAndSleep(l net.Listener) {
	if l != nil {
		l.Close()
		time.Sleep(2 * time.Second)
	}
}

// getInvokeProposal gets the proposal for the chaincode invocation
// Currently supported only for Invokes
// It returns the proposal and the transaction id associated to the proposal
func getInvokeProposal(cis *pb.ChaincodeInvocationSpec, chainID string, creator []byte) (*pb.Proposal, string, error) {
	return pbutils.CreateChaincodeProposal(common.HeaderType_ENDORSER_TRANSACTION, chainID, cis, creator)
}

// getInvokeProposalOverride allows to get a proposal for the chaincode invocation
// overriding transaction id and nonce which are by default auto-generated.
// It returns the proposal and the transaction id associated to the proposal
func getInvokeProposalOverride(txid string, cis *pb.ChaincodeInvocationSpec, chainID string, nonce, creator []byte) (*pb.Proposal, string, error) {
	return pbutils.CreateChaincodeProposalWithTxIDNonceAndTransient(txid, common.HeaderType_ENDORSER_TRANSACTION, chainID, cis, nonce, creator, nil)
}

func getDeployProposal(cds *pb.ChaincodeDeploymentSpec, chainID string, creator []byte) (*pb.Proposal, error) {
	return getDeployOrUpgradeProposal(cds, chainID, creator, false)
}

func getUpgradeProposal(cds *pb.ChaincodeDeploymentSpec, chainID string, creator []byte) (*pb.Proposal, error) {
	return getDeployOrUpgradeProposal(cds, chainID, creator, true)
}

//getDeployOrUpgradeProposal gets the proposal for the chaincode deploy or upgrade
//the payload is a ChaincodeDeploymentSpec
func getDeployOrUpgradeProposal(cds *pb.ChaincodeDeploymentSpec, chainID string, creator []byte, upgrade bool) (*pb.Proposal, error) {
	//we need to save off the chaincode as we have to instantiate with nil CodePackage
	var err error
	if err = ccprovider.PutChaincodeIntoFS(cds); err != nil {
		return nil, err
	}

	cds.CodePackage = nil

	b, err := proto.Marshal(cds)
	if err != nil {
		return nil, err
	}

	var propType string
	if upgrade {
		propType = "upgrade"
	} else {
		propType = "deploy"
	}
	sccver := util.GetSysCCVersion()
	//wrap the deployment in an invocation spec to lscc...
	lsccSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_GOLANG, ChaincodeId: &pb.ChaincodeID{Name: "lscc", Version: sccver}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte(propType), []byte(chainID), b}}}}

	//...and get the proposal for it
	var prop *pb.Proposal
	if prop, _, err = getInvokeProposal(lsccSpec, chainID, creator); err != nil {
		return nil, err
	}

	return prop, nil
}

func getSignedProposal(prop *pb.Proposal, signer msp.SigningIdentity) (*pb.SignedProposal, error) {
	propBytes, err := pbutils.GetBytesProposal(prop)
	if err != nil {
		return nil, err
	}

	signature, err := signer.Sign(propBytes)
	if err != nil {
		return nil, err
	}

	return &pb.SignedProposal{ProposalBytes: propBytes, Signature: signature}, nil
}

func getDeploymentSpec(context context.Context, spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	codePackageBytes, err := container.GetChaincodePackageBytes(spec)
	if err != nil {
		return nil, err
	}
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

func deploy(endorserServer pb.EndorserServer, chainID string, spec *pb.ChaincodeSpec, f func(*pb.ChaincodeDeploymentSpec)) (*pb.ProposalResponse, *pb.Proposal, error) {
	return deployOrUpgrade(endorserServer, chainID, spec, f, false)
}

func upgrade(endorserServer pb.EndorserServer, chainID string, spec *pb.ChaincodeSpec, f func(*pb.ChaincodeDeploymentSpec)) (*pb.ProposalResponse, *pb.Proposal, error) {
	return deployOrUpgrade(endorserServer, chainID, spec, f, true)
}

func deployOrUpgrade(endorserServer pb.EndorserServer, chainID string, spec *pb.ChaincodeSpec, f func(*pb.ChaincodeDeploymentSpec), upgrade bool) (*pb.ProposalResponse, *pb.Proposal, error) {
	var err error
	var depSpec *pb.ChaincodeDeploymentSpec

	ctxt := context.Background()
	depSpec, err = getDeploymentSpec(ctxt, spec)
	if err != nil {
		return nil, nil, err
	}

	if f != nil {
		f(depSpec)
	}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, nil, err
	}

	var prop *pb.Proposal
	if upgrade {
		prop, err = getUpgradeProposal(depSpec, chainID, creator)
	} else {
		prop, err = getDeployProposal(depSpec, chainID, creator)
	}
	if err != nil {
		return nil, nil, err
	}

	var signedProp *pb.SignedProposal
	signedProp, err = getSignedProposal(prop, signer)
	if err != nil {
		return nil, nil, err
	}

	var resp *pb.ProposalResponse
	resp, err = endorserServer.ProcessProposal(context.Background(), signedProp)

	return resp, prop, err
}

func invoke(chainID string, spec *pb.ChaincodeSpec) (*pb.Proposal, *pb.ProposalResponse, string, []byte, error) {
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, nil, "", nil, err
	}

	var prop *pb.Proposal
	prop, txID, err := getInvokeProposal(invocation, chainID, creator)
	if err != nil {
		return nil, nil, "", nil, fmt.Errorf("Error creating proposal  %s: %s\n", spec.ChaincodeId, err)
	}

	nonce, err := pbutils.GetNonce(prop)
	if err != nil {
		return nil, nil, "", nil, fmt.Errorf("Failed getting nonce  %s: %s\n", spec.ChaincodeId, err)
	}

	var signedProp *pb.SignedProposal
	signedProp, err = getSignedProposal(prop, signer)
	if err != nil {
		return nil, nil, "", nil, err
	}

	resp, err := endorserServer.ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return nil, nil, "", nil, err
	}

	return prop, resp, txID, nonce, err
}

func invokeWithOverride(txid string, chainID string, spec *pb.ChaincodeSpec, nonce []byte) (*pb.ProposalResponse, error) {
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, err
	}

	var prop *pb.Proposal
	prop, _, err = getInvokeProposalOverride(txid, invocation, chainID, nonce, creator)
	if err != nil {
		return nil, fmt.Errorf("Error creating proposal with override  %s %s: %s\n", txid, spec.ChaincodeId, err)
	}

	var signedProp *pb.SignedProposal
	signedProp, err = getSignedProposal(prop, signer)
	if err != nil {
		return nil, err
	}

	resp, err := endorserServer.ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return nil, fmt.Errorf("Error endorsing %s %s: %s\n", txid, spec.ChaincodeId, err)
	}

	return resp, err
}

func deleteChaincodeOnDisk(chaincodeID string) {
	os.RemoveAll(filepath.Join(config.GetPath("peer.fileSystemPath"), "chaincodes", chaincodeID))
}

//begin tests. Note that we rely upon the system chaincode and peer to be created
//once and be used for all the tests. In order to avoid dependencies / collisions
//due to deployed chaincodes, trying to use different chaincodes for different
//tests

//TestDeploy deploy chaincode example01
func TestDeploy(t *testing.T) {
	chainID := util.GetTestChainID()
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "ex01", Path: "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example01", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}}}
	defer deleteChaincodeOnDisk("ex01.0")

	cccid := ccprovider.NewCCContext(chainID, "ex01", "0", "", false, nil, nil)

	_, _, err := deploy(endorserServer, chainID, spec, nil)
	if err != nil {
		t.Fail()
		t.Logf("Deploy-error in deploy %s", err)
		chaincode.GetChain().Stop(context.Background(), cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}
	chaincode.GetChain().Stop(context.Background(), cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
}

//REMOVE WHEN JAVA CC IS ENABLED
func TestJavaDeploy(t *testing.T) {
	chainID := util.GetTestChainID()
	//pretend this is a java CC (type 4)
	spec := &pb.ChaincodeSpec{Type: 4, ChaincodeId: &pb.ChaincodeID{Name: "javacc", Path: "../../examples/chaincode/java/chaincode_example02", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}}}
	defer deleteChaincodeOnDisk("javacc.0")

	cccid := ccprovider.NewCCContext(chainID, "javacc", "0", "", false, nil, nil)

	_, _, err := deploy(endorserServer, chainID, spec, nil)
	if err == nil {
		t.Fail()
		t.Logf("expected java CC deploy to fail")
		chaincode.GetChain().Stop(context.Background(), cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}
	chaincode.GetChain().Stop(context.Background(), cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
}

func TestJavaCheckWithDifferentPackageTypes(t *testing.T) {
	//try SignedChaincodeDeploymentSpec with go chaincode (type 1)
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "gocc", Path: "path/to/cc", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("someargs")}}}
	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: []byte("some code")}
	env := &common.Envelope{Payload: pbutils.MarshalOrPanic(&common.Payload{Data: pbutils.MarshalOrPanic(&pb.SignedChaincodeDeploymentSpec{ChaincodeDeploymentSpec: pbutils.MarshalOrPanic(cds)})})}
	//wrap the package in an invocation spec to lscc...
	b := pbutils.MarshalOrPanic(env)

	lsccCID := &pb.ChaincodeID{Name: "lscc", Version: util.GetSysCCVersion()}
	lsccSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_GOLANG, ChaincodeId: lsccCID, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("install"), b}}}}

	e := &Endorser{}
	err := e.disableJavaCCInst(lsccCID, lsccSpec)
	assert.Nil(t, err)

	//now try plain ChaincodeDeploymentSpec...should succeed (go chaincode)
	b = pbutils.MarshalOrPanic(cds)

	lsccSpec = &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_GOLANG, ChaincodeId: lsccCID, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("install"), b}}}}
	err = e.disableJavaCCInst(lsccCID, lsccSpec)
	assert.Nil(t, err)
}

//TestRedeploy - deploy two times, second time should fail but example02 should remain deployed
func TestRedeploy(t *testing.T) {
	chainID := util.GetTestChainID()

	//invalid arguments
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "ex02", Path: "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}}}

	defer deleteChaincodeOnDisk("ex02.0")

	cccid := ccprovider.NewCCContext(chainID, "ex02", "0", "", false, nil, nil)

	_, _, err := deploy(endorserServer, chainID, spec, nil)
	if err != nil {
		t.Fail()
		t.Logf("error in endorserServer.ProcessProposal %s", err)
		chaincode.GetChain().Stop(context.Background(), cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	deleteChaincodeOnDisk("ex02.0")

	//second time should not fail as we are just simulating
	_, _, err = deploy(endorserServer, chainID, spec, nil)
	if err != nil {
		t.Fail()
		t.Logf("error in endorserServer.ProcessProposal %s", err)
		chaincode.GetChain().Stop(context.Background(), cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}
	chaincode.GetChain().Stop(context.Background(), cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
}

// TestDeployAndInvoke deploys and invokes chaincode_example01
func TestDeployAndInvoke(t *testing.T) {
	chainID := util.GetTestChainID()
	var ctxt = context.Background()

	url := "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example01"
	chaincodeID := &pb.ChaincodeID{Path: url, Name: "ex01", Version: "0"}

	defer deleteChaincodeOnDisk("ex01.0")

	args := []string{"10"}

	f := "init"
	argsDeploy := util.ToChaincodeArgs(f, "a", "100", "b", "200")
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: argsDeploy}}

	cccid := ccprovider.NewCCContext(chainID, "ex01", "0", "", false, nil, nil)

	resp, prop, err := deploy(endorserServer, chainID, spec, nil)
	chaincodeID1 := spec.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error deploying <%s>: %s", chaincodeID1, err)
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}
	var nextBlockNumber uint64 = 1 // first block needs to be block number = 1. Genesis block is block 0
	err = endorserServer.(*Endorser).commitTxSimulation(prop, chainID, signer, resp, nextBlockNumber)
	if err != nil {
		t.Fail()
		t.Logf("Error committing deploy <%s>: %s", chaincodeID1, err)
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	f = "invoke"
	invokeArgs := append([]string{f}, args...)
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: util.ToChaincodeArgs(invokeArgs...)}}
	prop, resp, txid, nonce, err := invoke(chainID, spec)
	if err != nil {
		t.Fail()
		t.Logf("Error invoking transaction: %s", err)
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}
	// Commit invoke
	nextBlockNumber++
	err = endorserServer.(*Endorser).commitTxSimulation(prop, chainID, signer, resp, nextBlockNumber)
	if err != nil {
		t.Fail()
		t.Logf("Error committing first invoke <%s>: %s", chaincodeID1, err)
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	// Now test for an invalid TxID
	f = "invoke"
	invokeArgs = append([]string{f}, args...)
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: util.ToChaincodeArgs(invokeArgs...)}}
	_, err = invokeWithOverride("invalid_tx_id", chainID, spec, nonce)
	if err == nil {
		t.Fail()
		t.Log("Replay attack protection faild. Transaction with invalid txid passed")
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	// Now test for duplicated TxID
	f = "invoke"
	invokeArgs = append([]string{f}, args...)
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: util.ToChaincodeArgs(invokeArgs...)}}
	_, err = invokeWithOverride(txid, chainID, spec, nonce)
	if err == nil {
		t.Fail()
		t.Log("Replay attack protection faild. Transaction with duplicaged txid passed")
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	// Test chaincode endorsement failure when invalid function name supplied
	f = "invokeinvalid"
	invokeArgs = append([]string{f}, args...)
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: util.ToChaincodeArgs(invokeArgs...)}}
	prop, resp, txid, nonce, err = invoke(chainID, spec)
	if err == nil {
		t.Fail()
		t.Logf("expecting fabric to report error from chaincode failure")
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	} else if _, ok := err.(*chaincodeError); !ok {
		t.Fail()
		t.Logf("expecting chaincode error but found %v", err)
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	if resp != nil {
		assert.Equal(t, int32(500), resp.Response.Status, "Unexpected response status")
	}

	t.Logf("Invoke test passed")

	chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
}

// TestUpgradeAndInvoke deploys chaincode_example01, upgrade it with chaincode_example02, then invoke it
func TestDeployAndUpgrade(t *testing.T) {
	chainID := util.GetTestChainID()
	var ctxt = context.Background()

	url1 := "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example01"
	url2 := "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"
	chaincodeID1 := &pb.ChaincodeID{Path: url1, Name: "upgradeex01", Version: "0"}
	chaincodeID2 := &pb.ChaincodeID{Path: url2, Name: "upgradeex01", Version: "1"}

	defer deleteChaincodeOnDisk("upgradeex01.0")
	defer deleteChaincodeOnDisk("upgradeex01.1")

	f := "init"
	argsDeploy := util.ToChaincodeArgs(f, "a", "100", "b", "200")
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID1, Input: &pb.ChaincodeInput{Args: argsDeploy}}

	cccid1 := ccprovider.NewCCContext(chainID, "upgradeex01", "0", "", false, nil, nil)
	cccid2 := ccprovider.NewCCContext(chainID, "upgradeex01", "1", "", false, nil, nil)

	resp, prop, err := deploy(endorserServer, chainID, spec, nil)

	chaincodeName := spec.ChaincodeId.Name
	if err != nil {
		t.Fail()
		chaincode.GetChain().Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID1}})
		chaincode.GetChain().Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID2}})
		t.Logf("Error deploying <%s>: %s", chaincodeName, err)
		return
	}

	var nextBlockNumber uint64 = 3 // something above created block 0
	err = endorserServer.(*Endorser).commitTxSimulation(prop, chainID, signer, resp, nextBlockNumber)
	if err != nil {
		t.Fail()
		chaincode.GetChain().Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID1}})
		chaincode.GetChain().Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID2}})
		t.Logf("Error committing <%s>: %s", chaincodeName, err)
		return
	}

	argsUpgrade := util.ToChaincodeArgs(f, "a", "150", "b", "300")
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID2, Input: &pb.ChaincodeInput{Args: argsUpgrade}}
	_, _, err = upgrade(endorserServer, chainID, spec, nil)
	if err != nil {
		t.Fail()
		chaincode.GetChain().Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID1}})
		chaincode.GetChain().Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID2}})
		t.Logf("Error upgrading <%s>: %s", chaincodeName, err)
		return
	}

	t.Logf("Upgrade test passed")

	chaincode.GetChain().Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID1}})
	chaincode.GetChain().Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID2}})
}

// TestWritersACLFail deploys a chaincode and then tries to invoke it;
// however we inject a special policy for writers to simulate
// the scenario in which the creator of this proposal is not among
// the writers for the chain
func TestWritersACLFail(t *testing.T) {
	//skip pending FAB-2457 fix
	t.Skip()
	chainID := util.GetTestChainID()
	var ctxt = context.Background()

	url := "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example01"
	chaincodeID := &pb.ChaincodeID{Path: url, Name: "ex01-fail", Version: "0"}

	defer deleteChaincodeOnDisk("ex01-fail.0")

	args := []string{"10"}

	f := "init"
	argsDeploy := util.ToChaincodeArgs(f, "a", "100", "b", "200")
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: argsDeploy}}

	cccid := ccprovider.NewCCContext(chainID, "ex01-fail", "0", "", false, nil, nil)

	resp, prop, err := deploy(endorserServer, chainID, spec, nil)
	chaincodeID1 := spec.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error deploying <%s>: %s", chaincodeID1, err)
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}
	var nextBlockNumber uint64 = 3 // The tests that ran before this test created blocks 0-2
	err = endorserServer.(*Endorser).commitTxSimulation(prop, chainID, signer, resp, nextBlockNumber)
	if err != nil {
		t.Fail()
		t.Logf("Error committing deploy <%s>: %s", chaincodeID1, err)
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	// here we inject a reject policy for writers
	// to simulate the scenario in which the invoker
	// is not authorized to issue this proposal
	rejectpolicy := &mockpolicies.Policy{
		Err: errors.New("The creator of this proposal does not fulfil the writers policy of this chain"),
	}
	pm := peer.GetPolicyManager(chainID)
	pm.(*mockpolicies.Manager).PolicyMap = map[string]policies.Policy{policies.ChannelApplicationWriters: rejectpolicy}

	f = "invoke"
	invokeArgs := append([]string{f}, args...)
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: util.ToChaincodeArgs(invokeArgs...)}}
	prop, resp, _, _, err = invoke(chainID, spec)
	if err == nil {
		t.Fail()
		t.Logf("Invocation should have failed")
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	t.Logf("TestWritersACLFail passed")

	chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
}

func TestHeaderExtensionNoChaincodeID(t *testing.T) {
	creator, _ := signer.Serialize()
	nonce := []byte{1, 2, 3}
	digest, err := factory.GetDefault().Hash(append(nonce, creator...), &bccsp.SHA256Opts{})
	txID := hex.EncodeToString(digest)
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: nil, Input: &pb.ChaincodeInput{Args: util.ToChaincodeArgs()}}
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}
	prop, _, _ := pbutils.CreateChaincodeProposalWithTxIDNonceAndTransient(txID, common.HeaderType_ENDORSER_TRANSACTION, util.GetTestChainID(), invocation, []byte{1, 2, 3}, creator, nil)
	signedProp, _ := getSignedProposal(prop, signer)
	_, err = endorserServer.ProcessProposal(context.Background(), signedProp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ChaincodeHeaderExtension.ChaincodeId is nil")
}

// TestAdminACLFail deploys tried to deploy a chaincode;
// however we inject a special policy for admins to simulate
// the scenario in which the creator of this proposal is not among
// the admins for the chain
func TestAdminACLFail(t *testing.T) {
	//skip pending FAB-2457 fix
	t.Skip()
	chainID := util.GetTestChainID()

	// here we inject a reject policy for admins
	// to simulate the scenario in which the invoker
	// is not authorized to issue this proposal
	rejectpolicy := &mockpolicies.Policy{
		Err: errors.New("The creator of this proposal does not fulfil the writers policy of this chain"),
	}
	pm := peer.GetPolicyManager(chainID)
	pm.(*mockpolicies.Manager).PolicyMap = map[string]policies.Policy{policies.ChannelApplicationAdmins: rejectpolicy}

	var ctxt = context.Background()

	url := "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example01"
	chaincodeID := &pb.ChaincodeID{Path: url, Name: "ex01-fail1", Version: "0"}

	defer deleteChaincodeOnDisk("ex01-fail1.0")

	f := "init"
	argsDeploy := util.ToChaincodeArgs(f, "a", "100", "b", "200")
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: argsDeploy}}

	cccid := ccprovider.NewCCContext(chainID, "ex01-fail1", "0", "", false, nil, nil)

	_, _, err := deploy(endorserServer, chainID, spec, nil)
	if err == nil {
		t.Fail()
		t.Logf("Deploying chaincode should have failed!")
		chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	t.Logf("TestATestAdminACLFailCLFail passed")

	chaincode.GetChain().Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
}

// TestInvokeSccFail makes sure that invoking a system chaincode fails
func TestInvokeSccFail(t *testing.T) {
	chainID := util.GetTestChainID()

	chaincodeID := &pb.ChaincodeID{Name: "escc"}
	args := util.ToChaincodeArgs("someFunc", "someArg")
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, _, _, err := invoke(chainID, spec)
	if err == nil {
		t.Logf("Invoking escc should have failed!")
		t.Fail()
		return
	}
}

func newTempDir() string {
	tempDir, err := ioutil.TempDir("", "fabric-")
	if err != nil {
		panic(err)
	}
	return tempDir
}

func TestMain(m *testing.M) {
	setupTestConfig()

	chainID := util.GetTestChainID()
	tev, err := initPeer(chainID)
	if err != nil {
		os.Exit(-1)
		fmt.Printf("Could not initialize tests")
		finitPeer(tev)
		return
	}

	endorserServer = NewEndorserServer()

	// setup the MSP manager so that we can sign/verify
	err = msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		fmt.Printf("Could not initialize msp/signer, err %s", err)
		finitPeer(tev)
		os.Exit(-1)
		return
	}
	signer, err = mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		fmt.Printf("Could not initialize msp/signer")
		finitPeer(tev)
		os.Exit(-1)
		return
	}

	retVal := m.Run()

	finitPeer(tev)

	os.Exit(retVal)
}

func setupTestConfig() {
	flag.Parse()

	// Now set the configuration file
	viper.SetEnvPrefix("CORE")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigName("endorser_test") // name of config file (without extension)
	viper.AddConfigPath("./")            // path to look for the config file in
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	testutil.SetupTestLogging()

	// Set the number of maxprocs
	runtime.GOMAXPROCS(viper.GetInt("peer.gomaxprocs"))

	// Init the BCCSP
	var bccspConfig *factory.FactoryOpts
	err = viper.UnmarshalKey("peer.BCCSP", &bccspConfig)
	if err != nil {
		bccspConfig = nil
	}

	msp.SetupBCCSPKeystoreConfig(bccspConfig, viper.GetString("peer.mspConfigPath")+"/keystore")

	err = factory.InitFactories(bccspConfig)
	if err != nil {
		panic(fmt.Errorf("Could not initialize BCCSP Factories [%s]", err))
	}
}
>>>>>>> release-1.0
