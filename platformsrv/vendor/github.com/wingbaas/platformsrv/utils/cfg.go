
package utils

import (
	"fmt"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
)

type Config interface {
	CfgInit(cfgFile string) error
	CfgPathInit() error
}

type BaasCfg struct {
	SrvAddr        			string `json:"SrvAddr"`
	SrvPort        			string `json:"SrvPort"`
	ClusterCfgPath 			string `json:"ClusterCfgPath"`
	ClusterPkiBasePath      string `json:"ClusterPkiBasePath"`
	BlockNetCfgBasePath     string `json:"BlockNetCfgBasePath"`
}

var BAAS_CFG *BaasCfg = nil

func (cfg *BaasCfg) CfgInit(cfgFile string) error {
	bytes,err := LoadFile(cfgFile)
	if err !=nil {
		logger.Debug("CfgInit: load baas config file error")
		return fmt.Errorf("%v", err)
	}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		logger.Debug("CfgInit: unmarshal error")
		return fmt.Errorf("%v", err)
	}
	return nil 
}

func (cfg *BaasCfg) CfgPathInit() error {
	err := DirCheck(cfg.ClusterCfgPath)
	if err != nil {
		logger.Debug("CfgPathInit: ClusterCfgPath init error")
		return fmt.Errorf("%v", err)
	}
	err = DirCheck(cfg.ClusterPkiBasePath)
	if err != nil {
		logger.Debug("CfgPathInit: ClusterPkiBasePath init error")
		return fmt.Errorf("%v", err)
	}
	err = DirCheck(cfg.BlockNetCfgBasePath)
	if err != nil {
		logger.Debug("CfgPathInit: BlockNetCfgBasePath init error")
		return fmt.Errorf("%v", err)
	}
	return nil
} 
