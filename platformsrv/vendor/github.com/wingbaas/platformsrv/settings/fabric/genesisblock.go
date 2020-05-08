

package fabric

import (
	"fmt"
	"time"
	"github.com/hyperledger/fabric-protos-go/orderer/etcdraft"
	genesisconfig "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	genesisconfigraft "github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/genesisconfig"
	"github.com/wingbaas/platformsrv/certgenerate/fabric/gopkg.in/yaml.v2"
	"github.com/wingbaas/platformsrv/logger" 
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric/txgenerate"
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
)

func GenerateGenesisBlockRaft(serviceRootPath string,deployConf public.DeployNetConfig) error { 
	var txCfg genesisconfigraft.TopLevel
	var orderOrgArray []genesisconfigraft.Organization
	for _,org := range deployConf.OrdererOrgs {
		var og genesisconfigraft.Organization
		og.Name = org.Name + "MSP"
		og.ID = org.Name + "MSP"
		og.MSPDir = serviceRootPath + "/crypto-config/ordererOrganizations/" + org.Domain + "/msp"

		var ogPolicy map[string]*genesisconfigraft.Policy
		ogPolicy = make(map[string]*genesisconfigraft.Policy)

		var policyR genesisconfigraft.Policy
		var policyW genesisconfigraft.Policy
		var policyA genesisconfigraft.Policy
		policyR.Type = "Signature"
		policyR.Rule = "OR('" + og.ID + ".member')"
		policyW.Type = "Signature"
		policyW.Rule = "OR('" + og.ID + ".member')"
		policyA.Type = "Signature"
		policyA.Rule = "OR('" + og.ID + ".admin')"
		ogPolicy["Readers"] = &policyR
		ogPolicy["Writers"] = &policyW
		ogPolicy["Admins"] = &policyA
		og.Policies = ogPolicy
		txCfg.Organizations = append(txCfg.Organizations,&og)
		orderOrgArray = append(orderOrgArray,og)
	}
	var orgArray []genesisconfigraft.Organization
	for _,org := range deployConf.PeerOrgs {
		var og genesisconfigraft.Organization
		og.Name = org.Name + "MSP"
		og.ID = org.Name + "MSP"
		og.MSPDir = serviceRootPath + "/crypto-config/peerOrganizations/" + org.Domain + "/msp"
		og.Policies = make(map[string]*genesisconfigraft.Policy)

		var ogPolicy map[string]*genesisconfigraft.Policy
		ogPolicy = make(map[string]*genesisconfigraft.Policy)

		var policyR genesisconfigraft.Policy
		var policyW genesisconfigraft.Policy
		var policyA genesisconfigraft.Policy
		var policyE genesisconfigraft.Policy
		policyR.Type = "Signature"
		policyR.Rule = "OR('" + og.ID + ".admin', " + "'" + og.ID + ".peer', " + "'" + og.ID + ".client')"
		policyW.Type = "Signature"
		policyW.Rule = "OR('" + og.ID + ".admin', " + "'" + og.ID + ".client')"
		policyA.Type = "Signature"
		policyA.Rule = "OR('" + og.ID + ".admin')"
		policyE.Type = "Signature"
		policyE.Rule = "OR('" + og.ID + ".peer')"
		ogPolicy["Readers"] = &policyR
		ogPolicy["Writers"] = &policyW
		ogPolicy["Admins"] = &policyA
		ogPolicy["Endorsement"] = &policyE
		og.Policies = ogPolicy

		var anchorPeer genesisconfigraft.AnchorPeer
		for _,h := range org.Specs {
			anchorPeer.Host = h.Hostname
			anchorPeer.Port = 7051
			break
		}
		og.AnchorPeers = append(og.AnchorPeers,&anchorPeer)
		txCfg.Organizations = append(txCfg.Organizations,&og)
		orgArray = append(orgArray,og)
	}
	capCh := make(map[string]bool)
	capCh["V2_0"] = true
	capOrder := make(map[string]bool)
	capOrder["V2_0"] = true
	capApp := make(map[string]bool)
	capApp["V2_0"] = true

	txCfg.Capabilities = make(map[string]map[string]bool)
	txCfg.Capabilities["Channel"] = capCh
	txCfg.Capabilities["Orderer"] = capOrder
	txCfg.Capabilities["Application"] = capApp

	var app genesisconfigraft.Application
	var policyAppR genesisconfigraft.Policy
	var policyAppW genesisconfigraft.Policy
	var policyAppA genesisconfigraft.Policy
	var policyAppL genesisconfigraft.Policy
	var policyAppE genesisconfigraft.Policy
	policyAppR.Type = "ImplicitMeta"
	policyAppR.Rule = "ANY Readers"
	policyAppW.Type = "ImplicitMeta"
	policyAppW.Rule = "ANY Writers"
	policyAppA.Type = "ImplicitMeta"
	policyAppA.Rule = "MAJORITY Admins"
	policyAppL.Type = "ImplicitMeta"
	policyAppL.Rule = "MAJORITY Endorsement"
	policyAppE.Type = "ImplicitMeta"
	policyAppE.Rule = "MAJORITY Endorsement"

	app.Policies = make(map[string]*genesisconfigraft.Policy)
	app.Policies["Readers"] = &policyAppR
	app.Policies["Writers"] = &policyAppW
	app.Policies["Admins"] = &policyAppA
	app.Policies["LifecycleEndorsement"] = &policyAppL
	app.Policies["Endorsement"] = &policyAppE
	app.Capabilities = make(map[string]bool)
	app.Capabilities = capApp

	txCfg.Application = &app

	var order genesisconfigraft.Orderer
	order.OrdererType = "etcdraft"
	for _,org := range deployConf.OrdererOrgs {
		for _,h := range org.Specs {
			order.Addresses = append(order.Addresses,h.Hostname/* + "." + org.Domain*/ + ":7050")
		}
	}
	order.BatchTimeout = 2
	order.BatchSize = genesisconfigraft.BatchSize{
		MaxMessageCount: 10,
		AbsoluteMaxBytes: 99*1024*1024,
		PreferredMaxBytes: 512*1024,
	}
	order.Capabilities = make(map[string]bool)
	order.Capabilities = capApp
	order.Policies = make(map[string]*genesisconfigraft.Policy)
	order.Policies["Readers"] = &policyAppR
	order.Policies["Writers"] = &policyAppW
	order.Policies["Admins"] = &policyAppA
	order.Policies["BlockValidation"] = &policyAppW
	var raft etcdraft.ConfigMetadata
	raft.Options = &etcdraft.Options{
		TickInterval: "500ms",
		ElectionTick: 10,
		HeartbeatTick: 1,
		MaxInflightBlocks: 5,
		SnapshotIntervalSize: 16*1024*1024,
	}
	for _,org := range deployConf.OrdererOrgs {
		for _,h := range org.Specs {
			var consent etcdraft.Consenter
			consent.Host = h.Hostname /*+ "." + org.Domain*/ 
			consent.Port = 7050
			certPath := serviceRootPath + "/crypto-config/ordererOrganizations/" + org.Domain + "/orderers/" + h.Hostname + "." + org.Domain + "/tls/server.crt"
			consent.ClientTlsCert = utils.ReadFileBytes(certPath)
			consent.ServerTlsCert = consent.ClientTlsCert
			raft.Consenters = append(raft.Consenters,&consent)
		}
	}
	order.EtcdRaft = &raft
	txCfg.Orderer = &order

	var chProfile genesisconfigraft.Profile
	chProfile.Policies = make(map[string]*genesisconfigraft.Policy)
	chProfile.Policies["Readers"] = &policyAppR
	chProfile.Policies["Writers"] = &policyAppW
	chProfile.Policies["Admins"] = &policyAppA
	chProfile.Capabilities = make(map[string]bool)
	chProfile.Capabilities = capApp
	txCfg.Channel = &chProfile

	var twoOrgsOrdererGenesisProfile genesisconfigraft.Profile
	twoOrgsOrdererGenesisProfile.Policies = make(map[string]*genesisconfigraft.Policy)
	twoOrgsOrdererGenesisProfile.Policies = chProfile.Policies
	twoOrgsOrdererGenesisProfile.Capabilities = make(map[string]bool)
	twoOrgsOrdererGenesisProfile.Capabilities = capApp
	for _,org := range orderOrgArray {
		order.Organizations = append(order.Organizations,&org)
	}
	twoOrgsOrdererGenesisProfile.Orderer = &order
	twoOrgsOrdererGenesisProfile.Consortiums = make(map[string]*genesisconfigraft.Consortium)
	var consortium genesisconfigraft.Consortium
	for _,org := range orgArray {
		consortium.Organizations = append(consortium.Organizations,&org)
	}
	twoOrgsOrdererGenesisProfile.Consortiums["SampleConsortium"] = &consortium
	txCfg.Profiles = make(map[string]*genesisconfigraft.Profile)
	txCfg.Profiles["TwoOrgsOrdererGenesis"] = &twoOrgsOrdererGenesisProfile 

	var twoOrgsChannelProfile genesisconfigraft.Profile
	twoOrgsChannelProfile.Consortium = "SampleConsortium"
	twoOrgsChannelProfile.Policies = make(map[string]*genesisconfigraft.Policy)
	twoOrgsChannelProfile.Policies = chProfile.Policies
	twoOrgsChannelProfile.Capabilities = make(map[string]bool)
	twoOrgsChannelProfile.Capabilities = capApp

	var profileApp genesisconfigraft.Application
	profileApp.Organizations = append(profileApp.Organizations,consortium.Organizations...)
	profileApp.Policies = make(map[string]*genesisconfigraft.Policy)
	profileApp.Policies = app.Policies
	profileApp.Capabilities = make(map[string]bool)
	profileApp.Capabilities = capApp
	twoOrgsChannelProfile.Application = &profileApp
	txCfg.Profiles["TwoOrgsChannel"] = &twoOrgsChannelProfile

	by,err := yaml.Marshal(txCfg)
	if err != nil {
		logger.Errorf("GenerateGenesisBlockRaft: marshal tx config failed,%s",err.Error())
		return fmt.Errorf("GenerateGenesisBlockRaft: marshal tx config failed,%s",err.Error())
	}
	cnfTxYaml := serviceRootPath + "/configtx.yaml"
	err = utils.WriteFile(cnfTxYaml,string(by))
	if err != nil {
		logger.Errorf("GenerateGenesisBlockRaft: write tx config to file failed,%s",err)
		return fmt.Errorf("GenerateGenesisBlockRaft: write tx config to file failed,%s",err)
	}
	txPath := serviceRootPath + "/channel-artifacts/"
	err = utils.DirCheck(txPath)
	if err != nil {
		logger.Errorf("GenerateGenesisBlockRaft: artifacts dir create error,%v",err)
		return fmt.Errorf("GenerateGenesisBlockRaft: artifacts dir create error,%v",err)
	}
	bl := txgeneratev2.CreateGenesisBlock(serviceRootPath, "testchainid", txPath)
	if !bl {
		logger.Errorf("GenerateGenesisBlockRaft: create blockchain genesis block failed")
		return fmt.Errorf("GenerateGenesisBlockRaft: create blockchain genesis block failed")
	}
	return nil
}

