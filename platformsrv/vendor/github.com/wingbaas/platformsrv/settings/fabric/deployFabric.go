
package fabric

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/certgenerate/fabric"
	"github.com/wingbaas/platformsrv/k8s/deployfabric"
)

const (
	SOLO_FABRIC     string = "SOLO_FABRIC"
	KAFKA_FABRIC    string = "KAFKA_FABRIC"
	RAFT_FABRIC     string = "RAFT_FABRIC"
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

func DeployFabric(p DeployPara,chainName string)(string,error) {
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
	dstNfsPath := utils.BAAS_CFG.NfsRootDir + blockId
	err = utils.DirCheck(dstNfsPath)
	if err != nil {
		logger.Errorf("DeployFabric: create nfs cert dir=%s, err=%v",dstNfsPath,err)
		return "",fmt.Errorf("DeployFabric: create nfs cert dir=%s, err=%v",dstNfsPath,err)
	}
	// _,err = utils.CopyDir(blockCertPath,dstNfsPath)
	// if err != nil {
	// 	logger.Errorf("DeployFabric: copy blockchain cert to nfs error,blockchain id=%s",blockId) 
	// 	return "",fmt.Errorf("DeployFabric: copy blockchain cert to nfs error,blockchain id=%s",blockId)
	// }

	res,err := deployfabric.CreateNamespace(p.ClusterId,chainName) 
	if err != nil {
		logger.Errorf("DeployFabric: CreateNamespace error") 
		return "",fmt.Errorf("DeployFabric: CreateNamespace error")
	} 
	logger.Debug("DeployFabric:CreateNamespace res=%s",string(res))

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


