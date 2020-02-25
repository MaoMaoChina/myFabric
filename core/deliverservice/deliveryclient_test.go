/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliverservice

import (
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/comm"
	"github.com/hyperledger/fabric/core/deliverservice/fake"
	"github.com/hyperledger/fabric/internal/pkg/peer/blocksprovider"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate counterfeiter -o fake/ledger_info.go --fake-name LedgerInfo . ledgerInfo
type ledgerInfo interface {
	blocksprovider.LedgerInfo
}

func TestStartDeliverForChannel(t *testing.T) {
	fakeLedgerInfo := &fake.LedgerInfo{}
	fakeLedgerInfo.LedgerHeightReturns(0, fmt.Errorf("fake-ledger-error"))

	grpcClient, err := comm.NewGRPCClient(comm.ClientConfig{
		SecOpts: comm.SecureOptions{
			UseTLS:            true,
			RequireClientCert: true,
			// The below certificates were taken from the peer TLS
			// dir as output by cryptogen.
			// They are server.crt and server.key respectively.
			Certificate: []byte(`-----BEGIN CERTIFICATE-----
MIIChTCCAiygAwIBAgIQOrr7/tDzKhhCba04E6QVWzAKBggqhkjOPQQDAjB2MQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEfMB0GA1UEAxMWdGxz
Y2Eub3JnMS5leGFtcGxlLmNvbTAeFw0xOTA4MjcyMDA2MDBaFw0yOTA4MjQyMDA2
MDBaMFsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH
Ew1TYW4gRnJhbmNpc2NvMR8wHQYDVQQDExZwZWVyMC5vcmcxLmV4YW1wbGUuY29t
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAExglppLxiAYSasrdFsrZJDxRULGBb
wHlArrap9SmAzGIeeIuqe9t3F23Q5Jry9lAnIh8h3UlkvZZpClXcjRiCeqOBtjCB
szAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMC
MAwGA1UdEwEB/wQCMAAwKwYDVR0jBCQwIoAgL35aqafj6SNnWdI4aMLh+oaFJvsA
aoHgYMkcPvvkiWcwRwYDVR0RBEAwPoIWcGVlcjAub3JnMS5leGFtcGxlLmNvbYIF
cGVlcjCCFnBlZXIwLm9yZzEuZXhhbXBsZS5jb22CBXBlZXIwMAoGCCqGSM49BAMC
A0cAMEQCIAiAGoYeKPMd3bqtixZji8q2zGzLmIzq83xdTJoZqm50AiAKleso2EVi
2TwsekWGpMaCOI6JV1+ZONyti6vBChhUYg==
-----END CERTIFICATE-----`),
			Key: []byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgxiyAFyD0Eg1NxjbS
U2EKDLoTQr3WPK8z7WyeOSzr+GGhRANCAATGCWmkvGIBhJqyt0WytkkPFFQsYFvA
eUCutqn1KYDMYh54i6p723cXbdDkmvL2UCciHyHdSWS9lmkKVdyNGIJ6
-----END PRIVATE KEY-----`,
			),
		},
	})
<<<<<<< HEAD
	require.NoError(t, err)

	t.Run("Green Path With Mutual TLS", func(t *testing.T) {
		ds := NewDeliverService(&Config{
			DeliverGRPCClient:    grpcClient,
			DeliverServiceConfig: &DeliverServiceConfig{},
		}).(*deliverServiceImpl)

		finalized := make(chan struct{})
		err := ds.StartDeliverForChannel("channel-id", fakeLedgerInfo, func() {
			close(finalized)
		})
		require.NoError(t, err)

		select {
		case <-finalized:
		case <-time.After(time.Second):
			assert.FailNow(t, "finalizer should have executed")
		}
=======
	assert.NoError(t, err)
	assert.NoError(t, service.StartDeliverForChannel("TEST_CHAINID", &mocks.MockLedgerInfo{0}, func() {}))

	// Lets start deliver twice
	assert.Error(t, service.StartDeliverForChannel("TEST_CHAINID", &mocks.MockLedgerInfo{0}, func() {}), "can't start delivery")
	// Lets stop deliver that not started
	assert.Error(t, service.StopDeliverForChannel("TEST_CHAINID2"), "can't stop delivery")

	// Let it try to simulate a few recv -> gossip rounds
	time.Sleep(time.Second)
	assert.NoError(t, service.StopDeliverForChannel("TEST_CHAINID"))
	time.Sleep(time.Duration(10) * time.Millisecond)
	// Make sure to stop all blocks providers
	service.Stop()
	time.Sleep(time.Duration(500) * time.Millisecond)
	assert.Equal(t, 0, connNumber)
	assertBlockDissemination(0, gossipServiceAdapter.GossipBlockDisseminations, t)
	assert.Equal(t, atomic.LoadInt32(&blocksDeliverer.RecvCnt), atomic.LoadInt32(&gossipServiceAdapter.AddPayloadsCnt))
	assert.Error(t, service.StartDeliverForChannel("TEST_CHAINID", &mocks.MockLedgerInfo{0}, func() {}), "Delivery service is stopping")
	assert.Error(t, service.StopDeliverForChannel("TEST_CHAINID"), "Delivery service is stopping")
}

func TestDeliverServiceRestart(t *testing.T) {
	defer ensureNoGoroutineLeak(t)()
	// Scenario: bring up ordering service instance, then shut it down, and then resurrect it.
	// Client is expected to reconnect to it, and to ask for a block sequence that is the next block
	// after the last block it got from the previous incarnation of the ordering service.

	os := mocks.NewOrderer(5611, t)

	time.Sleep(time.Second)
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}
>>>>>>> release-1.0

		bp, ok := ds.blockProviders["channel-id"]
		require.True(t, ok, "map entry must exist")
		assert.Equal(t, "76f7a03f8dfdb0ef7c4b28b3901fe163c730e906c70e4cdf887054ad5f608bed", fmt.Sprintf("%x", bp.TLSCertHash))
	})

<<<<<<< HEAD
	t.Run("Green Path without mutual TLS", func(t *testing.T) {
		grpcClient, err := comm.NewGRPCClient(comm.ClientConfig{
			SecOpts: comm.SecureOptions{
				UseTLS: true,
			},
		})
		require.NoError(t, err)

		ds := NewDeliverService(&Config{
			DeliverGRPCClient:    grpcClient,
			DeliverServiceConfig: &DeliverServiceConfig{},
		}).(*deliverServiceImpl)

		finalized := make(chan struct{})
		err = ds.StartDeliverForChannel("channel-id", fakeLedgerInfo, func() {
			close(finalized)
		})
		require.NoError(t, err)

		select {
		case <-finalized:
		case <-time.After(time.Second):
			assert.FailNow(t, "finalizer should have executed")
		}
=======
	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")
	// Check that delivery client requests blocks in order
	go os.SendBlock(uint64(100))
	assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)
	go os.SendBlock(uint64(101))
	assertBlockDissemination(101, gossipServiceAdapter.GossipBlockDisseminations, t)
	go os.SendBlock(uint64(102))
	assertBlockDissemination(102, gossipServiceAdapter.GossipBlockDisseminations, t)
	os.Shutdown()
	time.Sleep(time.Second * 3)
	os = mocks.NewOrderer(5611, t)
	li.Height = 103
	os.SetNextExpectedSeek(uint64(103))
	go os.SendBlock(uint64(103))
	assertBlockDissemination(103, gossipServiceAdapter.GossipBlockDisseminations, t)
	service.Stop()
	os.Shutdown()
}

func TestDeliverServiceFailover(t *testing.T) {
	defer ensureNoGoroutineLeak(t)()
	// Scenario: bring up 2 ordering service instances,
	// and shut down the instance that the client has connected to.
	// Client is expected to connect to the other instance, and to ask for a block sequence that is the next block
	// after the last block it got from the ordering service that was shut down.
	// Then, shut down the other node, and bring back the first (that was shut down first).

	os1 := mocks.NewOrderer(5612, t)
	os2 := mocks.NewOrderer(5613, t)

	time.Sleep(time.Second)
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}
>>>>>>> release-1.0

		bp, ok := ds.blockProviders["channel-id"]
		require.True(t, ok, "map entry must exist")
		assert.Nil(t, bp.TLSCertHash)
	})

<<<<<<< HEAD
	t.Run("Exists", func(t *testing.T) {
		ds := NewDeliverService(&Config{
			DeliverGRPCClient:    grpcClient,
			DeliverServiceConfig: &DeliverServiceConfig{},
		}).(*deliverServiceImpl)
=======
	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")
	// We need to discover to which instance the client connected to
	go os1.SendBlock(uint64(100))
	instance2fail := os1
	reincarnatedNodePort := 5612
	instance2failSecond := os2
	select {
	case seq := <-gossipServiceAdapter.GossipBlockDisseminations:
		assert.Equal(t, uint64(100), seq)
	case <-time.After(time.Second * 2):
		// Shutdown first instance and replace it, in order to make an instance
		// with an empty sending channel
		os1.Shutdown()
		time.Sleep(time.Second)
		os1 = mocks.NewOrderer(5612, t)
		instance2fail = os2
		instance2failSecond = os1
		reincarnatedNodePort = 5613
		// Ensure we really are connected to the second instance,
		// by making it send a block
		go os2.SendBlock(uint64(100))
		assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)
	}
>>>>>>> release-1.0

		err = ds.StartDeliverForChannel("channel-id", fakeLedgerInfo, func() {})
		require.NoError(t, err)

		err = ds.StartDeliverForChannel("channel-id", fakeLedgerInfo, func() {})
		assert.EqualError(t, err, "Delivery service - block provider already exists for channel-id found, can't start delivery")
	})
<<<<<<< HEAD
=======
	assert.NoError(t, err)
	li := &mocks.MockLedgerInfo{Height: 100}
	os1.SetNextExpectedSeek(li.Height)
	os2.SetNextExpectedSeek(li.Height)

	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")
>>>>>>> release-1.0

	t.Run("Stopping", func(t *testing.T) {
		ds := NewDeliverService(&Config{
			DeliverGRPCClient:    grpcClient,
			DeliverServiceConfig: &DeliverServiceConfig{},
		}).(*deliverServiceImpl)

		ds.Stop()

		err = ds.StartDeliverForChannel("channel-id", fakeLedgerInfo, func() {})
		assert.EqualError(t, err, "Delivery service is stopping cannot join a new channel channel-id")
	})
<<<<<<< HEAD
}

func TestStopDeliverForChannel(t *testing.T) {
	t.Run("Green path", func(t *testing.T) {
		ds := NewDeliverService(&Config{}).(*deliverServiceImpl)
		doneA := make(chan struct{})
		ds.blockProviders = map[string]*blocksprovider.Deliverer{
			"a": {
				DoneC: doneA,
			},
			"b": {
				DoneC: make(chan struct{}),
			},
		}
		err := ds.StopDeliverForChannel("a")
		assert.NoError(t, err)
		assert.Len(t, ds.blockProviders, 1)
		_, ok := ds.blockProviders["a"]
		assert.False(t, ok)
		select {
		case <-doneA:
		default:
			assert.Fail(t, "should have stopped the blocksprovider")
		}
=======
	assert.NoError(t, err)

	li := &mocks.MockLedgerInfo{Height: uint64(100)}
	os.SetNextExpectedSeek(uint64(100))
	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")

	// Check that delivery service requests blocks in order
	go os.SendBlock(uint64(100))
	assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)
	go os.SendBlock(uint64(101))
	assertBlockDissemination(101, gossipServiceAdapter.GossipBlockDisseminations, t)
	atomic.StoreUint64(&li.Height, uint64(102))
	os.SetNextExpectedSeek(uint64(102))
	// Now stop the delivery service and make sure we don't disseminate a block
	service.Stop()
	go os.SendBlock(uint64(102))
	select {
	case <-gossipServiceAdapter.GossipBlockDisseminations:
		assert.Fail(t, "Disseminated a block after shutting down the delivery service")
	case <-time.After(time.Second * 2):
	}
	os.Shutdown()
	time.Sleep(time.Second)
}

func TestDeliverServiceShutdownRespawn(t *testing.T) {
	// Scenario: Launch an ordering service node and let the client pull some blocks.
	// Then, wait a few seconds, and don't send any blocks.
	// Afterwards - start a new instance and shut down the old instance.
	SetReconnectTotalTimeThreshold(time.Second)
	defer func() {
		SetReconnectTotalTimeThreshold(time.Second * 60 * 5)
	}()
	defer ensureNoGoroutineLeak(t)()

	osn1 := mocks.NewOrderer(5614, t)

	time.Sleep(time.Second)
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}

	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"localhost:5614", "localhost:5615"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
	})
	assert.NoError(t, err)

	li := &mocks.MockLedgerInfo{Height: uint64(100)}
	osn1.SetNextExpectedSeek(uint64(100))
	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")

	// Check that delivery service requests blocks in order
	go osn1.SendBlock(uint64(100))
	assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)
	go osn1.SendBlock(uint64(101))
	assertBlockDissemination(101, gossipServiceAdapter.GossipBlockDisseminations, t)
	atomic.StoreUint64(&li.Height, uint64(102))
	// Now wait for a few seconds
	time.Sleep(time.Second * 2)
	// Now start the new instance
	osn2 := mocks.NewOrderer(5615, t)
	// Now stop the old instance
	osn1.Shutdown()
	// Send a block from osn2
	osn2.SetNextExpectedSeek(uint64(102))
	go osn2.SendBlock(uint64(102))
	// Ensure it is received
	assertBlockDissemination(102, gossipServiceAdapter.GossipBlockDisseminations, t)
	service.Stop()
	osn2.Shutdown()
}

func TestDeliverServiceBadConfig(t *testing.T) {
	// Empty endpoints
	service, err := NewDeliverService(&Config{
		Endpoints:   []string{},
		Gossip:      &mocks.MockGossipServiceAdapter{},
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
>>>>>>> release-1.0
	})

	t.Run("Already stopping", func(t *testing.T) {
		ds := NewDeliverService(&Config{}).(*deliverServiceImpl)
		ds.blockProviders = map[string]*blocksprovider.Deliverer{
			"a": {
				DoneC: make(chan struct{}),
			},
			"b": {
				DoneC: make(chan struct{}),
			},
		}

		ds.Stop()
		err := ds.StopDeliverForChannel("a")
		assert.EqualError(t, err, "Delivery service is stopping, cannot stop delivery for channel a")
	})

	t.Run("Non-existent", func(t *testing.T) {
		ds := NewDeliverService(&Config{}).(*deliverServiceImpl)
		ds.blockProviders = map[string]*blocksprovider.Deliverer{
			"a": {
				DoneC: make(chan struct{}),
			},
			"b": {
				DoneC: make(chan struct{}),
			},
		}

		err := ds.StopDeliverForChannel("c")
		assert.EqualError(t, err, "Delivery service - no block provider for c found, can't stop delivery")
	})
}

func TestStop(t *testing.T) {
	ds := NewDeliverService(&Config{}).(*deliverServiceImpl)
	ds.blockProviders = map[string]*blocksprovider.Deliverer{
		"a": {
			DoneC: make(chan struct{}),
		},
		"b": {
			DoneC: make(chan struct{}),
		},
	}
	assert.False(t, ds.stopping)
	for _, bp := range ds.blockProviders {
		select {
		case <-bp.DoneC:
			assert.Fail(t, "block providers should not be closed")
		default:
		}
	}

	ds.Stop()
	assert.True(t, ds.stopping)
	assert.Len(t, ds.blockProviders, 2)
	for _, bp := range ds.blockProviders {
		select {
		case <-bp.DoneC:
		default:
			assert.Fail(t, "block providers should te closed")
		}
	}

}
