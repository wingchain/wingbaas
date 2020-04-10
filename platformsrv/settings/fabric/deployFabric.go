
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
)

const (
	SOLO_FABRIC     string = "SOLO_FABRIC"
	KAFKA_FABRIC    string = "KAFKA_FABRIC"
	RAFT_FABRIC     string = "RAFT_FABRIC"
	ZOOK_COUNT      int    = 3
	KAFKA_COUNT     int   = 3
)

type NodeSpec struct {
    Hostname string   `json:"Hostname"`
}
  
type UsersSpec struct {
    Count int `json:"Count"`
}
  
type OrgSpec struct { 
    Name          string       `json:"Name"`
    Domain        string       `json:"Domain"`
    Specs         []NodeSpec   `json:"Specs"`
    Users         UsersSpec    `json:"Users"`
}
  
type DeployNetConfig struct {
    OrdererOrgs []OrgSpec `json:"OrdererOrgs"`
    PeerOrgs    []OrgSpec `json:"PeerOrgs"`
}

type DeployNodeInfo map[string]string
type DeployNodeGroup map[string]DeployNodeInfo

type DeployPara struct {
	DeployNetCfg        DeployNetConfig    	`json:"DeployNetCfg"`  
	DeployHost       	DeployNodeGroup 	`json:"DeployHost"`
	DeployType       	string          	`json:"DeployType"`
	Version    		 	string          	`json:"Version"`
	CryptoType       	string          	`json:"CryptoType"`
	ClusterId        	string          	`json:"ClusterId"` 
}

func DeployFabric(p DeployPara,chainName string,chainType string)(string,error) {
	if p.DeployType != SOLO_FABRIC && p.DeployType != KAFKA_FABRIC && p.DeployType != RAFT_FABRIC {
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

	if p.DeployType == KAFKA_FABRIC {
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
	return blockId,nil 
}

func DeployComponetsKafka(p DeployPara,chainName string,chainId string,chainType string)(string,error) {
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
	for i:=1; i<=ZOOK_COUNT; i++ {
		zkName := "zookeeper" + strconv.Itoa(i)
		_,err = deployfabric.CreateZookeeperDeployment(p.ClusterId,chainName,chainId,strconv.Itoa(i),zkImage,zkName)
		if err != nil {
			logger.Errorf("DeployComponets: CreateZkDeployment error=%s caName=%s",err.Error(),zkName)
			return "",fmt.Errorf("DeployComponets: CreateZkDeployment error=%s caName=%s",err.Error(),zkName)
		}
		_,err = deployfabric.CreateZookeeperService(p.ClusterId,chainName,chainId,zkName)
		if err != nil {
			logger.Errorf("DeployComponets: CreateZookeeperService error=%s caName=%s",err.Error(),zkName) 
			return "",fmt.Errorf("DeployComponets: CreateZookeeperService error=%s caName=%s",err.Error(),zkName)
		}
	}
	return "",nil 
}


