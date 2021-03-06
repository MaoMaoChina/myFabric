/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package multichain

import (
	"fmt"

	"github.com/hyperledger/fabric/common/config"
	"github.com/hyperledger/fabric/common/configtx"
	configtxapi "github.com/hyperledger/fabric/common/configtx/api"
	"github.com/hyperledger/fabric/common/policies"
	"github.com/hyperledger/fabric/orderer/ledger"
	cb "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/op/go-logging"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/crypto"
)

var logger = logging.MustGetLogger("orderer/multichain")

const (
	msgVersion = int32(0)
	epoch      = 0
)

// Manager coordinates the creation and access of chains
type Manager interface {
	// GetChain retrieves the chain support for a chain (and whether it exists)
	// 获取链对象
	GetChain(chainID string) (ChainSupport, bool)

	// SystemChannelID returns the channel ID for the system channel
	// 获取系统链,系统链用于引导其他链的生成
	SystemChannelID() string

	// NewChannelConfig returns a bare bones configuration ready for channel
	// creation request to be applied on top of it
	// 链的配置(New/Update)
	NewChannelConfig(envConfigUpdate *cb.Envelope) (configtxapi.Manager, error)
}

// 配置资源
// 每一个链都有自己的配置,配置被包裹在配置资源中
type configResources struct {
	configtxapi.Manager
}

func (cr *configResources) SharedConfig() config.Orderer {
	oc, ok := cr.OrdererConfig()
	if !ok {
		logger.Panicf("[channel %s] has no orderer configuration", cr.ChainID())
	}
	return oc
}
// 账本资源类 (配置资源&账本的读写对象)
type ledgerResources struct {
	*configResources
	ledger ledger.ReadWriter
}
// manager的实现类
type multiLedger struct {
	chains          map[string]*chainSupport  // 多链对象
	lock		sync.Mutex		  // 加锁*******************************************************
	consenters      map[string]Consenter      // 共识机制
	ledgerFactory   ledger.Factory            // 账本读写工厂
	signer          crypto.LocalSigner        // 签名
	systemChannelID string		          // 系统链名
	systemChannel   *chainSupport		  // 系统链对象
}
// 获取链最新配置交易
func getConfigTx(reader ledger.Reader) *cb.Envelope {
	lastBlock := ledger.GetBlock(reader, reader.Height()-1) // 获取链最新生成的块
	index, err := utils.GetLastConfigIndexFromBlock(lastBlock)
	if err != nil {
		logger.Panicf("Chain did not have appropriately encoded last config in its latest block: %s", err)
	}
	configBlock := ledger.GetBlock(reader, index)// 通过最新生成的块来获取配置块
	if configBlock == nil {
		logger.Panicf("Config block does not exist")
	}

	return utils.ExtractEnvelopeOrPanic(configBlock, 0)
}

// NewManagerImpl produces an instance of a Manager
func NewManagerImpl(ledgerFactory ledger.Factory, consenters map[string]Consenter, signer crypto.LocalSigner) Manager {
	ml := &multiLedger{
		chains:        make(map[string]*chainSupport),
		ledgerFactory: ledgerFactory,
		consenters:    consenters,
		signer:        signer,
	}
	// 读取本地存储的链的ID
	existingChains := ledgerFactory.ChainIDs()
	for _, chainID := range existingChains {
		// 实例化账本读对象 rl = read ledger
		rl, err := ledgerFactory.GetOrCreate(chainID)
		if err != nil {
			logger.Panicf("Ledger factory reported chainID %s but could not retrieve it: %s", chainID, err)
		}
		// 获取链最新配置交易
		configTx := getConfigTx(rl)
		if configTx == nil {
			logger.Panic("Programming error, configTx should never be nil here")
		}
		// 绑定配置与读写对象
		ledgerResources := ml.newLedgerResources(configTx)
		chainID := ledgerResources.ChainID()
		// 检查该链是否为系统链,系统链有联盟配置,联盟配置有创建其它链的权限
		if _, ok := ledgerResources.ConsortiumsConfig(); ok {
			// 检查系统链是否已经存在
			if ml.systemChannelID != "" {
				logger.Panicf("There appear to be two system chains %s and %s", ml.systemChannelID, chainID)
			}
			chain := newChainSupport(createSystemChainFilters(ml, ledgerResources),
				ledgerResources,
				consenters,
				signer)
			logger.Infof("Starting with system channel %s and orderer type %s", chainID, chain.SharedConfig().ConsensusType())
			ml.lock.Lock()
			ml.chains[chainID] = chain
			ml.lock.Unlock()
			ml.systemChannelID = chainID
			ml.systemChannel = chain
			// We delay starting this chain, as it might try to copy and replace the chains map via newChain before the map is fully built
			defer chain.start() // 延迟启动
		} else {
			logger.Debugf("Starting chain: %s", chainID)
			chain := newChainSupport(createStandardFilters(ledgerResources),
				ledgerResources,
				consenters,
				signer)
			ml.lock.Lock()
			ml.chains[chainID] = chain
			ml.lock.Unlock()
			chain.start() // 启动
		}
	//
	}
	// 检查是否有系统链存在,不存在则报错 
	if ml.systemChannelID == "" {
		logger.Panicf("No system chain found.  If bootstrapping, does your system channel contain a consortiums group definition?")
	}

	return ml
}

