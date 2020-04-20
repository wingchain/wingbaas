

package fabric

import (
	"fmt"
	"time"
	genesisconfig "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/wingbaas/platformsrv/certgenerate/fabric/gopkg.in/yaml.v2"
	"github.com/wingbaas/platformsrv/logger" 
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric/txgenerate"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
)

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