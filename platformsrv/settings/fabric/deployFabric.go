
package fabric

import (
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"time"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/certgenerate/fabric"
	"github.com/wingbaas/platformsrv/k8s/deployfabric" 
	"github.com/wingbaas/platformsrv/settings/fabric/public" 
	"github.com/wingbaas/platformsrv/sdk/sdkfabric"
)

func DeployFabric(p public.DeployPara,chainName string,chainType string)(string,error) {
	if p.DeployType != public.SOLO_FABRIC && p.DeployType != public.KAFKA_FABRIC && p.DeployType != public.RAFT_FABRIC {
		logger.Errorf("DeployFabric: unsupported deploy type")
		return "",fmt.Errorf("DeployFabric: unsupported deploy type")
	}
	if p.CryptoType  != fabric.CRYPTO_ECDSA && p.CryptoType != fabric.CRYPTO_SM { 
		logger.Errorf("DeployFabric: unsupported crypto type")
		return "",fmt.Errorf("DeployFabric: unsupported crypto type")
	}
	blockId := utils.GenerateRandomString(32)
	bytes, err := json.Marshal(p.DeployNetCfg) 
	if err != nil {
		logger.Errorf("DeployFabric: Marshal deploy net config error")
		return "",fmt.Errorf("DeployFabric: Marshal deploy net config error")
	}
	blockCertPath := utils.BAAS_CFG.BlockNetCfgBasePath + blockId 
	err = utils.DirCheck(blockCertPath)
	if err != nil {
		logger.Errorf("DeployFabric: create block config path error")
		return "",fmt.Errorf("DeployFabric: create block config path error")
	}
	bl := fabric.Generate(string(bytes),blockCertPath + "/crypto-config",p.CryptoType) 
	if !bl {
		logger.Errorf("DeployFabric: generate block certification error")
		return "",fmt.Errorf("DeployFabric: generate block certification error")
	}
	if p.DeployType == public.RAFT_FABRIC {
		err = GenerateGenesisBlockRaft(blockCertPath,p.DeployNetCfg) 
	}else{
		err = GenerateGenesisBlock(blockCertPath,p.DeployNetCfg,p.DeployType) 
	}
	if err != nil {
		logger.Errorf("DeployFabric: GenerateGenesisBlock error")
		return "",fmt.Errorf("DeployFabric: GenerateGenesisBlock error")
	}
	dstNfsPath := utils.BAAS_CFG.NfsLocalRootDir + blockId
	err = utils.DirCheck(dstNfsPath)
	if err != nil {
		logger.Errorf("DeployFabric: create nfs cert dir=%s, err=%v",dstNfsPath,err)
		return "",fmt.Errorf("DeployFabric: create nfs cert dir=%s, err=%v",dstNfsPath,err)
	}
	_,err = utils.CopyDir(blockCertPath,dstNfsPath)
	if err != nil {
		logger.Errorf("DeployFabric: copy blockchain cert to nfs error,blockchain id=%s",blockId) 
		return "",fmt.Errorf("DeployFabric: copy blockchain cert to nfs error,blockchain id=%s",blockId) 
	}
	_,err = DeployComponets(p,chainName,blockId,chainType,p.DeployType)
	if err != nil {
		logger.Errorf("DeployFabric: DeployComponets error=%s",err.Error())
		return "",fmt.Errorf("DeployFabric: DeployComponets error=%s",err.Error())
	} 
	bytes, err = json.Marshal(p)
	if err != nil {
		logger.Errorf("DeployFabric: Marshal deploy all config error") 
		return "",fmt.Errorf("DeployFabric: Marshal deploy all config error")
	} 
	blockCfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + blockId + ".json" 
	err = utils.WriteFile(blockCfgFile,string(bytes))
	if err != nil {
		logger.Errorf("DeployFabric: write block config error")
		return "",fmt.Errorf("DeployFabric: write block config error") 
	}
	var sdkCfg sdkfabric.GenerateParaSt
	sdkCfg.ClusterId = p.ClusterId
	sdkCfg.NamespaceId = chainName 
	sdkCfg.BlockId = blockId
	sdkfabric.GenerateOrgCfg(p.DeployNetCfg,sdkCfg) 
	time.Sleep(30*time.Second)
	return blockId,nil  
}