func GenerateGenesisBlock(serviceRootPath string,deployConf public.DeployNetConfig,deployType string) error { 
	var txCnf genesisconfig.TopLevel
	var orderArray []*genesisconfig.Organization
	var peerArray []*genesisconfig.Organization
	
	for i := 0; i < len(deployConf.OrdererOrgs); i++ { 
		var org genesisconfig.Organization
		org.ID = deployConf.OrdererOrgs[i].Name + "MSP"
		org.Name = deployConf.OrdererOrgs[i].Name + "MSP"
		org.MSPDir = serviceRootPath + "/crypto-config/ordererOrganizations/" + deployConf.OrdererOrgs[i].Domain + "/msp"
		txCnf.Organizations = append(txCnf.Organizations, &org)
		orderArray = append(orderArray, &org)
	}

	for i := 0; i < len(deployConf.PeerOrgs); i++ {
		var org genesisconfig.Organization
		org.ID = deployConf.PeerOrgs[i].Name + "MSP"
		org.Name = deployConf.PeerOrgs[i].Name + "MSP"
		org.MSPDir = serviceRootPath + "/crypto-config/peerOrganizations/" + deployConf.PeerOrgs[i].Domain + "/msp"
		if len(deployConf.PeerOrgs[i].Specs) > 0 {
			var anchor genesisconfig.AnchorPeer
			anchor.Host = deployConf.PeerOrgs[i].Specs[0].Hostname
			anchor.Port = 7051
			org.AnchorPeers = append(org.AnchorPeers, &anchor)
		}
		txCnf.Organizations = append(txCnf.Organizations, &org)
		peerArray = append(peerArray, &org)
	}

	var orders genesisconfig.Orderer
	if deployType == public.KAFKA_FABRIC {
		orders.OrdererType = "kafka"
	} else if deployType == public.SOLO_FABRIC {
		orders.OrdererType = "solo"
	} else if deployType == public.RAFT_FABRIC {
		orders.OrdererType = "etcdraft"
	}else {
		logger.Errorf("unsupported orderer type %s",deployType)
		return fmt.Errorf("unsupported orderer type %s",deployType)
	}
	for i := 0; i < len(deployConf.OrdererOrgs); i++ {
		if "Orderer" == deployConf.OrdererOrgs[i].Name {
			for j := 0; j < len(deployConf.OrdererOrgs[i].Specs); j++ {
				var tmpStr string
				tmpStr = deployConf.OrdererOrgs[i].Specs[j].Hostname + ":7050"
				orders.Addresses = append(orders.Addresses, tmpStr)
			}
		}
	}
	orders.BatchTimeout = 2 * time.Second
	orders.BatchSize.AbsoluteMaxBytes = 10 * 1024 * 1024
	orders.BatchSize.MaxMessageCount = 10
	orders.BatchSize.PreferredMaxBytes = 512 * 1024
	var kafka genesisconfig.Kafka
	var as = []string{"kafka0:9092", "kafka1:9092", "kafka2:9092", "kafka3:9092"}
	kafka.Brokers = as
	orders.Kafka = kafka
	txCnf.Orderer = &orders

	var app genesisconfig.Application
	txCnf.Application = &app

	var genesisProfile genesisconfig.Profile
	tmpOrder := orders
	tmpOrder.Organizations = orderArray
	genesisProfile.Orderer = &tmpOrder
	var tmpCorsu genesisconfig.Consortium
	tmpCorsu.Organizations = peerArray
	genesisProfile.Consortiums = make(map[string]*genesisconfig.Consortium)
	genesisProfile.Consortiums["SampleConsortium"] = &tmpCorsu

	txCnf.Profiles = make(map[string]*genesisconfig.Profile)
	txCnf.Profiles["TwoOrgsOrdererGenesis"] = &genesisProfile

	var channelProfile genesisconfig.Profile
	channelProfile.Consortium = "SampleConsortium"
	tmpApp := app
	tmpApp.Organizations = peerArray
	channelProfile.Application = &tmpApp

	txCnf.Profiles["TwoOrgsChannel"] = &channelProfile

	by, err := yaml.Marshal(txCnf)
	if err != nil {
		logger.Errorf("GenesisBlockGenerate: marshal tx config failed,%s",err.Error())
		return fmt.Errorf("GenesisBlockGenerate: marshal tx config failed,%s",err.Error())
	}
	cnfTxYaml := serviceRootPath + "/configtx.yaml"
	err = utils.WriteFile(cnfTxYaml, string(by))
	if err != nil {
		logger.Errorf("GenesisBlockGenerate: write tx config to file failed,%s",err)
		return fmt.Errorf("GenesisBlockGenerate: write tx config to file failed,%s",err)
	}
	txPath := serviceRootPath + "/channel-artifacts/"
	err = utils.DirCheck(txPath)
	if err != nil {
		logger.Errorf("GenesisBlockGenerate: artifacts dir create error,%v",err)
		return fmt.Errorf("GenesisBlockGenerate: artifacts dir create error,%v",err)
	}
	bl := txgenerate.CreateGenesisBlock(serviceRootPath, "testchainid", txPath)
	if !bl {
		logger.Errorf("GenesisBlockGenerate: create blockchain genesis block failed")
		return fmt.Errorf("GenesisBlockGenerate: create blockchain genesis block failed")
	}
	return nil  
}