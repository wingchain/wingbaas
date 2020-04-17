
package fabric

import (
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"time"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
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
	err = GenerateGenesisBlock(blockCertPath,p.DeployNetCfg,p.DeployType)
	if err != nil {
		logger.Errorf("DeployFabric: GenerateGenesisBlock error")
		return "",fmt.Errorf("DeployFabric: GenerateGenesisBlock error")
	}
	err = GenerateChannelTx(blockCertPath,"mychannel")
	if err != nil {
		logger.Errorf("DeployFabric: GenerateChannelTx error: %v",err)
		return "",fmt.Errorf("DeployFabric: GenerateChannelTx error: %v",err)
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

	if p.DeployType == public.KAFKA_FABRIC {
		_,err = DeployComponetsKafka(p,chainName,blockId,chainType)
		if err != nil {
			logger.Errorf("DeployFabric: DeployComponets error=%s",err.Error())
			return "",fmt.Errorf("DeployFabric: DeployComponets error=%s",err.Error())
		} 
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
	//var sdkCfg sdkfabric.GenerateParaSt
	//sdkfabric.GenerateCfg(p.DeployNetCfg,sdkCfg)
	return blockId,nil 
}

func DeployComponetsKafka(p public.DeployPara,chainName string,chainId string,chainType string)(string,error) {
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
		_,err = deployfabric.CreateCaDeployment(p.ClusterId,chainName,chainId,caImage,caName,org.Domain,priKey)
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
	//deploy zookeeper
	zkImage,err := utils.GetBlockImage(chainType,p.Version,"zookeeper") 
	if err != nil {
		logger.Errorf("DeployComponets: GetBlockImage zookeeper error,chainType=%s version=%s",chainType,p.Version)
		return "",fmt.Errorf("DeployComponets: GetBlockImage zookeeper error,chainType=%s version=%s",chainType,p.Version)
	}
	for i:=1; i<=public.ZOOK_COUNT; i++ {
		zkName := "zookeeper" + strconv.Itoa(i)
		_,err = deployfabric.CreateZookeeperDeployment(p.ClusterId,chainName,chainId,strconv.Itoa(i),zkImage,zkName)
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
	for i:=1; i<=public.KAFKA_COUNT; i++ {
		kafkaName := "kafka" + strconv.Itoa(i)
		_,err = deployfabric.CreateKafkaDeployment(p.ClusterId,chainName,chainId,strconv.Itoa(i),kafkaImage,kafkaName)
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
			_,err = deployfabric.CreateOrderKafkaDeployment(p.ClusterId,chainName,chainId,orderImage,orderName,domain)
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
			_,err = deployfabric.CreatePeerDeployment(p.ClusterId,chainName,chainId,dp)
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
	return "",nil 
}