func DeployComponets(p public.DeployPara,chainName string,chainId string,chainType string,consensusType string)(string,error) {
	_,err := deployfabric.CreateNamespace(p.ClusterId,chainName) 
	if err != nil {
		logger.Errorf("DeployComponets: CreateNamespace error") 
		return "",fmt.Errorf("DeployComponets: CreateNamespace error")
	}
	time.Sleep(10 * time.Second)
	//deploy ca
	for _,org := range p.DeployNetCfg.PeerOrgs { 
		caName := strings.ToLower(org.Name + "-ca") 
		caImage,err := utils.GetBlockImage(chainType,p.Version,"ca")
		if err != nil {
			logger.Errorf("DeployComponets: GetBlockImage ca error,chainType=%s version=%s",chainType,p.Version)
			return "",fmt.Errorf("DeployComponets: GetBlockImage ca error,chainType=%s version=%s",chainType,p.Version)
		}
		priKey,err := utils.GetCaPrivateKey(chainId,org.Domain)
		if err != nil {
			logger.Errorf("DeployComponets: GetCaPrivateKey error=%s caName=%s",err.Error(),caName) 
			return "",fmt.Errorf("DeployComponets: GetCaPrivateKey error=%s caName=%s",err.Error(),caName)
		}
		_,err = deployfabric.CreateCaDeployment(p.ClusterId,org.DeployNode,chainName,chainId,caImage,caName,org.Domain,priKey) 
		if err != nil {
			logger.Errorf("DeployComponets: CreateCaDeployment error=%s caName=%s",err.Error(),caName)
			return "",fmt.Errorf("DeployComponets: CreateCaDeployment error=%s caName=%s",err.Error(),caName)
		}
		_,err = deployfabric.CreateCaService(p.ClusterId,chainName,chainId,caName)
		if err != nil {
			logger.Errorf("DeployComponets: CreateCaService error=%s caName=%s",err.Error(),caName) 
			return "",fmt.Errorf("DeployComponets: CreateCaService error=%s caName=%s",err.Error(),caName)
		}	
	}
	if consensusType == public.KAFKA_FABRIC {
		//deploy zookeeper
		zkImage,err := utils.GetBlockImage(chainType,p.Version,"zookeeper") 
		if err != nil {
			logger.Errorf("DeployComponets: GetBlockImage zookeeper error,chainType=%s version=%s",chainType,p.Version)
			return "",fmt.Errorf("DeployComponets: GetBlockImage zookeeper error,chainType=%s version=%s",chainType,p.Version)
		}
		for i:=0; i<public.ZOOK_COUNT; i++ {
			zkName := "zookeeper" + strconv.Itoa(i)   
			_,err = deployfabric.CreateZookeeperDeployment(p.ClusterId,p.DeployNetCfg.ZookeeperDeployNode,chainName,chainId,strconv.Itoa(i+1),zkImage,zkName) 
			if err != nil {
				logger.Errorf("DeployComponets: CreateZkDeployment error=%s zkName=%s",err.Error(),zkName)
				return "",fmt.Errorf("DeployComponets: CreateZkDeployment error=%s zkName=%s",err.Error(),zkName)
			}
			_,err = deployfabric.CreateZookeeperService(p.ClusterId,chainName,chainId,zkName) 
			if err != nil {
				logger.Errorf("DeployComponets: CreateZookeeperService error=%s zkName=%s",err.Error(),zkName) 
				return "",fmt.Errorf("DeployComponets: CreateZookeeperService error=%s zkName=%s",err.Error(),zkName)
			}
		} 
		//deploy kafka
		kafkaImage,err := utils.GetBlockImage(chainType,p.Version,"kafka")
		if err != nil {
			logger.Errorf("DeployComponets: GetBlockImage kafka error,chainType=%s version=%s",chainType,p.Version)
			return "",fmt.Errorf("DeployComponets: GetBlockImage kafka error,chainType=%s version=%s",chainType,p.Version)
		}
		for i:=0; i<public.KAFKA_COUNT; i++ {
			kafkaName := "kafka" + strconv.Itoa(i)
			_,err = deployfabric.CreateKafkaDeployment(p.ClusterId,p.DeployNetCfg.KafkaDeployNode,chainName,chainId,strconv.Itoa(i),kafkaImage,kafkaName) 
			if err != nil {
				logger.Errorf("DeployComponets: CreateKafkaDeployment error=%s kafkaName=%s",err.Error(),kafkaName)
				return "",fmt.Errorf("DeployComponets: CreateKafkaDeployment error=%s kafkaName=%s",err.Error(),kafkaName)
			}
			_,err = deployfabric.CreateKafkaService(p.ClusterId,chainName,chainId,kafkaName)
			if err != nil {
				logger.Errorf("DeployComponets: CreateKafkaService error=%s kafkaName=%s",err.Error(),kafkaName) 
				return "",fmt.Errorf("DeployComponets: CreateKafkaService error=%s kafkaName=%s",err.Error(),kafkaName)
			}
		}
	}
	//deploy order
	orderImage,err := utils.GetBlockImage(chainType,p.Version,"orderer")
	if err != nil {
		logger.Errorf("DeployComponets: GetBlockImage orderer error,chainType=%s version=%s",chainType,p.Version)
		return "",fmt.Errorf("DeployComponets: GetBlockImage orderer error,chainType=%s version=%s",chainType,p.Version)
	}
	for _,org := range p.DeployNetCfg.OrdererOrgs {
		domain := org.Domain
		for _,spec := range org.Specs {
			orderName := spec.Hostname
			if consensusType == public.KAFKA_FABRIC {
				_,err = deployfabric.CreateOrderKafkaDeployment(p.ClusterId,org.DeployNode,chainName,chainId,orderImage,orderName,domain) 
			}else{
				_,err = deployfabric.CreateOrderRaftDeployment(p.ClusterId,org.DeployNode,chainName,chainId,orderImage,orderName,domain)
			}
			if err != nil {
				logger.Errorf("DeployComponets: CreateOrderKafkaDeployment error=%s orderName=%s",err.Error(),orderName)
				return "",fmt.Errorf("DeployComponets: CreateOrderKafkaDeployment error=%s orderName=%s",err.Error(),orderName)
			}
			_,err = deployfabric.CreateOrderService(p.ClusterId,chainName,chainId,orderName)
			if err != nil {
				logger.Errorf("DeployComponets: CreateOrderService error=%s orderName=%s",err.Error(),orderName) 
				return "",fmt.Errorf("DeployComponets: CreateOrderService error=%s orderName=%s",err.Error(),orderName)
			} 
		}
	}
	//deploy peer
	peerImage,err := utils.GetBlockImage(chainType,p.Version,"peer")
	if err != nil {
		logger.Errorf("DeployComponets: GetBlockImage peer error,chainType=%s version=%s",chainType,p.Version)
		return "",fmt.Errorf("DeployComponets: GetBlockImage peer error,chainType=%s version=%s",chainType,p.Version)
	}
	ccenvImage,err := utils.GetBlockImage(chainType,p.Version,"ccenv")
	if err != nil {
		logger.Errorf("DeployComponets: GetBlockImage ccenv error,chainType=%s version=%s",chainType,p.Version)
		return "",fmt.Errorf("DeployComponets: GetBlockImage ccenv error,chainType=%s version=%s",chainType,p.Version)
	}
	baseosImage,err := utils.GetBlockImage(chainType,p.Version,"baseos")
	if err != nil {
		logger.Errorf("DeployComponets: GetBlockImage baseos error,chainType=%s version=%s",chainType,p.Version)
		return "",fmt.Errorf("DeployComponets: GetBlockImage baseos error,chainType=%s version=%s",chainType,p.Version)
	}
	var dp deployfabric.PeerDeploymentPara
	dp.PeerImage = peerImage
	dp.CcenvImage = ccenvImage
	dp.BaseosImage = baseosImage
	for _,org := range p.DeployNetCfg.PeerOrgs {
		dp.PeerDomain = org.Domain 
		dp.OrgName = org.Name
		for _,spec := range org.Specs {
			dp.AnchorPeer = spec.Hostname
			break
		}
		for _,spec :=  range org.Specs { 
			dp.RawPeerName = spec.Hostname
			dp.PeerName = spec.Hostname 
			_,err = deployfabric.CreatePeerDeployment(p.ClusterId,org.DeployNode,chainName,chainId,dp) 
			if err != nil {
				logger.Errorf("DeployComponets: CreatePeerDeployment error=%s",err.Error())
				return "",fmt.Errorf("DeployComponets: CreatePeerDeployment error=%s",err.Error())
			}
			_,err = deployfabric.CreatePeerService(p.ClusterId,chainName,chainId,dp.PeerName)
			if err != nil {
				logger.Errorf("DeployComponets: CreatePeerService error=%s",err.Error())
				return "",fmt.Errorf("DeployComponets: CreatePeerService error=%s",err.Error())
			}
		}
	}
	toolsImage,err := utils.GetBlockImage(chainType,p.Version,"tools")
	if err != nil {
		logger.Errorf("DeployComponets: GetBlockImage tools error,chainType=%s version=%s",chainType,p.Version)
		return "",fmt.Errorf("DeployComponets: GetBlockImage tools error,chainType=%s version=%s",chainType,p.Version)
	}
	var dpt deployfabric.ToolsDeploymentPara
	for _,org := range p.DeployNetCfg.PeerOrgs {
		dpt.PeerDomain = org.Domain 
		dpt.OrgName = org.Name
		for _,spec :=  range org.Specs { 
			dpt.PeerName = spec.Hostname
			break
		}
		break
	}
	dpt.ToolsImage = toolsImage
	_,err = deployfabric.CreateToolsDeployment(p.ClusterId,p.DeployNetCfg.ToolsDeployNode,chainName,chainId,dpt)
	if err != nil {
		logger.Errorf("DeployComponets: CreateToolsDeployment error=%s",err.Error())
		return "",fmt.Errorf("DeployComponets: CreateToolsDeployment error=%s",err.Error())
	}
	return "",nil
}