func (ml *multiLedger) SystemChannelID() string {
	return ml.systemChannelID
}

// GetChain retrieves the chain support for a chain (and whether it exists)
func (ml *multiLedger) GetChain(chainID string) (ChainSupport, bool) {
	ml.lock.Lock()
	defer ml.lock.Unlock()
	cs, ok := ml.chains[chainID]
	return cs, ok
}

// 实例化账本资源对象
func (ml *multiLedger) newLedgerResources(configTx *cb.Envelope) *ledgerResources {
	initializer := configtx.NewInitializer()
	configManager, err := configtx.NewManagerImpl(configTx, initializer, nil)
	if err != nil {
		logger.Panicf("Error creating configtx manager and handlers: %s", err)
	}

	chainID := configManager.ChainID()

	ledger, err := ml.ledgerFactory.GetOrCreate(chainID)
	if err != nil {
		logger.Panicf("Error getting ledger for %s", chainID)
	}

	return &ledgerResources{
		configResources: &configResources{Manager: configManager},
		ledger:          ledger,
	}
}

// 新建一条链
func (ml *multiLedger) newChain(configtx *cb.Envelope) {
	ledgerResources := ml.newLedgerResources(configtx)
	ledgerResources.ledger.Append(ledger.CreateNextBlock(ledgerResources.ledger, []*cb.Envelope{configtx}))  // 将创世区块加在链上

	// Copy the map to allow concurrent reads from broadcast/deliver while the new chainSupport is
	// 使用Mutex加锁后就不再需要了
	//newChains := make(map[string]*chainSupport)
	//for key, value := range ml.chains {
	//	newChains[key] = value
	//}

	cs := newChainSupport(createStandardFilters(ledgerResources), ledgerResources, ml.consenters, ml.signer)
	chainID := ledgerResources.ChainID()

	logger.Infof("Created and starting new chain %s", chainID)

	newChains[string(chainID)] = cs
	cs.start()

	// ml.chains = newChains
	ml.lock.Lock()
	ml.chains[chainID] = cs
	ml.lock.Unlock()
}

func (ml *multiLedger) channelsCount() int {
	return len(ml.chains)
}

