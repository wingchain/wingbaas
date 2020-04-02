
package utils

import (
	"fmt"
	"encoding/json"
	"path/filepath"
	"os"
	"github.com/wingbaas/platformsrv/logger"
)

type Config interface {
	CfgInit(cfgFile string) error
	CfgPathInit() error
}

type BaasCfg struct {
	SrvAddr        			string		`json:"SrvAddr"`
	SrvPort        			string 		`json:"SrvPort"`
	ClusterCfgPath 			string 		`json:"ClusterCfgPath"`
	ClusterPkiBasePath      string 		`json:"ClusterPkiBasePath"`
	BlockNetCfgBasePath     string 		`json:"BlockNetCfgBasePath"`
	BlockChainVersionCfg    string		`json:"BlockChainVersionCfg"`
}

var BAAS_CFG *BaasCfg = nil
var BLOCK_CFG_MAP map[string]interface{}

func GetProcessRunRoot() (string,error) {
	root,err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Errorf("GetProcessRunRoot: get process run root dir error,%v",err)
		return "",fmt.Errorf("GetProcessRunRoot: get process run root dir error,%v",err)
	}
	return root,nil
}

func (cfg *BaasCfg) CfgInit(cfgFile string) error {
	bytes,err := LoadFile(cfgFile)
	if err !=nil {
		logger.Errorf("CfgInit: load baas config file error")
		return fmt.Errorf("%v", err)
	}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		logger.Errorf("CfgInit: unmarshal error")
		return fmt.Errorf("%v", err)
	}
	return nil 
}

func (cfg *BaasCfg) CfgPathInit() error {

	root,err := GetProcessRunRoot()
	if err != nil {
		logger.Errorf("CfgPathInit: %v",err)
		return fmt.Errorf("CfgPathInit: %v",err)
	}
	cfg.ClusterCfgPath = root + "/" + cfg.ClusterCfgPath
	cfg.ClusterPkiBasePath = root + "/" + cfg.ClusterPkiBasePath
	cfg.BlockNetCfgBasePath = root + "/" + cfg.BlockNetCfgBasePath

	err = DirCheck(cfg.ClusterCfgPath)
	if err != nil {
		logger.Errorf("CfgPathInit: ClusterCfgPath init error")
		return fmt.Errorf("%v", err)
	}
	err = DirCheck(cfg.ClusterPkiBasePath)
	if err != nil {
		logger.Errorf("CfgPathInit: ClusterPkiBasePath init error")
		return fmt.Errorf("%v", err)
	}
	err = DirCheck(cfg.BlockNetCfgBasePath)
	if err != nil {
		logger.Errorf("CfgPathInit: BlockNetCfgBasePath init error")
		return fmt.Errorf("%v", err)
	}
	return nil
} 

func (cfg *BaasCfg) CfgBlockCfgInit() error {
	bytes,err := LoadFile(cfg.BlockChainVersionCfg) 
	if err != nil {
		logger.Errorf("CfgBlockCfgInit: load blockchain version config error")
		return fmt.Errorf("CfgBlockCfgInit: load blockchain version config error")
	}
	BLOCK_CFG_MAP = make(map[string]interface{})
	err = json.Unmarshal(bytes, &BLOCK_CFG_MAP)
	if err != nil {
		logger.Errorf("CfgBlockCfgInit: unmarshal to map error")
		return fmt.Errorf("%v", err)
	}
	return nil 
} 