func DeployOrg(clusterId,chainName,chainId,chainType,version string,orgs []public.OrgSpec)(string,error) {
	//deploy ca
	for _,org := range orgs { 
		caName := strings.ToLower(org.Name + "-ca") 
		caImage,err := utils.GetBlockImage(chainType,version,"ca")
		if err != nil {
			logger.Errorf("DeployOrg: GetBlockImage ca error,chainType=%s version=%s",chainType,version)
			return "",fmt.Errorf("DeployOrg: GetBlockImage ca error,chainType=%s version=%s",chainType,version)
		}
		priKey,err := utils.GetCaPrivateKey(chainId,org.Domain)
		if err != nil {
			logger.Errorf("DeployOrg: GetCaPrivateKey error=%s caName=%s",err.Error(),caName) 
			return "",fmt.Errorf("DeployOrg: GetCaPrivateKey error=%s caName=%s",err.Error(),caName)
		}
		_,err = deployfabric.CreateCaDeployment(clusterId,org.DeployNode,chainName,chainId,caImage,caName,org.Domain,priKey) 
		if err != nil {
			logger.Errorf("DeployOrg: CreateCaDeployment error=%s caName=%s",err.Error(),caName)
			return "",fmt.Errorf("DeployOrg: CreateCaDeployment error=%s caName=%s",err.Error(),caName)
		}
		_,err = deployfabric.CreateCaService(clusterId,chainName,chainId,caName)
		if err != nil {
			logger.Errorf("DeployOrg: CreateCaService error=%s caName=%s",err.Error(),caName)
			return "",fmt.Errorf("DeployOrg: CreateCaService error=%s caName=%s",err.Error(),caName)
		}	
	}
	//deploy peer
	peerImage,err := utils.GetBlockImage(chainType,version,"peer")
	if err != nil {
		logger.Errorf("DeployOrg: GetBlockImage peer error,chainType=%s version=%s",chainType,version)
		return "",fmt.Errorf("DeployOrg: GetBlockImage peer error,chainType=%s version=%s",chainType,version)
	}
	ccenvImage,err := utils.GetBlockImage(chainType,version,"ccenv")
	if err != nil {
		logger.Errorf("DeployOrg: GetBlockImage ccenv error,chainType=%s version=%s",chainType,version)
		return "",fmt.Errorf("DeployOrg: GetBlockImage ccenv error,chainType=%s version=%s",chainType,version)
	}
	baseosImage,err := utils.GetBlockImage(chainType,version,"baseos")
	if err != nil {
		logger.Errorf("DeployOrg: GetBlockImage baseos error,chainType=%s version=%s",chainType,version)
		return "",fmt.Errorf("DeployOrg: GetBlockImage baseos error,chainType=%s version=%s",chainType,version)
	}
	var dp deployfabric.PeerDeploymentPara
	dp.PeerImage = peerImage
	dp.CcenvImage = ccenvImage
	dp.BaseosImage = baseosImage
	for _,org := range orgs {
		dp.PeerDomain = org.Domain 
		dp.OrgName = org.Name
		for _,spec := range org.Specs {
			dp.AnchorPeer = spec.Hostname
			break
		}
		for _,spec :=  range org.Specs { 
			dp.RawPeerName = spec.Hostname
			dp.PeerName = spec.Hostname 
			_,err = deployfabric.CreatePeerDeployment(clusterId,org.DeployNode,chainName,chainId,dp) 
			if err != nil {
				logger.Errorf("DeployOrg: CreatePeerDeployment error=%s",err.Error())
				return "",fmt.Errorf("DeployOrg: CreatePeerDeployment error=%s",err.Error())
			}
			_,err = deployfabric.CreatePeerService(clusterId,chainName,chainId,dp.PeerName)
			if err != nil {
				logger.Errorf("DeployOrg: CreatePeerService error=%s",err.Error())
				return "",fmt.Errorf("DeployOrg: CreatePeerService error=%s",err.Error())
			}
		}
	}
	logger.Debug("DeployOrg success")
	return "",nil
}

