/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
<<<<<<< HEAD
	"bytes"
=======
	"fmt"
>>>>>>> release-1.0
	"sync"
	"testing"
	"time"

<<<<<<< HEAD
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/comm"
	"github.com/hyperledger/fabric/core/deliverservice"
	"github.com/hyperledger/fabric/gossip/api"
	"github.com/hyperledger/fabric/gossip/election"
	"github.com/hyperledger/fabric/gossip/util"
	"github.com/hyperledger/fabric/internal/pkg/peer/blocksprovider"
	"github.com/hyperledger/fabric/internal/pkg/peer/orderers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type embeddingDeliveryService struct {
	startOnce sync.Once
	stopOnce  sync.Once
	deliverservice.DeliverService
=======
	"github.com/hyperledger/fabric/core/deliverservice"
	"github.com/hyperledger/fabric/core/deliverservice/blocksprovider"
	"github.com/hyperledger/fabric/gossip/api"
	"github.com/hyperledger/fabric/gossip/election"
	"github.com/hyperledger/fabric/gossip/identity"
	"github.com/hyperledger/fabric/gossip/state"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type embeddingDeliveryService struct {
	deliverclient.DeliverService
>>>>>>> release-1.0
	startSignal sync.WaitGroup
	stopSignal  sync.WaitGroup
}

