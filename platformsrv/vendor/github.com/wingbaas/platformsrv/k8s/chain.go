
package k8s

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
)

type Chain struct {
	BlockChainId 	string 	`json:"BlockChainId"`
	BlockChainName  string 	`json:"BlockChainName"`
	BlockChainType  string 	`json:"BlockChainType"`
	ClusterId       string 	`json:"ClusterId"`
}

const (
	CHAIN_FILE  string = "chain.json" 
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
	}else {
		logger.Errorf("GetChains: not find chain list file")
		return nil,fmt.Errorf("%s", "GetChains:  not find chain list file")
	}
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
		logger.Errorf("GetChainByName: get chains error,%v", err)
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
			newChains = append(newChains,c)  
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