func OrgCreateChannel(chainId string,orgName string,channelId string) error {
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return fmt.Errorf("OrgCreateChannel: load chain cfg error,channelId=%s\n", channelId)
	}
	chain,_ := k8s.GetChain(chainId,obj.ClusterId)
	if chain == nil {
		logger.Errorf("OrgCreateChannel: not find this chain,id=%s\n",chainId)
		return fmt.Errorf("OrgCreateChannel: not find this chain,id=%s\n",chainId)
	}
	blockCertPath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId
	if obj.DeployType == public.RAFT_FABRIC {
		err = GenerateChannelTxRaft(blockCertPath,channelId)
	}else if obj.DeployType == public.KAFKA_FABRIC {
		err = GenerateChannelTx(blockCertPath,channelId)
	}
	if err != nil {
		logger.Errorf("OrgCreateChannel: GenerateChannelTx error: %v",err)
		return fmt.Errorf("OrgCreateChannel: GenerateChannelTx error: %v",err)
	}
	var sdkCfg sdkfabric.GenerateParaSt
	sdkCfg.ClusterId = obj.ClusterId
	sdkCfg.NamespaceId = chain.BlockChainName  
	sdkCfg.BlockId = chainId
	sdkCfg.ChannelName = channelId
	_,err = sdkfabric.UpdateChannelCfg(obj.DeployNetCfg,sdkCfg) 
	if err != nil {
		logger.Errorf("OrgCreateChannel: UpdateChannelCfg error")
		return fmt.Errorf("OrgCreateChannel: UpdateChannelCfg error")
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs { 
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	}
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	chSetup := sdkfabric.ChannnelSetup{
		ChannelID: channelId, 
		ChannelConfig: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/channel-artifacts/" + channelId + ".tx",
	}
	err = fSetup.Initialize()
	if err != nil {
		return fmt.Errorf("OrgCreateChannel: init SDK error: %s\n", err)
	}
	defer fSetup.CloseSDK() 

	err = fSetup.CreateChannel(chSetup)
	if err != nil {
		return fmt.Errorf("OrgCreateChannel: create channel error: %s\n", err)
	}
	logger.Debug("OrgCreateChannel success")
	return nil
}

