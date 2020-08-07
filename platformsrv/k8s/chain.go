
package k8s

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
)

type Chain struct {
	AllianceId	    string  `json:"AllianceId,omitempty"`
	BlockChainId 	string 	`json:"BlockChainId"`
	BlockChainName  string 	`json:"BlockChainName"`
	BlockChainType  string 	`json:"BlockChainType"`
	ClusterId       string 	`json:"ClusterId"`
	Version			string	`json:"Version"`
	Status			int		`json:"Status"` 
}

const ( 
	CHAIN_FILE  			  string = "chain.json" 
	CHAIN_STATUS_FREE   	  int = 0
	CHAIN_STATUS_ORGADD 	  int = 1
	CHAIN_STATUS_PKGCC  	  int = 2 
	CHAIN_STATUS_DEPLOYCC	  int = 3
	CHAIN_STATUS_CALLCC		  int = 4 
	CHAIN_STATUS_QUERY		  int = 5 
	CHAIN_STATUS_CREATEING    int = 101 
	CHAIN_STATUS_CREAT_ERROR  int = 102 
	CHAIN_STATUS_DELETING     int = 201 
)

func AddChain(chain Chain) error {
	cfgPath := utils.BAAS_CFG.BlockNetCfgBasePath + CHAIN_FILE
	exsist,_ := utils.PathExists(cfgPath)
	var chains []Chain
	if exsist {
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&chains)
			if err != nil {
				logger.Errorf("AddChain: unmarshal chains error,%v", err)
				return fmt.Errorf("%v", err)
			}
		}else {
			logger.Errorf("AddChain: load chain list file error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}
	for _,c := range chains {
		if c.BlockChainName == chain.BlockChainName {
			logger.Errorf("%s%s","AddChain: chain name already exsist ",chain.BlockChainName)
			return fmt.Errorf("%s%s","AddChain: chain name already exsist ",chain.BlockChainName)
		}
	}
	chains = append(chains,chain)
	bytes, err := json.Marshal(chains)
	if err != nil {
		logger.Errorf("AddChain: marshal chains error,%v", err)
		return fmt.Errorf("%v", err)
	}
	err = utils.WriteFile(cfgPath,string(bytes))
	if err != nil {
		logger.Errorf("AddChain: Write chain list file error,%v", err)
		return fmt.Errorf("%v", err)
	}
	return nil
}
 
func GetChains(clusterId string)([]Chain,error) {
	cfgPath := utils.BAAS_CFG.BlockNetCfgBasePath + CHAIN_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var chains []Chain
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&chains)
			if err != nil {
				logger.Errorf("GetChains: unmarshal chains error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			var chs []Chain
			for _,c := range chains {
				if c.ClusterId == clusterId {
					chs = append(chs,c)
				}
			}
			return chs,nil
		}else {
			logger.Errorf("GetChains: load chain list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Debug("GetChains: not find chain list file")
	return nil,fmt.Errorf("%s", "GetChains:not find chain list file")
}  

func GetChain(chainId string,clusterId string) (*Chain,error) { 
	chains,err := GetChains(clusterId)
	if err != nil {
		logger.Errorf("GetChain: get chains error,%v", err)
		return nil,fmt.Errorf("%v", err)
	}
	for _,c := range chains {
		if c.BlockChainId == chainId {
			return &c,nil
		}
	}
	logger.Debug("GetChain: chain id not exsist") 
	return nil,nil 
}

func GetChainByName(chainName string,clusterId string) (*Chain,error) { 
	chains,err := GetChains(clusterId)
	if err != nil {
		logger.Debug("GetChainByName: chain not find")
		return nil,fmt.Errorf("%v", err)
	}
	for _,c := range chains {
		if c.BlockChainName == chainName {
			return &c,nil
		}
	}
	logger.Debug("GetChainByName: chain name not exsist %s",chainName)
	return nil,nil
}

func GetChainById(blockChainId string)(*Chain,error) {
	cfgPath := utils.BAAS_CFG.BlockNetCfgBasePath + CHAIN_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var chains []Chain
		bytes,err := utils.LoadFile(cfgPath) 
		if err == nil {
			err = json.Unmarshal(bytes,&chains)
			if err != nil {
				logger.Errorf("GetChainById: unmarshal chains error,%v", err)
				return nil,fmt.Errorf("%v", err)
			} 
			for _,c := range chains {
				if c.BlockChainId == blockChainId {
					return &c,nil
				}
			}
			return nil,fmt.Errorf("not find this chain")
		}else {
			logger.Errorf("GetChainById: load chain list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Debug("GetChainById: not find chain list file")
	return nil,fmt.Errorf("%s", "GetChainById:not find chain list file")
}  

func GetChainOnlyByName(blockChainName string)(*Chain,error) {
	cfgPath := utils.BAAS_CFG.BlockNetCfgBasePath + CHAIN_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var chains []Chain
		bytes,err := utils.LoadFile(cfgPath) 
		if err == nil {
			err = json.Unmarshal(bytes,&chains)
			if err != nil {
				logger.Errorf("GetChainOnlyByName: unmarshal chains error,%v", err)
				return nil,fmt.Errorf("%v", err)
			} 
			for _,c := range chains {
				if c.BlockChainName == blockChainName {
					return &c,nil
				}
			}
			return nil,fmt.Errorf("not find this chain")
		}else {
			logger.Errorf("GetChainOnlyByName: load chain list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Debug("GetChainOnlyByName: not find chain list file")
	return nil,fmt.Errorf("%s", "GetChainOnlyByName:not find chain list file")
}  

func DeleteChain(chain Chain) error {
	cfgPath := utils.BAAS_CFG.BlockNetCfgBasePath + CHAIN_FILE
	exsist,_ := utils.PathExists(cfgPath)
	var chains []Chain
	var newChains []Chain
	if exsist {
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&chains)
			if err != nil {
				logger.Errorf("DeleteChain: unmarshal chains error,%v", err)
				return fmt.Errorf("%v", err)
			}
		}else {
			logger.Errorf("DeleteChain: load chain list file error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}
	for _,c := range chains {
		if c.BlockChainName != chain.BlockChainName {
			tmpChain := c
			newChains = append(newChains,tmpChain)  
		}
	}
	bytes, err := json.Marshal(newChains)
	if err != nil {
		logger.Errorf("DeleteChain: marshal chains error,%v", err)
		return fmt.Errorf("%v", err)
	}
	err = utils.WriteFile(cfgPath,string(bytes))
	if err != nil {
		logger.Errorf("DeleteChain: Write chain list file error,%v", err)
		return fmt.Errorf("%v", err)
	}
	return nil
}

func UpdateChainStatus(chain Chain) error {
	cfgPath := utils.BAAS_CFG.BlockNetCfgBasePath + CHAIN_FILE
	var chains []Chain
	bytes,err := utils.LoadFile(cfgPath)
	if err != nil {
		logger.Errorf("UpdateChainStatus: read chains error,%v", err)
		return fmt.Errorf("%v", err)
	}
	err = json.Unmarshal(bytes,&chains)
	if err != nil {
		logger.Errorf("UpdateChainStatus: unmarshal chains error,%v", err)
		return fmt.Errorf("%v", err)
	}
	for k,c := range chains {
		if c.BlockChainName == chain.BlockChainName {
			chains[k].Status = chain.Status
		}
	}
	bytes, err = json.Marshal(chains)
	if err != nil {
		logger.Errorf("UpdateChainStatus: marshal chains error,%v", err)
		return fmt.Errorf("%v", err)
	}
	err = utils.WriteFile(cfgPath,string(bytes))
	if err != nil {
		logger.Errorf("UpdateChainStatus: Write chain list file error,%v", err)
		return fmt.Errorf("%v", err)
	}
	return nil
}

func CheckChainStatus(chainId string) error {
	cfgPath := utils.BAAS_CFG.BlockNetCfgBasePath + CHAIN_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var chains []Chain
		bytes,err := utils.LoadFile(cfgPath) 
		if err == nil {
			err = json.Unmarshal(bytes,&chains)
			if err != nil {
				logger.Errorf("CheckChainStatus: unmarshal chains error,%v", err)
				return fmt.Errorf("%v", err)
			} 
			for _,c := range chains {
				if c.BlockChainId == chainId {
					if c.Status == CHAIN_STATUS_CREATEING || c.Status == CHAIN_STATUS_CREAT_ERROR || c.Status == CHAIN_STATUS_DELETING {
						logger.Errorf("CheckChainStatus: chain status exception,status=%d,  chainid=%s", c.Status,chainId)
						return fmt.Errorf("CheckChainStatus: chain status exception,status=%d,  chainid=%s", c.Status,chainId)
					}
					return nil 
				}
			}
			return fmt.Errorf("CheckChainStatus:not find this chain")
		}else {
			logger.Errorf("CheckChainStatus: load chain list file error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}
	logger.Debug("CheckChainStatus: not find chain list file")
	return fmt.Errorf("%s", "CheckChainStatus:not find chain list file")
}