<<<<<<< HEAD
func newEmbeddingDeliveryService(ds deliverservice.DeliverService) *embeddingDeliveryService {
=======
func newEmbeddingDeliveryService(ds deliverclient.DeliverService) *embeddingDeliveryService {
>>>>>>> release-1.0
	eds := &embeddingDeliveryService{
		DeliverService: ds,
	}
	eds.startSignal.Add(1)
	eds.stopSignal.Add(1)
	return eds
}

func (eds *embeddingDeliveryService) waitForDeliveryServiceActivation() {
	eds.startSignal.Wait()
}

func (eds *embeddingDeliveryService) waitForDeliveryServiceTermination() {
	eds.stopSignal.Wait()
}

func (eds *embeddingDeliveryService) StartDeliverForChannel(chainID string, ledgerInfo blocksprovider.LedgerInfo, finalizer func()) error {
<<<<<<< HEAD
	eds.startOnce.Do(func() {
		eds.startSignal.Done()
	})
=======
	eds.startSignal.Done()
>>>>>>> release-1.0
	return eds.DeliverService.StartDeliverForChannel(chainID, ledgerInfo, finalizer)
}

func (eds *embeddingDeliveryService) StopDeliverForChannel(chainID string) error {
<<<<<<< HEAD
	eds.stopOnce.Do(func() {
		eds.stopSignal.Done()
	})
=======
	eds.stopSignal.Done()
>>>>>>> release-1.0
	return eds.DeliverService.StopDeliverForChannel(chainID)
}

func (eds *embeddingDeliveryService) Stop() {
	eds.DeliverService.Stop()
}

type embeddingDeliveryServiceFactory struct {
	DeliveryServiceFactory
}

<<<<<<< HEAD
func (edsf *embeddingDeliveryServiceFactory) Service(g GossipServiceAdapter, endpoints *orderers.ConnectionSource, mcs api.MessageCryptoService) deliverservice.DeliverService {
	ds := edsf.DeliveryServiceFactory.Service(g, endpoints, mcs)
	return newEmbeddingDeliveryService(ds)
=======
func (edsf *embeddingDeliveryServiceFactory) Service(g GossipService, endpoints []string, mcs api.MessageCryptoService) (deliverclient.DeliverService, error) {
	ds, _ := edsf.DeliveryServiceFactory.Service(g, endpoints, mcs)
	return newEmbeddingDeliveryService(ds), nil
>>>>>>> release-1.0
}

func TestLeaderYield(t *testing.T) {
	// Scenario: Spawn 2 peers and wait for the first one to be the leader
	// There isn't any orderer present so the leader peer won't be able to
	// connect to the orderer, and should relinquish its leadership after a while.
	// Make sure the other peer declares itself as the leader soon after.
<<<<<<< HEAD
	takeOverMaxTimeout := time.Minute
	// It's enough to make single re-try
	// There is no ordering service available anyway, hence connection timeout
	// could be shorter
	serviceConfig := &ServiceConfig{
		UseLeaderElection:          true,
		OrgLeader:                  false,
		ElectionStartupGracePeriod: election.DefStartupGracePeriod,
		// Since we ensuring gossip has stable membership, there is no need for
		// leader election to wait for stabilization
		ElectionMembershipSampleInterval: time.Millisecond * 100,
		ElectionLeaderAliveThreshold:     time.Second * 5,
		// Test case has only two instance + making assertions only after membership view
		// is stable, hence election duration could be shorter
		ElectionLeaderElectionDuration: time.Millisecond * 500,
	}
	n := 2
	gossips := startPeers(t, serviceConfig, n, 0, 1)
=======
	deliverclient.SetReconnectTotalTimeThreshold(time.Second * 5)
	viper.Set("peer.gossip.useLeaderElection", true)
	viper.Set("peer.gossip.orgLeader", false)
	n := 2
	portPrefix := 30000
	gossips := startPeers(t, n, portPrefix)
>>>>>>> release-1.0
	defer stopPeers(gossips)
	channelName := "channelA"
	peerIndexes := []int{0, 1}
	// Add peers to the channel
<<<<<<< HEAD
	addPeersToChannel(t, n, channelName, gossips, peerIndexes)
	// Prime the membership view of the peers
	waitForFullMembershipOrFailNow(t, channelName, gossips, n, time.Second*30, time.Millisecond*100)

	grpcClient, err := comm.NewGRPCClient(comm.ClientConfig{})
	require.NoError(t, err)

	store := newTransientStore(t)
	defer store.tearDown()

	// Helper function that creates a gossipService instance
	newGossipService := func(i int) *GossipService {
		gs := gossips[i].GossipService
		gs.deliveryFactory = &embeddingDeliveryServiceFactory{&deliveryFactoryImpl{
			credentialSupport: comm.NewCredentialSupport(),
			deliverServiceConfig: &deliverservice.DeliverServiceConfig{
				PeerTLSEnabled:              false,
				ReConnectBackoffThreshold:   deliverservice.DefaultReConnectBackoffThreshold,
				ReconnectTotalTimeThreshold: time.Second,
				ConnectionTimeout:           time.Millisecond * 100,
			},
			deliverGRPCClient: grpcClient,
		}}
		gs.InitializeChannel(channelName, orderers.NewConnectionSource(flogging.MustGetLogger("peer.orderers"), nil), store.Store, Support{
			Committer: &mockLedgerInfo{1},
		})
		return gs
	}

	// The first leader is determined by the peer with the lower PKIid (lower TCP port in this case).
	// We set p0 to be the peer with the lower PKIid to ensure it'll be elected as leader before p1 and spare time.
	pkiID0 := gossips[0].peerIdentity
	pkiID1 := gossips[1].peerIdentity
	var firstLeaderIdx, secondLeaderIdx int
	if bytes.Compare(pkiID0, pkiID1) < 0 {
		firstLeaderIdx = 0
		secondLeaderIdx = 1
	} else {
		firstLeaderIdx = 1
		secondLeaderIdx = 0
	}
	p0 := newGossipService(firstLeaderIdx)
	p1 := newGossipService(secondLeaderIdx)

	// Returns index of the leader or -1 if no leader elected
	getLeader := func() int {
		p0.lock.RLock()
		p1.lock.RLock()
		defer p0.lock.RUnlock()
		defer p1.lock.RUnlock()

		if p0.leaderElection[channelName].IsLeader() {
=======
	addPeersToChannel(t, n, portPrefix, channelName, gossips, peerIndexes)
	// Prime the membership view of the peers
	waitForFullMembership(t, gossips, n, time.Second*30, time.Second*2)
	mcs := &naiveCryptoService{}
	// Helper function that creates a gossipService instance
	newGossipService := func(i int) *gossipServiceImpl {
		peerIdentity := api.PeerIdentityType(fmt.Sprintf("localhost:%d", portPrefix+i))
		gs := &gossipServiceImpl{
			mcs:             mcs,
			gossipSvc:       gossips[i],
			chains:          make(map[string]state.GossipStateProvider),
			leaderElection:  make(map[string]election.LeaderElectionService),
			deliveryFactory: &embeddingDeliveryServiceFactory{&deliveryFactoryImpl{}},
			idMapper:        identity.NewIdentityMapper(mcs, peerIdentity),
			peerIdentity:    peerIdentity,
			secAdv:          &secAdvMock{},
		}
		gossipServiceInstance = gs
		gs.InitializeChannel(channelName, &mockLedgerInfo{1}, []string{"localhost:7050"})
		return gs
	}

	p0 := newGossipService(0)
	p1 := newGossipService(1)

	// Returns index of the leader or -1 if no leader elected
	getLeader := func() int {
		if p0.leaderElection[channelName].IsLeader() {
			// Ensure p1 isn't a leader at the same time
			assert.False(t, p1.leaderElection[channelName].IsLeader())
>>>>>>> release-1.0
			return 0
		}
		if p1.leaderElection[channelName].IsLeader() {
			return 1
		}
		return -1
	}

<<<<<<< HEAD
	ds0 := p0.deliveryService[channelName].(*embeddingDeliveryService)
=======
	ds0 := p0.deliveryService.(*embeddingDeliveryService)
	ds1 := p1.deliveryService.(*embeddingDeliveryService)
>>>>>>> release-1.0

	// Wait for p0 to connect to the ordering service
	ds0.waitForDeliveryServiceActivation()
	t.Log("p0 started its delivery service")
	// Ensure it's a leader
	assert.Equal(t, 0, getLeader())
	// Wait for p0 to lose its leadership
	ds0.waitForDeliveryServiceTermination()
	t.Log("p0 stopped its delivery service")
<<<<<<< HEAD
	// Ensure p0 is not a leader
	assert.NotEqual(t, 0, getLeader())
	// Wait for p1 to take over. It should take over before time reaches timeLimit
	timeLimit := time.Now().Add(takeOverMaxTimeout)
	for getLeader() != 1 && time.Now().Before(timeLimit) {
		time.Sleep(100 * time.Millisecond)
	}
	if time.Now().After(timeLimit) && getLeader() != 1 {
		util.PrintStackTrace()
		t.Fatalf("p1 hasn't taken over leadership within %v: %d", takeOverMaxTimeout, getLeader())
	}
	t.Log("p1 has taken over leadership")
	p0.chains[channelName].Stop()
	p1.chains[channelName].Stop()
	p0.deliveryService[channelName].Stop()
	p1.deliveryService[channelName].Stop()
=======
	// Ensure there is no leader
	assert.Equal(t, -1, getLeader())
	// Wait for p1 to take over
	ds1.waitForDeliveryServiceActivation()
	t.Log("p1 started its delivery service")
	// Ensure it's a leader now
	assert.Equal(t, 1, getLeader())
	p0.chains[channelName].Stop()
	p1.chains[channelName].Stop()
	p0.deliveryService.Stop()
	p1.deliveryService.Stop()
>>>>>>> release-1.0
}
