

package main

import (
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/api"
)


func main() {
	logger.Debug("start wingbaas main.....")
	logger.SetLogInit()
	utils.BAAS_CFG = new(utils.BaasCfg)
	err := utils.BAAS_CFG.CfgInit("./conf/config.json")
	if err != nil {
		logger.Debug("baas config init error")
		return
	}
	if utils.BAAS_CFG == nil { 
		logger.Debug("baas config nil,wingbaas server exit!")
		return
	}
	err = utils.BAAS_CFG.CfgPathInit()
	if err != nil { 
		logger.Debug("baas config init error, server exit!")
		return
	}
	logger.Debug("start wing baas server")
	api.Start(utils.BAAS_CFG.SrvPort)
	logger.Debug("wingbaas server exit!") 
}
