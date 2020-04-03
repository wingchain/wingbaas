

package main

import (
	"os"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/api"
)

func main() {
	logger.SetLogInit()
	var cfgFile string
	if len(os.Args) >= 2 {
		cfgFile = os.Args[1]
	}else {
		logger.Debug("wingbaas start with default config file conf/config.json")
		cfgFile = "./conf/config.json"
	}
	utils.BAAS_CFG = new(utils.BaasCfg)
	err := utils.BAAS_CFG.CfgInit(cfgFile)
	if err != nil {
		logger.Errorf("baas config init error")
		return
	}
	if utils.BAAS_CFG == nil { 
		logger.Errorf("baas config nil,wingbaas server exit!")
		return
	}
	err = utils.BAAS_CFG.CfgPathInit()
	if err != nil { 
		logger.Errorf("baas config init error, server exit!")
		return
	}
	err = utils.BAAS_CFG.CfgBlockCfgInit()
	if err != nil { 
		logger.Errorf("baas blockchain config init error, server exit!")
		return
	}
	logger.Debug("start wing baas server")
	api.Start(utils.BAAS_CFG.SrvPort)
	logger.Debug("wingbaas server exit!") 
}