// 生成新链的配置
func (ml *multiLedger) NewChannelConfig(envConfigUpdate *cb.Envelope) (configtxapi.Manager, error) {
	configUpdatePayload, err := utils.UnmarshalPayload(envConfigUpdate.Payload)
	if err != nil {
		return nil, fmt.Errorf("Failing initial channel config creation because of payload unmarshaling error: %s", err)
	}

	configUpdateEnv, err := configtx.UnmarshalConfigUpdateEnvelope(configUpdatePayload.Data)
	if err != nil {
		return nil, fmt.Errorf("Failing initial channel config creation because of config update envelope unmarshaling error: %s", err)
	}

	if configUpdatePayload.Header == nil {
		return nil, fmt.Errorf("Failed initial channel config creation because config update header was missing")
	}
	channelHeader, err := utils.UnmarshalChannelHeader(configUpdatePayload.Header.ChannelHeader)

	configUpdate, err := configtx.UnmarshalConfigUpdate(configUpdateEnv.ConfigUpdate)
	if err != nil {
		return nil, fmt.Errorf("Failing initial channel config creation because of config update unmarshaling error: %s", err)
	}

	if configUpdate.ChannelId != channelHeader.ChannelId {
		return nil, fmt.Errorf("Failing initial channel config creation: mismatched channel IDs: '%s' != '%s'", configUpdate.ChannelId, channelHeader.ChannelId)
	}

	if configUpdate.WriteSet == nil {
		return nil, fmt.Errorf("Config update has an empty writeset")
	}

	if configUpdate.WriteSet.Groups == nil || configUpdate.WriteSet.Groups[config.ApplicationGroupKey] == nil {
		return nil, fmt.Errorf("Config update has missing application group")
	}

	if uv := configUpdate.WriteSet.Groups[config.ApplicationGroupKey].Version; uv != 1 {
		return nil, fmt.Errorf("Config update for channel creation does not set application group version to 1, was %d", uv)
	}

	consortiumConfigValue, ok := configUpdate.WriteSet.Values[config.ConsortiumKey]
	if !ok {
		return nil, fmt.Errorf("Consortium config value missing")
	}

	consortium := &cb.Consortium{}
	err = proto.Unmarshal(consortiumConfigValue.Value, consortium)
	if err != nil {
		return nil, fmt.Errorf("Error reading unmarshaling consortium name: %s", err)
	}

	applicationGroup := cb.NewConfigGroup()
	consortiumsConfig, ok := ml.systemChannel.ConsortiumsConfig()
	if !ok {
		return nil, fmt.Errorf("The ordering system channel does not appear to support creating channels")
	}

	consortiumConf, ok := consortiumsConfig.Consortiums()[consortium.Name]
	if !ok {
		return nil, fmt.Errorf("Unknown consortium name: %s", consortium.Name)
	}

	applicationGroup.Policies[config.ChannelCreationPolicyKey] = &cb.ConfigPolicy{
		Policy: consortiumConf.ChannelCreationPolicy(),
	}
	applicationGroup.ModPolicy = config.ChannelCreationPolicyKey

	// Get the current system channel config
	systemChannelGroup := ml.systemChannel.ConfigEnvelope().Config.ChannelGroup

	// If the consortium group has no members, allow the source request to have no members.  However,
	// if the consortium group has any members, there must be at least one member in the source request
	if len(systemChannelGroup.Groups[config.ConsortiumsGroupKey].Groups[consortium.Name].Groups) > 0 &&
		len(configUpdate.WriteSet.Groups[config.ApplicationGroupKey].Groups) == 0 {
		return nil, fmt.Errorf("Proposed configuration has no application group members, but consortium contains members")
	}

	// If the consortium has no members, allow the source request to contain arbitrary members
	// Otherwise, require that the supplied members are a subset of the consortium members
	if len(systemChannelGroup.Groups[config.ConsortiumsGroupKey].Groups[consortium.Name].Groups) > 0 {
		for orgName := range configUpdate.WriteSet.Groups[config.ApplicationGroupKey].Groups {
			consortiumGroup, ok := systemChannelGroup.Groups[config.ConsortiumsGroupKey].Groups[consortium.Name].Groups[orgName]
			if !ok {
				return nil, fmt.Errorf("Attempted to include a member which is not in the consortium")
			}
			applicationGroup.Groups[orgName] = consortiumGroup
		}
	}

	channelGroup := cb.NewConfigGroup()

	// Copy the system channel Channel level config to the new config
	for key, value := range systemChannelGroup.Values {
		channelGroup.Values[key] = value
		if key == config.ConsortiumKey {
			// Do not set the consortium name, we do this later
			continue
		}
	}

	for key, policy := range systemChannelGroup.Policies {
		channelGroup.Policies[key] = policy
	}

	// Set the new config orderer group to the system channel orderer group and the application group to the new application group
	channelGroup.Groups[config.OrdererGroupKey] = systemChannelGroup.Groups[config.OrdererGroupKey]
	channelGroup.Groups[config.ApplicationGroupKey] = applicationGroup
	channelGroup.Values[config.ConsortiumKey] = config.TemplateConsortium(consortium.Name).Values[config.ConsortiumKey]

	templateConfig, _ := utils.CreateSignedEnvelope(cb.HeaderType_CONFIG, configUpdate.ChannelId, ml.signer, &cb.ConfigEnvelope{
		Config: &cb.Config{
			ChannelGroup: channelGroup,
		},
	}, msgVersion, epoch)

	initializer := configtx.NewInitializer()

	// This is a very hacky way to disable the sanity check logging in the policy manager
	// for the template configuration, but it is the least invasive near a release
	pm, ok := initializer.PolicyManager().(*policies.ManagerImpl)
	if ok {
		pm.SuppressSanityLogMessages = true
		defer func() {
			pm.SuppressSanityLogMessages = false
		}()
	}
	// 配置manager实例化
	return configtx.NewManagerImpl(templateConfig, initializer, nil)
}