func OrgJoinChannel(chainId string,orgName string,channelId string) error {
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return fmt.Errorf("OrgJoinChannel: load chain cfg error,channelId=%s\n", channelId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	}
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	chSetup := sdkfabric.ChannnelSetup{
		ChannelID: channelId, 
		ChannelConfig: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/channel-artifacts/" + channelId + ".tx",
	}
	err = fSetup.Initialize()
	if err != nil {
		return fmt.Errorf("OrgJoinChannel: init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 

	err = fSetup.JoinChannel(chSetup) 
	if err != nil {
		return fmt.Errorf("OrgJoinChannel: join channel error=%s\n", err)
	}
	logger.Debug("OrgJoinChannel success")
	return nil
}

func OrgDeployChaiCode(chainId string,orgName string,channelId string,ChainCodeID string,ChaincodeVersion string,initArgs []string) error {
	rootPath,err := utils.GetProcessRunRoot()
	if err != nil {
		return fmt.Errorf("OrgDeployChaiCode: get process run path error")
	}
	obj,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		return fmt.Errorf("OrgDeployChaiCode:load chain cfg error,chainId=%s\n", chainId)
	}
	var orderId string
	for _,org := range obj.DeployNetCfg.OrdererOrgs {
		for _,p := range org.Specs {
			orderId = p.Hostname + "." + org.Domain
			break
		}
	} 
	fSetup := sdkfabric.FabricSetup{ 
		OrdererID: orderId,
		OrgAdmin:  "Admin",
		OrgName:   orgName, 
		ChannelId: channelId,
		ConfigFile: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/network-config-" + orgName + ".yaml",
	}
	for _,org := range obj.DeployNetCfg.PeerOrgs {
		if org.Name == orgName {
			for _,p := range org.Specs {
				peer := p.Hostname + "." + org.Domain
				fSetup.Peers = append(fSetup.Peers,peer)
			}
			break
		}
	} 
	chSetup := sdkfabric.ChannnelSetup{
		ChannelID: channelId, 
		ChannelConfig: utils.BAAS_CFG.BlockNetCfgBasePath + chainId + "/channel-artifacts/" + channelId + ".tx",
	}
	ccSetup := sdkfabric.ChaincodeSetup {
		ChainCodeID:     	ChainCodeID,
		ChaincodeGoPath: 	rootPath + "/tmp/" + chainId,
		ChaincodePath:   	ChainCodeID + ChaincodeVersion,
		ChaincodeVersion:	ChaincodeVersion,
		InitOrg:			orgName,
	}
	ccSetup.InitArgs = append(ccSetup.InitArgs,initArgs...)

	err = fSetup.Initialize()
	if err != nil {
		return fmt.Errorf("OrgDeployChaiCode:init SDK failed: org=%s  err=%s\n", orgName,err)
	}
	defer fSetup.CloseSDK() 
	err = fSetup.InstallCC(ccSetup)
	if err != nil {
		return fmt.Errorf("OrgDeployChaiCode: install cc failed,org=%s\n", orgName)
	}
	time.Sleep(3*time.Second) 
	err = fSetup.InstantiateCC(chSetup,ccSetup) 
	if err != nil {
		return fmt.Errorf("OrgDeployChaiCode: instatiate cc failed,org=%s\n", orgName)
	}
	logger.Debug("OrgDeployChaiCode success")
	return nil
}


